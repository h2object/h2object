package app

import (
	"path"
	"strings"
	"net/http"
	"github.com/h2object/h2object/httpext"
)

type service struct{
	ctx  	*context
}

func new_service(ctx *context) *service {
	return &service{
		ctx: ctx,
	}
}

func (srv *service) init() error {
	append_filter(acl_filter)
	append_filter(params_filter)
	append_filter(routes_filter)
	append_filter(suffix_filter)	
	append_filter(static_filter)
	append_filter(nofound_filter)
	append_filter(nil_filter)
	return nil
}

func (srv *service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request := ext.NewRequest(req)
	response := ext.NewResponse(w)
	controller := ext.NewController(request, response)

	if strings.HasSuffix(controller.Request.URI(), "/") {
		controller.Request.URL.Path = path.Join(controller.Request.URL.Path, srv.ctx.index)
	}

	filters[0](srv.ctx, controller, filters[1:])
}
