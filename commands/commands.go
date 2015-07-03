package commands 

import (
	"github.com/codegangsta/cli"
)

const version = "1.0.0"
const author = ""
const support = "liujianping@h2object.io"

func App() *cli.App {
	app := cli.NewApp()

	//! app settings
	app.Name = "h2object"
	app.Usage = "another http server with themes to build statics & markdown web site"
	app.Version = version
	app.Author = author
	app.Email = support

	//! app flags
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "host, l",
			Value: "127.0.0.1",
			Usage: "local server host",
		},
		cli.IntFlag{
			Name: "port, p",
			Value: 9000,
			Usage: "local server port",
		},
		cli.StringFlag{
			Name: "workdir, w",
			Value: "",
			Usage: "local working directory",
		},
		cli.BoolFlag{
			Name: "daemon, d",
			Usage: "run @ daemon mode",
		},
	}

	//! app commands
	app.Commands = []cli.Command{
		{
			Name:  "auth",
			Usage: "auth commands @ h2object.io",
			Subcommands: []cli.Command{
				{
					Name:  "new",
					Usage: "sign up @ remote",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "api.h2object.io",
							Usage: "remote server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "remote server port",
						},
					},
					Action: func(ctx *cli.Context) {
						authNewCommand(ctx)	
					},
				},
				{
					Name:  "login",
					Usage: "sign in @ remote",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "api.h2object.io",
							Usage: "remote server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "remote server port",
						},
					},
					Action: func(ctx *cli.Context) {
						authLoginCommand(ctx)	
					},
				},
				{
					Name:  "status",
					Usage: "sign info @ remote",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "api.h2object.io",
							Usage: "remote server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "remote server port",
						},
					},
					Action: func(ctx *cli.Context) {
						authStatusCommand(ctx)	
					},
				},
				{
					Name:  "logout",
					Usage: "sign off @ remote",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "api.h2object.io",
							Usage: "remote server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "remote server port",
						},
					},
					Action: func(ctx *cli.Context) {
						authLogoutCommand(ctx)	
					},
				},
			},
		},
		{
			Name:  "http",
			Usage: "http commands @ local",
			Subcommands: []cli.Command{
				{
					Name:  "start",
					Usage: "local server start",
					Action: func(ctx *cli.Context) {
						httpStartCommand(ctx)	
					},
				},
				{
					Name:  "reload",
					Usage: "local server reload",
					Action: func(ctx *cli.Context) {
						httpReloadCommand(ctx)	
					},
				},
				{
					Name:  "stop",
					Usage: "local server stop",
					Action: func(ctx *cli.Context) {
						httpStopCommand(ctx)	
					},
				},
			},
		},
		{
			Name:  "deploy",
			Usage: "deploy commands @ local <=> remote",
			Subcommands: []cli.Command{
				{
					Name:  "push",
					Usage: "deploy @ local => remote ",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "",
							Usage: "remote server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "remote server port",
						},
					},
					Action: func(ctx *cli.Context) {
						deployPushCommand(ctx)	
					},
				},
				{
					Name:  "pull",
					Usage: "deploy @ local <= remote",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "",
							Usage: "remote server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "remote server port",
						},
					},
					Action: func(ctx *cli.Context) {
						deployPullCommand(ctx)	
					},
				},
			},
		},
		{
			Name:  "theme",
			Usage: "theme commands @ h2object.io",
			Subcommands: []cli.Command{
				{
					Name:  "search",
					Usage: "search themes @ remote",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "api.h2object.io",
							Usage: "remote server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "remote server port",
						},
						cli.IntFlag{
							Name: "page, p",
							Value: 0,
							Usage: "search page number",
						},
						cli.IntFlag{
							Name: "size, s",
							Value: 50,
							Usage: "search page size",
						},
						cli.IntFlag{
							Name: "catagory, c",
							Value: 0,
							Usage: "search catagory",
						},
					},
					Action: func(ctx *cli.Context) {
						themeSearchCommand(ctx)	
					},
				},
				{
					Name:  "pull",
					Usage: "pull a theme @ local <= remote",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "api.h2object.io",
							Usage: "remote server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "remote server port",
						},
					},
					Action: func(ctx *cli.Context) {
						themePullCommand(ctx)	
					},
				},
				{
					Name:  "push",
					Usage: "push current theme @ local => remote",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "api.h2object.io",
							Usage: "remote server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "remote server port",
						},
					},
					Action: func(ctx *cli.Context) {
						themePushCommand(ctx)	
					},
				},
			},
		},
	}

	return app
}
