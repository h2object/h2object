package app

import (
	"github.com/h2object/h2object/httpext"
)

func routes_filter(ctx *context, c *ext.Controller, filters []filter) {
	if c.Request.URI() == "/stats" {
		ctx.Info("request (%s) (%s) done by routes", c.Request.MethodToLower(), c.Request.URI())
		c.Json(ctx.app.statistics.Data())
		return
	}

	filters[0](ctx, c, filters[1:])
}
