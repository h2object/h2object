package commands 

import (
	"github.com/codegangsta/cli"
)

const version = "1.0.2"
const author = ""
const support = "liujianping@h2object.io"

func App() *cli.App {
	app := cli.NewApp()

	//! app settings
	app.Name = "h2object"
	app.Usage = "another fast & flexible static website generator & deployment tool"
	app.Version = version
	app.Author = author
	app.Email = support

	//! app flags
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "host, l",
			Value: "0.0.0.0",
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
		cli.StringFlag{
			Name: "refresh, r",
			Value: "10m",
			Usage: "refresh interval",
		},
		cli.StringFlag{
			Name: "storage, s",
			Value: "",
			Usage: "storage max size(kb, kib, mb, mib, gb, gib, tb, tib, pb, pib)",
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
					Flags: []cli.Flag {
						cli.StringFlag{
							Name: "appid",
							Value: "",
							Usage: "h2object application id",
						},
						cli.StringFlag{
							Name: "secret",
							Value: "",
							Usage: "h2object application secret",
						},
					},
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
					Usage: "push files @ local => remote ",
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
					Usage: "pull files @ local <= remote",
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
				{
					Name:  "delete",
					Usage: "delete files @ remote",
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
						deployDeleteCommand(ctx)	
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
		{
			Name:  "container",
			Usage: "container commands @ h2object.io",
			Subcommands: []cli.Command{
				{
					Name:  "get",
					Usage: "get containers @ remote with option argments: [container-id] ...",
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
						containerGetCommand(ctx)	
					},
				},
				{
					Name:  "create",
					Usage: "create container @ remote with required argment: [invitation]",
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
						containerCreateCommand(ctx)	
					},
				},
				{
					Name:  "start",
					Usage: "start container @ remote with required argment: [container-id]",
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
						containerStartCommand(ctx)	
					},
				},
				{
					Name:  "stop",
					Usage: "stop container @ remote with required argment: [container-id]",
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
						containerStopCommand(ctx)	
					},
				},
				{
					Name:  "restart",
					Usage: "restart container @ remote with required argment: [container-id]",
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
						containerRestartCommand(ctx)	
					},
				},
				{
					Name:  "pause",
					Usage: "pause container @ remote with required argment: [container-id]",
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
						containerPauseCommand(ctx)	
					},
				},
				{
					Name:  "unpause",
					Usage: "unpause container @ remote with required argment: [container-id]",
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
						containerUnpauseCommand(ctx)	
					},
				},
				{
					Name:  "domain",
					Usage: "set container domain @ remote  with required argments: [container-id] [custom-domain]",
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
						containerDomainCommand(ctx)	
					},
				},
			},
		},
	}

	return app
}
