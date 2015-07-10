package commands

import (
	"github.com/codegangsta/cli"
	"path"
	"path/filepath"
	"fmt"
	"os"
	"io/ioutil"
	"strings"
	"github.com/h2object/pb"
	"github.com/h2object/content-type"
	"github.com/h2object/h2object/app"
	"github.com/h2object/h2object/api"
	"github.com/h2object/h2object/util"
	"github.com/docker/docker/pkg/archive"
)

func uri(workdir, fn string) string {
	u := fn
	u = strings.TrimPrefix(u, path.Join(workdir, "markdowns"))
	u = strings.TrimPrefix(u, path.Join(workdir, "templates"))
	u = strings.TrimPrefix(u, path.Join(workdir, "statics"))
	return u
}

func push(client *api.Client, auth api.Auth, bar *pb.ProgressBar, workdir string, dir string, exclude_suffixes []string) error {
	stat, err := os.Stat(dir)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		infos, err := ioutil.ReadDir(path.Join(dir))
		if err != nil {
			return err
		}

		for _, info := range infos {
			if info.IsDir() {
				if err := push(client, auth, bar,
							   workdir, 
							   path.Join(dir, info.Name()), exclude_suffixes); err != nil {
					return err
				}
				continue
			}

			hit := false
			for _, suffix := range exclude_suffixes {
				if util.HasSuffix(info.Name(), suffix) {
					hit = true
					break
				}
			}
			if hit == true {
				continue
			}

			u := uri(workdir, path.Join(dir, info.Name()))
			if err := file_push(client, auth, bar, path.Join(dir, info.Name()), u); err != nil {
				return err
			}
		}
	} else {
		u := uri(workdir, dir)
		if err := file_push(client, auth, bar, dir, u); err != nil {
			return err
		}		
	}	
	return nil
}

func file_push(client *api.Client, auth api.Auth, bar *pb.ProgressBar, fn string, uri string) error {
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer file.Close()
	
	st, err := file.Stat()
	if err != nil {
		return err
	}

	rd := pb.NewPbReader(file, bar)

	_, filename := path.Split(fn)
	contentType := content_type.DefaultContentTypeHelper.ContentTypeByFilename(filename)

	return client.Upload(nil, auth, uri, contentType, rd, st.Size())
}

func deployPullCommand(ctx *cli.Context) {
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

	stderr := os.Stderr
	
	h2oconf, err := app.LoadCONFIG(path.Join(directory, "h2object.conf"))
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		os.Exit(1)
	}

	h2oconf.SetSection("deploy")
	Host := h2oconf.StringDefault("host", "")
	Port := h2oconf.IntDefault("port", 80)
	AppID := h2oconf.StringDefault("appid", "")
	Secret := h2oconf.StringDefault("secret", "")

	client := api.NewClient(Host, Port)
	auth := api.NewAdminAuth(AppID, Secret)

	dirs := ctx.Args()
	if len(dirs) == 0 {
		body, size, err := client.Download(nil, auth, path.Join("/", ".export"))
		if err != nil {
			fmt.Fprintln(stderr, err.Error())
			os.Exit(1)
		}

		bar := pb.New(int(size)).SetUnits(pb.U_BYTES)
		bar.Prefix("/ ")
		bar.Start()
		// create multi writer
		rd := pb.NewPbReader(body, bar)
		if err := archive.Untar(rd, directory, nil); err != nil {
			fmt.Fprintln(stderr, err.Error())
			os.Exit(1)
		}	
		bar.FinishPrint(fmt.Sprintf("/ pulled succussfully without <h2object.conf> file."))
	} else {
		for _, dir := range dirs {
			if !strings.HasPrefix(dir, "markdowns") &&
			   !strings.HasPrefix(dir, "templates") &&
			   !strings.HasPrefix(dir, "statics") &&
			   !strings.HasPrefix(dir, "storage") &&
			   !strings.HasPrefix(dir, "indexes") {
				fmt.Fprintf(stderr, "push path ignored: %s\n", dir)	
				continue
			}
			
			body, size, err := client.Download(nil, auth, path.Join("/", dir + ".export"))
			if err != nil {
				fmt.Fprintln(stderr, err.Error())
				os.Exit(1)
			}

			bar := pb.New(int(size)).SetUnits(pb.U_BYTES)
			bar.Prefix(dir + " ")
			bar.Start()
			// create multi writer
			rd := pb.NewPbReader(body, bar)
			if err := archive.Untar(rd, path.Join(directory, dir), nil); err != nil {
				fmt.Fprintln(stderr, err.Error())
				os.Exit(1)
			}	
			bar.FinishPrint(fmt.Sprintf("%s pulled succussfully.", dir))			
		}
	}
}

