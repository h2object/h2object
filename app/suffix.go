package app

import (
	"strings"
	"github.com/h2object/h2object/httpext"
)

type handler func(ctx *context, ctrl *ext.Controller) bool

var handlers map[string]handler

func suffix_filter(ctx *context, c *ext.Controller, filters []filter) {
	if do, ok := handlers[strings.ToLower(c.Request.Suffix())]; ok {
		if done := do(ctx, c); done {
			ctx.Info("request (%s) (%s) done by suffix", c.Request.MethodToLower(), c.Request.URI())
			return
		}
	}
	filters[0](ctx, c, filters[1:])
}

func init() {
	handlers = make(map[string]handler)
	handlers["md"] = do_markdown
	handlers["html"] = do_template
	handlers["conf"] = do_configure
	handlers["export"] = do_export
	handlers["xml"] = do_xml
}