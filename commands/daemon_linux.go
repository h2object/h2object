package commands

import (
	"os"
	"fmt"
	"path"
	"syscall"
	"github.com/h2object/pidfile"
	"github.com/h2object/go-daemon"
	"github.com/h2object/h2object/app"
)

func daemonize(application *app.Application) {
	termHandler := func(sig os.Signal) error {
		application.Exit()
		return nil
	}
	reloadHandler := func(sig os.Signal) error {
		application.Refresh()
		return nil
	}
	daemon.AddCommand(nil, syscall.SIGQUIT, termHandler)
	daemon.AddCommand(nil, syscall.SIGTERM, termHandler)
	daemon.AddCommand(nil, syscall.SIGHUP, reloadHandler)
	
	cntxt := &daemon.Context{
		PidFileName: "",
		PidFilePerm: 0644,
		LogFileName: "",
		LogFilePerm: 0640,
		WorkDir:     "",
		Umask:       027,
		Args:        []string{},
	}
	d, err := cntxt.Reborn()
	if err != nil {
		fmt.Println("daemonize ", err)
		os.Exit(1)
	}
	if d != nil {
		os.Exit(1)
	}
	defer cntxt.Release()
	
	pid, err := pidfile.New(path.Join(application.Options.Root, "h2object.pid"))
	if err != nil {
		fmt.Println("[h2object] failed:", err)
		os.Exit(1)
	}
	defer pid.Kill()

	application.Main()
	fmt.Println("[h2object] start", success)

	err = daemon.ServeSignals()
	if err != nil {
		fmt.Println("[h2object] failed:", err)
		os.Exit(1)
	}
}
