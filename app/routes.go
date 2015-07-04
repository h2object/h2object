package app

import (
	"github.com/h2object/h2object/httpext"
)

var ItWorks string = `
<html>
	<head>
		<title>h2object works</title>
	</head>
	<body>
		<h3>H2OBJECT</h3>
		====<br>
		<p>It works!!!</p>
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
		c.Html(ItWorks)
		return true
	}
	
	if c.Request.URI() == "/stats" {
		c.Json(ctx.app.statistics.Data())
		return true
	}
	
	return false
}
