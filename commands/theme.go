package commands

import (
	"os"
	// "io"
	"strings"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"text/tabwriter"
	"github.com/codegangsta/cli"
	"github.com/docker/docker/pkg/term"
	"github.com/docker/docker/pkg/archive"
	"github.com/h2object/pb"
	"github.com/h2object/h2object/app"
)

type Theme struct{
	ID 			int64 	`db:"id" json:"id"`
	Stat    	int64   `db:"stat" json:"stat"`  
	Catagory 	int64 	`db:"catagory" json:"catagory"`
	Account 	int64	`db:"account" json:"account"`
	Provider 	string 	`db:"provider" json:"provider"`
	Name     	string  `db:"name" json:"name"`
	Picture 	string  `db:"picture" json:"picture"`
	Description string 	`db:"description" json:"description"`
	Version 	string 	`db:"version" json:"version"`
	Url 		string 	`db:"url" json:"url"`
	Price 		float64 `db:"price" json:"price"`
	Apps    	int64	`db:"apps" json:"apps"`	 
	CreateAt    int64	`db:"create_at" json:"create_at"`
	ModifyAt    int64	`db:"modify_at" json:"modify_at"`
}


func catagory_print(catagory int64) string {
	switch catagory {
	case 0:
		return fmt.Sprintf("free")	
	case 1:
		return fmt.Sprintf("member")
	case 2:
		return fmt.Sprintf("buy")
	}
	return fmt.Sprintf("%d", catagory)
}

func status_print(stat int64) string {
	switch stat {
	case 0:
		return fmt.Sprintf("private")
	case 1:
		return fmt.Sprintf("publish")
	}
	return fmt.Sprintf("%d", stat)
}

func themeSearchCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	_, stdout, stderr := term.StdStreams()

	host := ctx.String("Host")
	port := ctx.Int("Port")
	page := int64(ctx.Int("page"))
	size := int64(ctx.Int("size"))
	catagory := int64(ctx.Int("catagory"))

	var token string = ""
	config, err := LoadConfigFile(workdir)
	if err == nil {
		token = config.Auth.Token
	}	
	client := NewClient(workdir, host, port)

	keyword := ""
	if len(ctx.Args()) > 0 {
		keyword = ctx.Args()[0]
	}

	type Result struct{
		Total int64 `json:"total"`
		Count int64 `json:"count"`
		Themes []Theme `json:"themes"`
	}
	var result Result
	if err := client.ThemeSearch(token, keyword, catagory, page, size, &result); err != nil {
		fmt.Fprintln(stderr, err.Error())
		os.Exit(1)
	}

	w := tabwriter.NewWriter(stdout, 10, 0, 2, ' ', 0)

	fmt.Fprintf(w, "provider/name:version\tstatus\tcatagory\tdownloads\tdescription\n")
	fmt.Fprintf(w, "---------------------\t-------\t--------\t---------\t-----------\n")
	for _, theme := range result.Themes {
		fmt.Fprintf(w, "%s/%s:%s\t%s\t%s\t%d\t%s\n", theme.Provider, theme.Name, theme.Version, status_print(theme.Stat),catagory_print(theme.Catagory), theme.Apps, theme.Description)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "The repository has %d themes totally. ", result.Total)
	if result.Total > result.Count {
		fmt.Fprintf(w, "You can use -page or -size to show the left themes\n")
	} else {
		fmt.Fprintf(w, "\n")
	}
}

func themePushCommand(ctx *cli.Context) {
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

	_, stdout, stderr := term.StdStreams()
	
	// auth check
	config, err := LoadConfigFile(directory)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		os.Exit(1)
	}

	if config.Auth.Token == "" || config.Auth.Secret == "" {
		fmt.Fprintln(stdout, "theme push need login first. ")
		os.Exit(1)
	}
	
	h2oconf, err := app.LoadCONFIG(path.Join(directory, "h2object.conf"))
	if err != nil {
		fmt.Fprintln(stdout, err)
		os.Exit(1)
	}

	host := ctx.String("Host")
	port := ctx.Int("Port")
	client := NewClient(directory, host, port)

	var pkg Package
	h2oconf.SetSection("theme")
	pkg.Provider = h2oconf.StringDefault("provider", "")
	if pkg.Provider == "" {
		fmt.Fprintln(stderr, "please set h2object.conf [theme]provider first.")
		os.Exit(1)
	}
	pkg.Name = h2oconf.StringDefault("name", "")
	if pkg.Name == "" {
		fmt.Fprintln(stderr, "please set h2object.conf [theme]name first.")
		os.Exit(1)
	}

	pkg.Description = h2oconf.StringDefault("description", "")
	if pkg.Description == "" {
		fmt.Fprintln(stderr, "please set h2object.conf [theme]description first.")
		os.Exit(1)
	}

	pkg.Version = h2oconf.StringDefault("version", "1.0.0")
	pkg.Price = h2oconf.FloatDefault("price", 0.0)
	pkg.Catagory = int64(h2oconf.IntDefault("catagory", 0))

	var tarOpt archive.TarOptions
	tarOpt.ExcludePatterns = append(tarOpt.ExcludePatterns, ".h2object")
	tarOpt.Compression = archive.Gzip
	rd, err := archive.TarWithOptions(directory, &tarOpt)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		os.Exit(1)
	}
	defer rd.Close()
	pkg.ArchiveReader = rd
	pkg.ArchiveName = pkg.Version + ".tar.gz"

	if err := client.ThemePush(config.Auth.Token, &pkg); err != nil {
		fmt.Fprintln(stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func parse_theme(pkg string) (string, string, string, error) {
	ps := strings.Split(pkg, ":")
	if len(ps) != 2 {
		return "", "", "", errors.New("absent version")
	}
	version := ps[1]

	ns := strings.Split(ps[0], "/")
	if len(ns) != 2 {
		return "", "", "", errors.New("unknown provider and name")
	}
	return ns[0], ns[1], version, nil
}

func themePullCommand(ctx *cli.Context) {
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
	_, _, stderr := term.StdStreams()
	
	if len(ctx.Args()) != 1 {
		fmt.Fprintln(stderr, "please input theme id with format: provider/name:version")
		os.Exit(1)
	}

	provider, name, version, err := parse_theme(ctx.Args()[0])
	if err != nil {
		fmt.Fprintln(stderr, "package format: provider/name:version, please try again")
		os.Exit(1)
	}

	// auth config
	config, err := LoadConfigFile(directory)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		os.Exit(1)
	}

	host := ctx.String("Host")
	port := ctx.Int("Port")
	client := NewClient(directory, host, port)

	var pkg Package
	pkg.Provider = provider
	pkg.Name = name
	pkg.Version = version

	if err := client.ThemePull(config.Auth.Token, &pkg); err != nil {
		fmt.Fprintln(stderr, err.Error())
		os.Exit(1)
	}

	bar := pb.New(int(pkg.ArchiveLen)).SetUnits(pb.U_BYTES)
	bar.Prefix(fmt.Sprintf("%s/%s:%s ", pkg.Provider, pkg.Name, pkg.Version))
	bar.Start()
	// create multi writer
	rd := pb.NewPbReader(pkg.ArchiveReader, bar)

	if err := archive.Untar(rd, directory, nil); err != nil {
		fmt.Fprintln(stderr, err.Error())
		os.Exit(1)
	}	
	bar.FinishPrint(fmt.Sprintf("%s/%s:%s pulled succussfully.", pkg.Provider, pkg.Name, pkg.Version))
}