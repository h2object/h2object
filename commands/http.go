package commands

import (
	"github.com/codegangsta/cli"
	"github.com/h2object/pidfile"
	
	"path"
	"path/filepath"
	"fmt"
	"os"
	"time"
	"github.com/h2object/h2object/app"
	"github.com/h2object/h2object/log"
)

const (
	success        = "\t\t\t\t\t[  \033[32mOK\033[0m  ]" // Show colored "OK"
	failed         = "\t\t\t\t\t[\033[31mFAILED\033[0m]" // Show colored "FAILED"
)

func httpStartCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	if workdir == "" {
		fmt.Println("unknown working directory, please use -w to provide.")
		os.Exit(1)
	}
	// daemon 
	daemon := ctx.GlobalBool("daemon")
	
	// options
	options := app.NewOptions(ctx.GlobalString("host"),ctx.GlobalInt("port"))
	if err := options.Prepare(workdir); err != nil {
		fmt.Println("options prepare failed:", err)
		os.Exit(1)
	}

	refresh := ctx.GlobalString("refresh")
	options.SetRefreshDefault(refresh, time.Minute * 10)

	// configs
	configs, err := app.LoadCONFIG(path.Join(options.Root, "h2object.conf"))
	if err != nil {
		configs = app.DefaultCONFIG()
		if err := configs.Save(path.Join(options.Root, "h2object.conf")); err != nil {
			fmt.Println("h2object.conf saving failed:", err)
			os.Exit(1)
		}
	}
	
	logger := log.NewH2OLogger()
	defer logger.Close()
	logger.SetConsole(true)
	
	configs.SetSection("logs")
	fenable := configs.BoolDefault("file.enable", false)
	fname := configs.StringDefault("file.name", "h2o.log")
	flevel := configs.StringDefault("file.level", "info")
	fsize := configs.IntDefault("file.rotate_max_size", 1024*1024*1024)
	fline := configs.IntDefault("file.rotate_max_line", 102400)
	fdaily := configs.BoolDefault("file.rotate_daily", true)
	fn := path.Join(options.LogsRoot, fname)
	if fenable == true {
		logger.SetFileLog(fn, flevel, fsize, fline, fdaily)	
	}	

	application := app.NewApplication(options, configs, logger)

	if err := application.Init(); err != nil {
		fmt.Println("[h2object] init failed:", err)
		os.Exit(1)
	}

	start(application, daemon)	
}

func httpStopCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	if workdir == "" {
		fmt.Println("unknown working directory, please use -w to provide.")
		os.Exit(1)
	}

	// directory
	directory, err := filepath.Abs(workdir)
	if err != nil {
		fmt.Println("workdir:", err)
		return
	}

	pid, err := pidfile.Load(path.Join(directory, "h2object.pid"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := pid.Kill(); err != nil {
		fmt.Println(err.Error())
		return	
	}

	fmt.Println("[h2object] stop", success)
	return	
}

func httpReloadCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	if workdir == "" {
		fmt.Println("unknown working directory, please use -w to provide.")
		os.Exit(1)
	}

	// directory
	directory, err := filepath.Abs(workdir)
	if err != nil {
		fmt.Println("workdir:", err)
		return
	}

	pid, err := pidfile.Load(path.Join(directory, "h2object.pid"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := pid.HUP(); err != nil {
		fmt.Println(err.Error())
		return	
	}

	fmt.Println("[h2object] reload", success)
	return	
}

