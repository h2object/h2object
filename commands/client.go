package commands

import (
	"io"
	"fmt"	
	"bytes"
	"strconv"
	"net/url"
	"net/http"
	"encoding/json"
	"mime/multipart"
	"github.com/h2object/rpc"
	"github.com/h2object/pb"
)

type Analyser struct{
}

type JSResponse struct{
	Err string `json:"error"`
}

func (a Analyser) Analyse(ret interface{}, resp *http.Response) error {
	defer resp.Body.Close()

	if resp.StatusCode/100 == 2 {
		if resp.ContentLength != 0 && ret != nil {
			if err := json.NewDecoder(resp.Body).Decode(ret); err != nil {
				return err
			}
		}
	} else {
		if resp.ContentLength != 0 {
			var js JSResponse
			if err := json.NewDecoder(resp.Body).Decode(&js); err != nil {
				return err
			}
			return fmt.Errorf("code: %d, reason: %s", resp.StatusCode, js.Err)
		}
		return fmt.Errorf("code: %d, reason: failed", resp.StatusCode)
	}
	return nil
}


type Client struct{
	workdir string
	addr 	string
	conn *rpc.Client
}

func NewClient(workdir string, host string, port int) *Client {
	connection := rpc.NewClient(Analyser{})
	cli := &Client{
		workdir: workdir,
		addr: fmt.Sprintf("%s:%d", host, port),	
		conn: connection,
	}
	return cli
}

func (cli *Client) SignUp(authid, password string, auth interface{}) error {
	URL := rpc.BuildHttpURL(cli.addr, "/auth/email/signup", nil)	
	if err := cli.conn.PostForm(nil, URL, map[string][]string{
		"authid":{authid},
		"password":{password},
	}, auth); err != nil {
		return err
	}
	return nil
}


func (cli *Client) SignInPassword(authid, password, remember string, auth interface{}) error {
	URL := rpc.BuildHttpURL(cli.addr, "/auth/email/signinpassword", nil)	
	if err := cli.conn.PostForm(nil, URL, map[string][]string{
		"authid":{authid},
		"password":{password},
		"remember":{remember},
	}, auth); err != nil {
		return err
	}
	return nil
}

func (cli *Client) SignInSecret(authid, secret, remember string, auth interface{}) error {
	URL := rpc.BuildHttpURL(cli.addr, "/auth/email/signinsecret", nil)	
	if err := cli.conn.PostForm(nil, URL, map[string][]string{
		"authid":{authid},
		"secret":{secret},
		"remember":{remember},
	}, auth); err != nil {
		return err
	}
	return nil
}

func (cli *Client) SignOff(token string) error {
	
	params := url.Values{}
	params.Set("token", token)

	URL := rpc.BuildHttpURL(cli.addr, "/auth/signoff", params)	
	if err := cli.conn.Get(nil, URL, nil); err != nil {
		return err
	}
	return nil
}

func (cli *Client) Auth(token string, auth interface{}) error {
	params := url.Values{}
	params.Set("token", token)

	URL := rpc.BuildHttpURL(cli.addr, "/auth/info", params)	
	if err := cli.conn.Get(nil, URL, auth); err != nil {
		return err
	}
	return nil
}

func (cli *Client) ThemeSearch(token string, keyword string, catagory int64, page int64, size int64, themes interface{}) error {
	params := url.Values{}
	params.Set("keyword", keyword)
	if token != "" {
		params.Set("token", token)
	}
	if catagory != 0 {
		params.Set("catagory", strconv.FormatInt(catagory, 64))
	}
	if page != 0 {
		params.Set("page", strconv.FormatInt(page, 64))
	}
	if size != 0 {
		params.Set("size", strconv.FormatInt(size, 64))
	}
	URL := rpc.BuildHttpURL(cli.addr, "/themes/search", params)	
	if err := cli.conn.Get(nil, URL, themes); err != nil {
		return err
	}
	return nil
}

type Package struct{
	Provider string
	Name string
	Version string
	Catagory int64
	Description string
	Price float64
	ArchiveReader io.ReadCloser
	ArchiveName	  string
	ArchiveLen	int64
}

func (cli *Client) ThemePush(token string, pkg *Package) error {
	params := url.Values{}
	params.Set("token", token)

	URL := rpc.BuildHttpURL(cli.addr, 
		   fmt.Sprintf("/themes/push/%s/%s/%s", pkg.Provider, pkg.Name, pkg.Version), 
		   params)

	var b bytes.Buffer
	wr := multipart.NewWriter(&b)

	if err := wr.WriteField("catagory", fmt.Sprintf("%d", pkg.Catagory)); err != nil {
		return err
	}
	if err := wr.WriteField("description", pkg.Description); err != nil {
		return err
	}
	if err := wr.WriteField("price", fmt.Sprintf("%f", pkg.Price)); err != nil {
		return err
	}
	p, err := wr.CreateFormFile("file", pkg.ArchiveName)
	if err != nil {
		return err
	}
	defer pkg.ArchiveReader.Close()

	if _, err := io.Copy(p, pkg.ArchiveReader); err != nil {
		return err
	}
	if err := wr.Close(); err != nil {
		return err
	}

	bar := pb.New(b.Len()).SetUnits(pb.U_BYTES)
	bar.Prefix(fmt.Sprintf("%s/%s:%s ", pkg.Provider, pkg.Name, pkg.Version))
	bar.Start()

	// create multi writer
	rd := pb.NewPbReader(&b, bar)

	if err := cli.conn.Post(nil, URL, wr.FormDataContentType(), rd, int64(b.Len()), nil); err != nil {
		return err
	}
	bar.FinishPrint(fmt.Sprintf("%s/%s:%s pushed succussfully.", pkg.Provider, pkg.Name, pkg.Version))
	return nil
}

func (cli *Client) ThemePull(token string, pkg *Package) error {
	params := url.Values{}
	params.Set("token", token)

	URL := rpc.BuildHttpURL(cli.addr, 
		   fmt.Sprintf("/themes/pull/%s/%s/%s", pkg.Provider, pkg.Name, pkg.Version), 
		   params)
	resp, err := cli.conn.GetResponse(nil, URL)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		if resp.ContentLength != 0 {
			var js JSResponse
			if err := json.NewDecoder(resp.Body).Decode(&js); err != nil {
				return err
			}
			return fmt.Errorf("code: %d, reason: %s", resp.StatusCode, js.Err)
		}
		return fmt.Errorf("code: %d, reason: failed", resp.StatusCode)
	}
	pkg.ArchiveLen = resp.ContentLength
	pkg.ArchiveReader = resp.Body	
	return nil
}

