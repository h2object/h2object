package commands 

import (
	"github.com/codegangsta/cli"
)

const version = "0.1.0"
const author = "support"
const support = "support@h2object.io"

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
			Usage: "h2object http server host",
		},
		cli.IntFlag{
			Name: "port, p",
			Value: 9000,
			Usage: "h2object http server port",
		},
		cli.StringFlag{
			Name: "workdir, w",
			Value: "",
			Usage: "h2object working directory",
		},
		cli.BoolFlag{
			Name: "daemon, d",
			Usage: "run at daemon mode",
		},
	}

	//! app commands
	app.Commands = []cli.Command{
		{
			Name:  "auth",
			Usage: "auth commands default @ h2object.io",
			Subcommands: []cli.Command{
				{
					Name:  "new",
					Usage: "auth account sign up.",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "api.h2object.io",
							Usage: "h2object auth server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "h2object auth server port",
						},
					},
					Action: func(ctx *cli.Context) {
						authNewCommand(ctx)	
					},
				},
				{
					Name:  "login",
					Usage: "auth login with username & password",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "api.h2object.io",
							Usage: "h2object auth server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "h2object auth server port",
						},
					},
					Action: func(ctx *cli.Context) {
						authLoginCommand(ctx)	
					},
				},
				{
					Name:  "status",
					Usage: "auth status",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "api.h2object.io",
							Usage: "h2object auth server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "h2object auth server port",
						},
					},
					Action: func(ctx *cli.Context) {
						authStatusCommand(ctx)	
					},
				},
				{
					Name:  "logout",
					Usage: "auth logout",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "api.h2object.io",
							Usage: "h2object auth server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "h2object auth server port",
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
					Usage: "start http server",
					Action: func(ctx *cli.Context) {
						httpStartCommand(ctx)	
					},
				},
				{
					Name:  "reload",
					Usage: "reload http server",
					Action: func(ctx *cli.Context) {
						httpReloadCommand(ctx)	
					},
				},
				{
					Name:  "stop",
					Usage: "stop http server",
					Action: func(ctx *cli.Context) {
						httpStopCommand(ctx)	
					},
				},
			},
		},
		{
			Name:  "deploy",
			Usage: "deploy commands @ local",
			Subcommands: []cli.Command{
				{
					Name:  "push",
					Usage: "push local application to remote http server",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "",
							Usage: "h2object push remote host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "h2object push remote port",
						},
					},
					Action: func(ctx *cli.Context) {
						deployPushCommand(ctx)	
					},
				},
				{
					Name:  "pull",
					Usage: "pull remote application to the local http server",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "",
							Usage: "h2object pull remote host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "h2object pull remote port",
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
			Usage: "theme commands default @ h2object.io",
			Subcommands: []cli.Command{
				{
					Name:  "search",
					Usage: "search themes with keyword from theme server",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "api.h2object.io",
							Usage: "h2object theme server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "h2object theme server port",
						},
						cli.IntFlag{
							Name: "page, p",
							Value: 0,
							Usage: "theme search page number",
						},
						cli.IntFlag{
							Name: "size, s",
							Value: 0,
							Usage: "theme search page size",
						},
						cli.IntFlag{
							Name: "catagory, c",
							Value: 0,
							Usage: "theme search catagory",
						},
					},
					Action: func(ctx *cli.Context) {
						themeSearchCommand(ctx)	
					},
				},
				{
					Name:  "pull",
					Usage: "pull dest theme from theme server",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "api.h2object.io",
							Usage: "h2object theme server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "h2object theme server port",
						},
					},
					Action: func(ctx *cli.Context) {
						themePullCommand(ctx)	
					},
				},
				{
					Name:  "push",
					Usage: "push current theme to theme server",
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "Host, H",
							Value: "api.h2object.io",
							Usage: "h2object theme server host",
						},
						cli.IntFlag{
							Name: "Port, P",
							Value: 80,
							Usage: "h2object theme server port",
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
