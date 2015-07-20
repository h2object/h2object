package app

import (
	"errors"
	"net/http"
	"strings"
	"github.com/h2object/h2object/httpext"
)

func acl_filter(ctx *context, c *ext.Controller, filters []filter) {
	if done := do_authentic(ctx, c); done {
		ctx.Info("request (%s) (%s) done by acl", c.Request.MethodToLower(), c.Request.URI())
		return
	}

	filters[0](ctx, c, filters[1:])
}

func do_authentic(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request
	
	required := false
	switch r.MethodToLower() {
	case "get":
		switch r.Suffix() {
		case "page":
			fallthrough
		case "click":
			fallthrough	
		case "system":
			required = true
		}

		if r.URI() == "/stats" {
			required = true	
		}
	case "put":
		required = true
		if ctx.storage_full() {
			ctrl.JsonError(http.StatusForbidden, errors.New("application storage reach max limit."))
			return true
		}
	case "delete":
		required = true
	}

	token := r.Param("token")
	if token == "" {
		authorization := r.Header.Get("Authorization")
		if strings.HasPrefix(authorization, "H2OBJECT ") {
			token = authorization[len("H2OBJECT "):]
		}
	}
	if required {
		if token != ctx.signature {
			ctrl.JsonError(http.StatusUnauthorized, errors.New("require administrator right"))
			return true
		}
	}

	return false
}