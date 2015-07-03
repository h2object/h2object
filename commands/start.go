package commands

import (
	"os"
	"path"
	"syscall"
	"github.com/h2object/pidfile"
	"github.com/h2object/h2object/app"
)

func run(application *app.Application) {
	pid, err := pidfile.New(path.Join(application.Options.Root, "h2object.pid"))
	if err != nil {

	}
	defer pid.Kill()

	exitc := make(chan int)
	signc := make(chan os.Signal, 1)

	go func(){
		for {
			sig := <- signc
			switch sig {
			case syscall.SIGHUP:
				application.Refresh()
				continue
			default:
				exitc <- 1
				break
			}
		}
	}()

	application.Main()
	<- exitc 
	application.Exit()
}

func start(application *app.Application, daemon bool) {
	if !daemon {
		run(application)
	} else {
		daemonize(application)
	}
}