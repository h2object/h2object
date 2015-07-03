package commands

import (
	"fmt"
	"github.com/h2object/h2object/app"
)

func daemonize(application *app.Application) {
	fmt.Println("[h2object] warn: ", "windows not support daemon mode")
	run(application)	
}