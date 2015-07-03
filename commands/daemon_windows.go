package commands

import (
	"github.com/h2object/h2object/app"
)

func daemonize(application *app.Application) {
	if mode {
		fmt.Println("[h2object] warn: ", "windows not support daemon mode")
	}
	run(application)	
}