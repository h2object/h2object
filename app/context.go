package app

import (
	"os"
	"time"
	"path"
	"errors"
	"io/ioutil"
	"github.com/h2object/h2object/log"
	"github.com/h2object/h2object/util"
	"github.com/h2object/h2object/page"
)

var (
	default_context *context = nil
)

type context struct{
	log.Logger
	app *Application
	index 	  string
	signature string
	markdowns []string
	templates []string
	cache_duration time.Duration
	devmode   bool
}

func set_context(ctx *context) {
	default_context = ctx
}

func get_context() *context {
	return default_context
}

func new_context(app *Application) *context {
	return &context{
		Logger: app.Logger,
		app: app,
	}
}

func (ctx *context) init() error {
	
	for _, suffix := range ctx.markdowns {
		handlers[suffix] = do_markdown
	}
	for _, suffix := range ctx.templates {
		handlers[suffix] = do_template
	}	
	if ctx.devmode {
		handlers["page"] = do_page
		handlers["system"] = do_page
	}


	if err := ctx.init_pages(); err != nil {
		return err
	}

	return nil
}



func (ctx *context) get_page(uri string) *page.Page {
	if pg, ok := ctx.app.cache.Get(uri); ok {
		defer ctx.app.cache.Set(uri, pg, ctx.cache_duration)
		return pg.(*page.Page)
	}

	pg := page.NewPage(uri)

	if val, err := ctx.app.pages.Get(uri, true); err == nil {
		if err := pg.SetData(val); err == nil {
			defer ctx.app.cache.Set(uri, pg, ctx.cache_duration)
			return pg
		}	
	}

	if err := pg.Load(path.Join(ctx.app.Options.MarkdownRoot, uri)); err == nil {
		defer ctx.app.cache.Set(uri, pg, ctx.cache_duration)
		defer ctx.put_page(uri, pg)
		return pg
	}

	return nil
}

func (ctx *context) put_page(uri string, pg *page.Page) {
	if pg != nil {
		if err := ctx.app.pages.Put(uri, pg.GetData()); err != nil {
			ctx.Warn("context put (%s) page failed:(%s)", uri, err)
		}
		if err := ctx.app.pageIndexes.IndexIfNotExist(uri, pg.GetData()); err != nil {
			ctx.Warn("context index (%s) page failed:(%s)", uri, err)	
		}
	}
}

func (ctx *context) get_pages(uri string) *Pages {
	return NewPages(uri, ctx)
}

func (ctx *context) folder_pages(root, folder string, suffixes []string) error {
	stat, err := os.Stat(folder)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return errors.New(folder + " is not folder")
	}
	if infos, err := ioutil.ReadDir(folder); err == nil {
		for _, info := range infos {
			if info.IsDir() {
				if err := ctx.folder_pages(path.Join(root, info.Name()), 
						  path.Join(folder, info.Name()), suffixes); err != nil {
					return err
				}
				continue
			}

			for _,suffix := range suffixes {
				if !util.HasSuffix(info.Name(), suffix) {
					continue
				}

				uri := path.Join(root, info.Name())
				fn := path.Join(folder, info.Name())
				pg := page.NewPage(uri)
				if err := pg.Load(fn); err != nil {
					return err
				}

				ctx.put_page(uri, pg)
			}
		}
	}
	return nil
}

func (ctx *context) init_pages() error {
	return ctx.folder_pages("/", ctx.app.Options.MarkdownRoot, ctx.markdowns)
}
