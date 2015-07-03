package ext

import (	
	"io"	
	"time"
	"net/http"
	"github.com/h2object/h2object/template"
)

type Controller struct{
	Request  *Request
	Response *Response
}

func NewController(request *Request, response *Response) *Controller {
	return &Controller{
		Request:  request,
		Response: response,
	}
}

func (r *Controller) JsonError(status int, err error) {
	JsonErrorResult{
		code: status,
		err: err,
		pretty: r.Request.Pretty,
	}.Apply(r.Request, r.Response)
}

func (r *Controller) TextError(status int, err error) {
	TextErrorResult{
		code: status,
		err: err,
	}.Apply(r.Request, r.Response)
}

func (r *Controller) Json(data interface{}) {
	JsonResult{
		obj: data,
		pretty: r.Request.Pretty,
		callback: r.Request.Callback,
	}.Apply(r.Request, r.Response)
}

func (r *Controller) Xml(data interface{}) {
	XmlResult{
		obj: data,
		pretty: r.Request.Pretty,
	}.Apply(r.Request, r.Response)
}

func (r *Controller) Text(text string) {
	TextResult{
		text: text,
	}.Apply(r.Request, r.Response)
}

func (r *Controller) Html(html string) {
	HtmlResult{
		html: html,
	}.Apply(r.Request, r.Response)
}

func (r *Controller) Read(content_type string, rd io.Reader, sz int64) {
	ReaderResult{
		content_type: content_type,
		read: rd,
		size: sz,
	}.Apply(r.Request, r.Response)
}

func (r *Controller) Template(t template.Template, data interface{}) {
	TemplateResult{
		t: t,
		data: data,
	}.Apply(r.Request, r.Response)
}

func (r *Controller) Redirect(uri string, query string) {
	RedirectResult{
		url: uri,
		rawQuery: query,
	}.Apply(r.Request, r.Response)
}

func (r *Controller) File(path string) {
	http.ServeFile(r.Response.Out, r.Request.Request, path)
}

func (r *Controller) Binary(name string, rd io.Reader, size int64) {
	BinaryResult{
		Reader: rd,
		Name: name,
		Length: size,
		Delivery: Attachment,
		ModTime: time.Now(),
	}.Apply(r.Request, r.Response)
}




