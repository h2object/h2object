package commands

import (
	"os"
	"io"
	"fmt"
	"text/tabwriter"
	"github.com/codegangsta/cli"
)

type OUTContainer struct{
	ID string 	`json:"id"`
	SystemDomain string `json:"system_domain"`
	CustomDomain string `json:"custom_domain"`
	Port int64 `json:"port"`
	AppID 	string `json:"appid"`
	AppSecret string `json:"secret"`
	Status string `json:"status"`
	Storage string `json:"storage"`
	Version string `json:"version"`
}

func printOUTContainer(out io.Writer, container OUTContainer) {
	fmt.Fprintf(out, "------------------------\t\n")
	fmt.Fprintf(out, "Container ID:\t%s\n", container.ID)
	fmt.Fprintf(out, "Application Version:\t%s\n", container.Version)
	fmt.Fprintf(out, "Container System Domain:\t%s\n", container.SystemDomain)
	fmt.Fprintf(out, "Container Custom Domain:\t%s\n", container.CustomDomain)
	fmt.Fprintf(out, "Container Status:\t%s\n", container.Status)
	fmt.Fprintf(out, "Container AppID:\t%s\n", container.AppID)
	fmt.Fprintf(out, "Container AppSecret:\t%s\n", container.AppSecret)
	fmt.Fprintf(out, "Container Storage Max:\t%s\n", container.Storage)
	
	fmt.Fprintf(out, "\n")
}

func containerGetCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	stdout := os.Stdout
	stderr := os.Stderr

	host := ctx.String("Host")
	port := ctx.Int("Port")
	
	var token string = ""
	config, err := LoadConfigFile(workdir)
	if err == nil {
		token = config.Auth.Token
	}
	if token == "" {
		fmt.Fprintln(stdout, "container command need login first. ")
		os.Exit(1)
	}

	client := NewClient(workdir, host, port)

	var containers []OUTContainer
	if err := client.GetContainers(token, ctx.Args(), &containers); err != nil {
		fmt.Fprintln(stderr, err.Error())
		os.Exit(1)
	}

	w := tabwriter.NewWriter(stdout, 10, 0, 2, ' ', 0)

	fmt.Fprintf(w, "Container ID\tStatus\tSystem Domain\tCustom Domain\tPort\tApplication ID\tApplciation Secret\t\n")
	for _, c := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\t%s\n", c.ID, c.Status, c.SystemDomain, c.CustomDomain, c.Port, c.AppID, c.AppSecret)
	}
	fmt.Fprintf(w, "\n")
}

func containerCreateCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	stdout := os.Stdout
	
	host := ctx.String("Host")
	port := ctx.Int("Port")
	
	var token string = ""
	config, err := LoadConfigFile(workdir)
	if err == nil {
		token = config.Auth.Token
	}
	if token == "" {
		fmt.Fprintln(stdout, "container command need login first. ")
		os.Exit(1)
	}

	var invitation string = ""
	if len(ctx.Args()) == 1 {
		invitation = ctx.Args()[0]
	}

	client := NewClient(workdir, host, port)
	var container OUTContainer
	
	if err := client.CreateContainer(token, invitation, version, &container); err != nil {
		fmt.Fprintln(stdout, err.Error())
		os.Exit(1)
	}
	
	fmt.Fprintf(stdout, "Container Created OK.\n")
	printOUTContainer(tabwriter.NewWriter(stdout, 10, 0, 2, ' ', 0), container)
}

func containerStartCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	stdout := os.Stdout
	
	host := ctx.String("Host")
	port := ctx.Int("Port")
	
	var token string = ""
	config, err := LoadConfigFile(workdir)
	if err == nil {
		token = config.Auth.Token
	}
	if token == "" {
		fmt.Fprintln(stdout, "container command need login first. ")
		os.Exit(1)
	}

	if len(ctx.Args()) != 1 {
		fmt.Fprintln(stdout, "command need [container id] args. ")
		os.Exit(1)
	}

	client := NewClient(workdir, host, port)
	var container OUTContainer

	if err := client.StartContainer(token, ctx.Args()[0], version, &container); err != nil {
		fmt.Fprintln(stdout, err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(stdout, "Container Start OK.\n")
	printOUTContainer(tabwriter.NewWriter(stdout, 10, 0, 2, ' ', 0), container)
}

func containerStopCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	stdout := os.Stdout
	
	host := ctx.String("Host")
	port := ctx.Int("Port")
	
	var token string = ""
	config, err := LoadConfigFile(workdir)
	if err == nil {
		token = config.Auth.Token
	}
	if token == "" {
		fmt.Fprintln(stdout, "container command need login first. ")
		os.Exit(1)
	}

	if len(ctx.Args()) != 1 {
		fmt.Fprintln(stdout, "command need [container id] args. ")
		os.Exit(1)
	}

	client := NewClient(workdir, host, port)
	var container OUTContainer

	if err := client.StopContainer(token, ctx.Args()[0], &container); err != nil {
		fmt.Fprintln(stdout, err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(stdout, "Container Stop OK.\n")
	printOUTContainer(tabwriter.NewWriter(stdout, 10, 0, 2, ' ', 0), container)
}

func containerRestartCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	stdout := os.Stdout
	
	host := ctx.String("Host")
	port := ctx.Int("Port")
	
	var token string = ""
	config, err := LoadConfigFile(workdir)
	if err == nil {
		token = config.Auth.Token
	}
	if token == "" {
		fmt.Fprintln(stdout, "container command need login first. ")
		os.Exit(1)
	}

	if len(ctx.Args()) != 1 {
		fmt.Fprintln(stdout, "command need [container id] args. ")
		os.Exit(1)
	}

	client := NewClient(workdir, host, port)
	var container OUTContainer

	if err := client.RestartContainer(token, ctx.Args()[0], &container); err != nil {
		fmt.Fprintln(stdout, err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(stdout, "Container Restart OK.\n")
	printOUTContainer(tabwriter.NewWriter(stdout, 10, 0, 2, ' ', 0), container)
}

func containerPauseCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	stdout := os.Stdout
	
	host := ctx.String("Host")
	port := ctx.Int("Port")
	
	var token string = ""
	config, err := LoadConfigFile(workdir)
	if err == nil {
		token = config.Auth.Token
	}
	if token == "" {
		fmt.Fprintln(stdout, "container command need login first. ")
		os.Exit(1)
	}

	if len(ctx.Args()) != 1 {
		fmt.Fprintln(stdout, "command need [container id] args. ")
		os.Exit(1)
	}

	client := NewClient(workdir, host, port)
	var container OUTContainer

	if err := client.PauseContainer(token, ctx.Args()[0], &container); err != nil {
		fmt.Fprintln(stdout, err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(stdout, "Container Pause OK.\n")
	printOUTContainer(tabwriter.NewWriter(stdout, 10, 0, 2, ' ', 0), container)
}

func containerUnpauseCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	stdout := os.Stdout
	
	host := ctx.String("Host")
	port := ctx.Int("Port")
	
	var token string = ""
	config, err := LoadConfigFile(workdir)
	if err == nil {
		token = config.Auth.Token
	}
	if token == "" {
		fmt.Fprintln(stdout, "container command need login first. ")
		os.Exit(1)
	}

	if len(ctx.Args()) != 1 {
		fmt.Fprintln(stdout, "command need [container id] args. ")
		os.Exit(1)
	}

	client := NewClient(workdir, host, port)
	var container OUTContainer

	if err := client.UnpauseContainer(token, ctx.Args()[0], &container); err != nil {
		fmt.Fprintln(stdout, err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(stdout, "Container Unpause OK.\n")
	printOUTContainer(tabwriter.NewWriter(stdout, 10, 0, 2, ' ', 0), container)
}

func containerDomainCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	stdout := os.Stdout
	
	host := ctx.String("Host")
	port := ctx.Int("Port")
	
	var token string = ""
	config, err := LoadConfigFile(workdir)
	if err == nil {
		token = config.Auth.Token
	}
	if token == "" {
		fmt.Fprintln(stdout, "container command need login first. ")
		os.Exit(1)
	}

	if len(ctx.Args()) != 2 {
		fmt.Fprintln(stdout, "command need [container id] [domain] args. ")
		os.Exit(1)
	}

	client := NewClient(workdir, host, port)
	var container OUTContainer

	if err := client.DomainContainer(token, ctx.Args()[0], ctx.Args()[1], &container); err != nil {
		fmt.Fprintln(stdout, err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(stdout, "Container Domain Set OK.\n")
	printOUTContainer(tabwriter.NewWriter(stdout, 10, 0, 2, ' ', 0), container)
}
