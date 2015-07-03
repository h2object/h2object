package ext

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"io"
	"fmt"
	"strconv"
	"time"
	"bytes"
	"github.com/h2object/content-type"
	"github.com/h2object/h2object/template"
)

type Result interface {
	Apply(req *Request, resp *Response)
}

type TextErrorResult struct {
	code int
	err  error
}
// This method is used when the template loader or error template is not available.
func (r TextErrorResult) Apply(req *Request, resp *Response) {
	resp.WriteHeader(r.code, "text/plain; charset=utf-8")
	resp.Out.Write([]byte(r.err.Error()))
}

type JsonErrorResult struct {
	code int
	err  error
	pretty bool
}

func (r JsonErrorResult) Apply(req *Request, resp *Response) {
	m := map[string]interface{}{
		"error": r.err.Error(),
	}

	var b []byte
	if r.pretty {
		b, _ = json.MarshalIndent(m, "", "  ")
	} else {
		b, _ = json.Marshal(m)
	}

	req.ResultErr = r.err
	resp.WriteHeader(r.code, "application/json; charset=utf-8")	
	resp.Out.Write(b)
}

type JsonResult struct {
	obj      interface{}
	pretty	 bool
	callback string
}

func (r JsonResult) Apply(req *Request, resp *Response) {
	var b []byte
	var err error
	if r.pretty {
		b, err = json.MarshalIndent(r.obj, "", "  ")
	} else {
		b, err = json.Marshal(r.obj)
	}

	if err != nil {
		JsonErrorResult{code: http.StatusInternalServerError, err: err, pretty: req.Pretty}.Apply(req, resp)
		return
	}

	if r.callback == "" {
		resp.WriteHeader(http.StatusOK, "application/json; charset=utf-8")
		resp.Out.Write(b)
		return
	}

	resp.WriteHeader(http.StatusOK, "application/javascript; charset=utf-8")
	resp.Out.Write([]byte(r.callback + "("))
	resp.Out.Write(b)
	resp.Out.Write([]byte(");"))
}

type XmlResult struct {
	obj interface{}
	pretty bool
}

func (r XmlResult) Apply(req *Request, resp *Response) {
	var b []byte
	var err error
	if r.pretty {
		b, err = xml.MarshalIndent(r.obj, "", "  ")
	} else {
		b, err = xml.Marshal(r.obj)
	}

	if err != nil {
		JsonErrorResult{code: http.StatusInternalServerError, err: err, pretty: req.Pretty}.Apply(req, resp)
		return
	}

	resp.WriteHeader(http.StatusOK, "application/xml; charset=utf-8")
	resp.Out.Write(b)
}

type TextResult struct {
	text string
}

func (r TextResult) Apply(req *Request, resp *Response) {
	resp.WriteHeader(http.StatusOK, "text/plain; charset=utf-8")
	resp.Out.Write([]byte(r.text))
}

type HtmlResult struct {
	html string
}

func (r HtmlResult) Apply(req *Request, resp *Response) {
	resp.WriteHeader(http.StatusOK, "text/html; charset=utf-8")
	resp.Out.Write([]byte(r.html))
}

type ReaderResult struct {
	content_type 	string
	read io.Reader 
	size int64
}

func (r ReaderResult) Apply(req *Request, resp *Response) {
	resp.Out.Header().Set("Content-Length", strconv.FormatInt(r.size, 10))
	resp.WriteHeader(http.StatusOK, r.content_type)	
	io.Copy(resp.Out, r.read)
}


type TemplateResult struct {
	t   template.Template
	data interface{}
}

func (r TemplateResult) Apply(req *Request, resp *Response) {
	// Handle panics when rendering templates.
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		TextErrorResult{
	// 			code: http.StatusInternalServerError, 
	// 			err: fmt.Errorf("template (%s) error: %s",r.t.Name(), err), 
	// 		}.Apply(req, resp)
	// 		log.Println("template result failed: ", err)
	// 	}
	// }()


	var b bytes.Buffer
	if err := r.t.Render(&b, r.data); err != nil {
		TextErrorResult{
			code: http.StatusInternalServerError, 
			err: err, 
		}.Apply(req, resp)
		return		
	}

	resp.WriteHeader(http.StatusOK, "text/html; charset=utf-8")
	b.WriteTo(resp.Out)
	return
}

type RedirectResult struct {
	url string
	rawQuery string
}

func (r RedirectResult) Apply(req *Request, resp *Response) {
	if r.rawQuery != "" {
		resp.Out.Header().Set("Location", r.url + "?" + r.rawQuery)
	} else {
		resp.Out.Header().Set("Location", r.url)
	}	
	resp.WriteHeader(http.StatusFound, "")
}

type ContentDisposition string

var (
	Attachment ContentDisposition = "attachment"
	Inline     ContentDisposition = "inline"
)

type BinaryResult struct {
	Reader   io.Reader
	Name     string
	Length   int64
	Delivery ContentDisposition
	ModTime  time.Time
}

func (r BinaryResult) Apply(req *Request, resp *Response) {
	disposition := string(r.Delivery)
	if r.Name != "" {
		disposition += fmt.Sprintf(`; filename="%s"`, r.Name)
	}
	resp.Out.Header().Set("Content-Disposition", disposition)

	// If we have a ReadSeeker, delegate to http.ServeContent
	if rs, ok := r.Reader.(io.ReadSeeker); ok {
		// http.ServeContent doesn't know about response.ContentType, so we set the respective header.
		if resp.ContentType != "" {
			resp.Out.Header().Set("Content-Type", resp.ContentType)
		} else {
			contentType := content_type.DefaultContentTypeHelper.ContentTypeByFilename(r.Name)
			resp.Out.Header().Set("Content-Type", contentType)
		}
		http.ServeContent(resp.Out, req.Request, r.Name, r.ModTime, rs)
	} else {
		// Else, do a simple io.Copy.
		if r.Length != -1 {
			resp.Out.Header().Set("Content-Length", strconv.FormatInt(r.Length, 10))
		}
		resp.WriteHeader(http.StatusOK, content_type.DefaultContentTypeHelper.ContentTypeByFilename(r.Name))
		io.Copy(resp.Out, r.Reader)
	}

	// Close the Reader if we can
	if v, ok := r.Reader.(io.Closer); ok {
		v.Close()
	}
}
