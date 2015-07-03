package app

import (
	"errors"
	"net/http"
	"github.com/h2object/h2object/httpext"
)

type filter func(ctx *context, c *ext.Controller, filterChain []filter)

// Filters is the default set of global filters.
// It may be set by the application on initialization.
var filters = []filter{}

func append_filter(f filter) {
	filters = append(filters, f)
}

func nil_filter(ctx *context, c *ext.Controller, filterChain []filter) {}

func nofound_filter(ctx *context, c *ext.Controller, filterChain []filter) {
	ctx.Info("request (%s) (%s) not found", c.Request.MethodToLower(), c.Request.URI())
	c.JsonError(http.StatusNotFound, errors.New("not found"))
}