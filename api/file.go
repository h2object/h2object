package api

import (
	"io"
	"os"
	"path"
	"fmt"
	"encoding/json"
	"github.com/h2object/rpc"
	"github.com/h2object/content-type"
)

func (h2o *Client) Download(l Logger, auth Auth, dest_uri string) (io.ReadCloser, int64, error)  {
	URL := rpc.BuildHttpURL(h2o.addr, dest_uri, nil)

	h2o.Lock()
	defer h2o.Unlock()

	h2o.conn.Prepare(auth)
	defer h2o.conn.Prepare(nil)

	resp, err := h2o.conn.GetResponse(l, URL)
	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode/100 != 2 {
		if resp.ContentLength != 0 {
			var ret1 rpc.ErrorRet
			if err := json.NewDecoder(resp.Body).Decode(&ret1); err != nil {
				return nil, 0, err
			}
			return nil, 0, fmt.Errorf("code: %d, reason: %s", resp.StatusCode, ret1.Error)
		}
		return nil, 0, fmt.Errorf("code: %d, reason: failed", resp.StatusCode)
	}
	return resp.Body, resp.ContentLength, nil
}


func (h2o *Client) UploadFile(l Logger, auth Auth, dest_uri string, file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fd.Close()

	st, err := fd.Stat()
	if err != nil {
		return err
	}

	_, filename := path.Split(file)
	contentType := content_type.DefaultContentTypeHelper.ContentTypeByFilename(filename)

	return h2o.Upload(l, auth, dest_uri, contentType, fd, st.Size())
}

func (h2o *Client) Upload(l Logger, auth Auth, dest_uri string, contentType string, rd io.Reader, sz int64) error {
	URL := rpc.BuildHttpURL(h2o.addr, dest_uri, nil)

	h2o.Lock()
	defer h2o.Unlock()

	h2o.conn.Prepare(auth)
	defer h2o.conn.Prepare(nil)

	if err := h2o.conn.Put(l, URL, contentType, rd, sz, nil); err != nil {
		return err
	}
	return nil
}
