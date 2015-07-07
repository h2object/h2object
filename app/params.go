package app

import (
	"strconv"
	"net/http"
	"github.com/h2object/h2object/httpext"
)


func params_filter(ctx *context, c *ext.Controller, filters []filter) {
	if done := do_params(ctx, c); done {
		ctx.Info("request (%s) (%s) done by params", c.Request.MethodToLower(), c.Request.URI())
		return
	}
	filters[0](ctx, c, filters[1:])
}

func do_params(ctx *context, ctrl *ext.Controller) bool {
	qr := ctrl.Request.Param("qrcode")
	if qr != "" {
		if size, err := strconv.ParseInt(qr, 10, 64); err == nil {
			file, err := ctx.qrcode_path(ctrl.Request.URL, int(size))
			if err != nil {
				ctrl.JsonError(http.StatusInternalServerError, err)
				return true
			}
			ctrl.File(file)
			return true
		}
	}
	return false
}


