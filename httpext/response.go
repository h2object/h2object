package ext

import (
	"net/http"
)

type Response struct {
	Status      int
	ContentType string
	Out http.ResponseWriter
}

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{
		Out: w,
	}
}

func (resp *Response) SetCookie(cookie *http.Cookie) {
	http.SetCookie(resp.Out, cookie)
}

// Write the header (for now, just the status code).
// The status may be set directly by the application (c.Response.Status = 501).
// if it isn't, then fall back to the provided status code.
func (resp *Response) WriteHeader(defaultStatusCode int, defaultContentType string) {
	if resp.Status == 0 {
		resp.Status = defaultStatusCode
	}
	if resp.ContentType == "" {
		resp.ContentType = defaultContentType
	}
	resp.Out.Header().Set("Content-Type", resp.ContentType)
	resp.Out.WriteHeader(resp.Status)
}
