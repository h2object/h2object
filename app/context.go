package app

import (
	"os"
	"time"
	"path"
	"sync"
	"errors"
	"strconv"
	"net/url"
	"io/ioutil"
	"github.com/h2object/h2object/log"
	"github.com/h2object/h2object/util"
	"github.com/h2object/h2object/page"
	qrcode "github.com/skip2/go-qrcode"
)

var (
	default_context *context = nil
)

type context struct{
	sync.RWMutex
	log.Logger
	app *Application
	site_name string
	site_description string
	site_author string
	site_contact string
	host 	  string
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
	ctx.load()
	return ctx.init_pages()
}

func (ctx *context) load() {	
	ctx.Lock()
	defer ctx.Unlock()

	// clear handlers
	for _, suffix := range ctx.markdowns {
		delete(handlers, suffix)
	}
	for _, suffix := range ctx.templates {
		delete(handlers, suffix)
	}
	delete(handlers, "page")
	delete(handlers, "system")

	// load
	conf := ctx.app.Configs
	conf.SetSection("h2object")

	ctx.site_name = conf.StringDefault("site.name", "")
	ctx.site_description = conf.StringDefault("site.description", "")
	ctx.site_author = conf.StringDefault("site.author", "h2object")
	ctx.site_contact = conf.StringDefault("site.author", "support@h2object.io")

	appid_dft, _ := util.AlphaStringRange(24, 32)
	secret_dft, _ := util.AlphaStringRange(32, 36)
	appid := conf.StringDefault("appid", appid_dft)
	secret := conf.StringDefault("secret", secret_dft)
	ctx.host = conf.StringDefault("host", ctx.app.Options.HTTPAddress)
	ctx.index = conf.StringDefault("index", "")
	ctx.signature = util.SignString(secret, appid)
	ctx.markdowns = conf.MultiStringDefault("markdown.suffix", ",", []string{"md", "markdown"})
	ctx.templates = conf.MultiStringDefault("template.suffix", ",", []string{"html", "htm", "tpl"})
	ctx.cache_duration = conf.DurationDefault("markdown.cache", 10 * time.Minute)
	ctx.devmode = conf.BoolDefault("develope.mode", false)

	// reset handlers
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

func (ctx *context) del_page(uri string) {
	ctx.app.cache.Delete(uri)
	if err := ctx.app.pages.Del(uri); err != nil {
		ctx.Warn("context del (%s) page failed:(%s)", uri, err)
	}
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

type QRCode struct{
	Value string
	File  string
	Link  string
}

func (ctx *context) qrcode_generate(u *url.URL, size int) (*QRCode, error) {
	u2 := ctx.qrcode_url(u)

	var qc QRCode
	qc.Value = u2.String()
	qc.File = path.Join(ctx.app.Options.StaticRoot, "qrcode", util.QrcodeKey(u2.Path, size))

	query := u2.Query()
	query.Set("qrcode", strconv.Itoa(size))

	u3 := u2
	u3.RawQuery = query.Encode()
	qc.Link = u3.String()

	ctx.Info("qrcode value:(%s) link:(%s)", qc.Value, qc.Link)
	if util.Exist(qc.File) {
		return &qc, nil
	}
	if err := qrcode.WriteFile(qc.Value, qrcode.High, size, qc.File); err != nil {
		return nil, err
	}

	return &qc, nil
}


func (ctx *context) qrcode_url(u *url.URL) *url.URL {
	q := u.Query()
	q.Del("qrcode")
	u.RawQuery = q.Encode()
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	if u.Host == "" {
		u.Host = ctx.host	
	}
	return u
}

func (ctx *context) qrcode_path(u *url.URL, size int) (string, error) {
	code, err := ctx.qrcode_generate(u, size)
	if err != nil {
		return "", err
	}
	return code.File, nil
}

func (ctx *context) qrcode_value(u *url.URL, size int) (string, error) {
	code, err := ctx.qrcode_generate(u, size)
	if err != nil {
		return "", err
	}
	return code.Value, nil
}

func (ctx *context) qrcode_link(u *url.URL, size int) (string, error) {
	code, err := ctx.qrcode_generate(u, size)
	if err != nil {
		return "", err
	}
	return code.Link, nil
}
