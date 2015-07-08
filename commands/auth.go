package commands

import (
	"os"
	"io"
	"fmt"
	"bufio"
	"time"
	"path/filepath"
	"github.com/h2object/h2object/util"
	"github.com/codegangsta/cli"
	"github.com/docker/docker/pkg/term"
)

func prompt(out io.Writer, prompt string, configDefault string) {
	if configDefault == "" {
		fmt.Fprintf(out, "%s: ", prompt)
	} else {
		fmt.Fprintf(out, "%s (%s): ", prompt, configDefault)
	}
}

func readInput(in io.Reader, out io.Writer) string {
	reader := bufio.NewReader(in)
	line, _, err := reader.ReadLine()
	if err != nil {
		fmt.Fprintln(out, err.Error())
		os.Exit(1)
	}
	return string(line)
}

func authNewCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	// stdin, stdout, stderr := term.StdStreams()

	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr

	fmt.Fprintf(stdout, "new account need to be actived through email, please use your valid email address.\n")
	prompt(stdout, "account email", "")
	authid := readInput(stdin, stdout)
	
	v := util.ValidEmail()
	if v.IsSatisfied(authid) != true {
		fmt.Fprintln(stderr, v.DefaultMessage())
		os.Exit(1)	
	}

	inFd, _ := term.GetFdInfo(stdin)
	oldState, _ := term.SaveState(inFd)
	prompt(stdout, "account password", "")
	term.DisableEcho(inFd, oldState)
	password := readInput(stdin, stdout)

	prompt(stdout, "\nconfirm password", "")
	password2 := readInput(stdin, stdout)
	term.RestoreTerminal(inFd, oldState)

	if password != password2 {
		fmt.Fprintln(stderr, "\npassword not same, try again aplease.")
		os.Exit(1)
	}

	host := ctx.String("Host")
	port := ctx.Int("Port")
	client := NewClient(workdir, host, port)

	if err := client.SignUp(authid, password, nil); err != nil {
		fmt.Fprintln(stderr, "\n", err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(stdout, "\naccount (%s) signup succeed. \n", authid)
	os.Exit(0)
}

func authLoginCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	// directory
	directory, err := filepath.Abs(workdir)
	if err != nil {
		fmt.Println("workdir:", err)
		return
	}

	// stdin, stdout, stderr := term.StdStreams()
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr

	config, err := LoadConfigFile(directory)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		os.Exit(1)
	}

	if config.Auth.Token != "" {
		fmt.Fprintln(stdout, "account already logined.")
		os.Exit(1)
	}

	host := ctx.String("Host")
	port := ctx.Int("Port")
	client := NewClient(directory, host, port)

	var remember string = ""
	if len(ctx.Args()) > 0 {
		remember = ctx.Args()[0]
	}

	var authid, password string
	if config.Auth.Secret != "" {

		prompt(stdout, "account email", "")
		authid = readInput(stdin, stdout)

		v := util.ValidEmail()
		if v.IsSatisfied(authid) != true {
			fmt.Fprintln(stderr, v.DefaultMessage())
			os.Exit(1)	
		}

		// signinsecret
		if err := client.SignInSecret(authid, config.Auth.Secret, remember, config.Auth); err != nil {
			fmt.Fprintln(stderr, err.Error())
			os.Exit(1)	
		}

		if err := config.Save(); err != nil {
			fmt.Fprintln(stderr, err.Error())
			os.Exit(1)		
		}
		// save
		fmt.Fprintf(stdout, "account (%s) login succeed. \n", authid)
		os.Exit(0)

	} else {
		prompt(stdout, "account email", "")
		authid = readInput(stdin, stdout)

		v := util.ValidEmail()
		if v.IsSatisfied(authid) != true {
			fmt.Fprintln(stderr, v.DefaultMessage())
			os.Exit(1)	
		}

		inFd, _ := term.GetFdInfo(stdin)
		oldState, _ := term.SaveState(inFd)
		prompt(stdout, "account password", "")
		term.DisableEcho(inFd, oldState)
		password = readInput(stdin, stdout)
		term.RestoreTerminal(inFd, oldState)

		// signinpassword
		if err := client.SignInPassword(authid, password, remember, config.Auth); err != nil {
			fmt.Fprintln(stderr, err.Error())
			os.Exit(1)	
		}

		if err := config.Save(); err != nil {
			fmt.Fprintln(stderr, err.Error())
			os.Exit(1)		
		}
		// save
		fmt.Fprintf(stdout, "account (%s) login succeed. \n", authid)
		os.Exit(0)
	}
}

func authLogoutCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	// directory
	directory, err := filepath.Abs(workdir)
	if err != nil {
		fmt.Println("workdir:", err)
		return
	}
	// _, stdout, stderr := term.StdStreams()
	// stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr

	config, err := LoadConfigFile(directory)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		os.Exit(1)
	}

	if config.Auth.Token == "" {
		fmt.Fprintln(stdout, "none account logined. ")
		os.Exit(0)
	}

	host := ctx.String("Host")
	port := ctx.Int("Port")
	client := NewClient(directory, host, port)

	client.SignOff(config.Auth.Token)
	config.Remove()
	fmt.Fprintln(stdout, "account logout. ")
	os.Exit(0)
}

func authStatusCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	// directory
	directory, err := filepath.Abs(workdir)
	if err != nil {
		fmt.Println("workdir:", err)
		return
	}
	// _, stdout, stderr := term.StdStreams()
	// stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr

	config, err := LoadConfigFile(directory)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		os.Exit(1)
	}

	if config.Auth.Token == "" && config.Auth.Secret == "" {
		fmt.Fprintln(stdout, "none account logined. ")
		os.Exit(0)
	}	

	host := ctx.String("Host")
	port := ctx.Int("Port")
	client := NewClient(directory, host, port)

	if err := client.Auth(config.Auth.Token, config.Auth); err != nil {
		fmt.Fprintln(stderr, err.Error())
		os.Exit(1)
	}

	expire := time.Unix(config.Auth.ExpireAt, 0)

	fmt.Fprintf(stdout, "current account: %s\n", config.Auth.AuthID)
	fmt.Fprintf(stdout, "token secret: %s\n", config.Auth.Token)
	fmt.Fprintf(stdout, "token expire: %s\n", expire.Format("2006-01-02 15:04:05"))
}