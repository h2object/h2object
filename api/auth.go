package api

import (
	"fmt"
	"net/http"
	"github.com/h2object/h2object/util"
)

type HeaderAuth struct{
	headers http.Header
}

func NewHeaderAuth() *HeaderAuth {
	return &HeaderAuth{
		headers: http.Header{},
	}
}

func (hds *HeaderAuth) Add(key, value string) {
	hds.headers.Add(key, value)
}

func (hds *HeaderAuth) Set(key, value string) {
	hds.headers.Set(key, value)
}

func (hds *HeaderAuth) Del(key string) {
	hds.headers.Del(key)
}

func (hds *HeaderAuth) Do(req *http.Request) *http.Request {
	for k, vs := range hds.headers {
		req.Header.Del(k)
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	return req
}

func NewAdminAuth(appid, secret string) Auth {
	auth := NewHeaderAuth()
	auth.Del("Authorization")
	auth.Add("Authorization", fmt.Sprintf("H2OBJECT %s", util.SignString(secret, appid)))
	return auth
}
