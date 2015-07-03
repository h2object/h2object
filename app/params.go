package app

import (
	"github.com/h2object/h2object/httpext"
)


func params_filter(ctx *context, c *ext.Controller, filters []filter) {
	
	filters[0](ctx, c, filters[1:])
}