func deployPushCommand(ctx *cli.Context) {
	workdir := ctx.GlobalString("workdir")
	if workdir == "" {
		fmt.Println("unknown working directory, please use -w to provide.")
		os.Exit(1)
	}

	// directory
	absworkdir, err := filepath.Abs(workdir)
	if err != nil {
		fmt.Println("workdir:", err)
		return
	}

	dirs := ctx.Args()
	if len(dirs) == 0 {
		dirs = append(dirs, "h2object.conf:h2object:third", "markdowns", "templates", "statics")
	}

	h2oconf, err := app.LoadCONFIG(path.Join(absworkdir, "h2object.conf"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	h2oconf.SetSection("deploy")
	Host := h2oconf.StringDefault("host", "")
	Port := h2oconf.IntDefault("port", 80)
	AppID := h2oconf.StringDefault("appid", "")
	Secret := h2oconf.StringDefault("secret", "")

	h2oconf.SetSection("h2object")
	markdown_suffixs := h2oconf.MultiStringDefault("markdown.suffix", ",", []string{"md", "markdown"})

	client := api.NewClient(Host, Port)
	auth := api.NewAdminAuth(AppID, Secret)
	
	for _, directory := range dirs {
		if strings.HasPrefix(directory, "markdowns") {

			size := util.FolderSize(path.Join(absworkdir, directory), []string{})

			bar := pb.New(int(size)).SetUnits(pb.U_BYTES)
			bar.Prefix(directory + " ")
			bar.Start()

			if err := push(client, auth, bar, absworkdir, path.Join(absworkdir, directory), []string{}); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			bar.FinishPrint(fmt.Sprintf("%s push completed.", directory))
			continue
		}

		if strings.HasPrefix(directory, "templates") {

			size := util.FolderSize(path.Join(absworkdir, directory), []string{})

			bar := pb.New(int(size)).SetUnits(pb.U_BYTES)

			bar.Prefix(directory + " ")
			bar.Start()

			if err := push(client, auth, bar, absworkdir, path.Join(absworkdir, directory), []string{}); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			bar.FinishPrint(fmt.Sprintf("%s push completed.", directory))
			continue
		}

		if strings.HasPrefix(directory, "statics") {
			size := util.FolderSize(path.Join(absworkdir, directory), markdown_suffixs)
			bar := pb.New(int(size)).SetUnits(pb.U_BYTES)
			bar.Prefix(directory + " ")
			bar.Start()

			if err := push(client, auth, bar, absworkdir, path.Join(absworkdir, directory), markdown_suffixs); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			bar.FinishPrint(fmt.Sprintf("%s push completed.", directory))
			continue
		}

		if strings.HasPrefix(directory, "h2object.conf") {
			ds := strings.Split(directory, ":")
			
			sections := []string{}
			if len(ds) > 1 {
				sections = append(sections, ds[1:]...)
			} else {
				sections = append(sections, "h2object", "third")
			}

			for _, section := range sections {
				if h2oconf.HasSection(section) {
					h2oconf.SetSection(section)
					opts := h2oconf.Options("")

					data := map[string]string{}
					for _, opt := range opts {
						if section == "h2object" && (opt == "appid" || opt == "secret") {
							continue
						}
						data[opt] = h2oconf.StringDefault(opt, "")
					}

					if err := client.SetConfig(nil, auth, section, data); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					fmt.Printf("section (%s) push succeed.\n", section)
				}
			}
			continue
		}

		fmt.Printf("push path ignored: %s\n", directory)
	}
}
