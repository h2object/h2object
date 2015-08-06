package app

import (
	"fmt"
	"github.com/h2object/h2object/httpext"
)

var ItWorks string = `
<html>
	<head>
		<title>h2object running</title>
	</head>
	<body style="margin: 0px 0px 0px 0px; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;">
		<div style="background-color:#A9F16C; min-height: 360px;">
			<br><br>
			<div style="width: 360px; margin-right: auto; margin-left: auto;">
			<h1 style="font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 60px;line-height: 1;">It works!</h1>
			<br>
			<p>h2object current version: %s</p>
			</div>
		</div>
	</body>
</html>
`

func routes_filter(ctx *context, c *ext.Controller, filters []filter) {
	if done := routes(ctx, c); done {
		ctx.Info("request (%s) (%s) done by routes", c.Request.MethodToLower(), c.Request.URI())
		return
	}
	filters[0](ctx, c, filters[1:])
}

func routes(ctx *context, c *ext.Controller) bool {
	if c.Request.URI() == "/" {
		c.Html(fmt.Sprintf(ItWorks, ctx.app.version))
		return true
	}
	
	if c.Request.URI() == "/stats" {
		c.Json(ctx.app.statistics.Data())
		return true
	}
	
	return false
}
