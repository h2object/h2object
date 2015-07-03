package app

import (
	"time"
	"sync"
	"net"
	"net/http"
	"github.com/h2object/stats"
	"github.com/h2object/h2object/log"	
	"github.com/h2object/h2object/util"
	"github.com/h2object/h2object/object"
	"github.com/h2object/h2object/httpext"
	"github.com/h2object/h2object/template"
)

type Application struct{
	sync.RWMutex
	log.Logger

	// option & config
	options		 *Options
	configs		 *CONFIG

	// http
	httpAddr 	 *net.TCPAddr
	httpListener  net.Listener
	// service
	service 	  http.Handler 

	// cache
	cache 		  object.Cache

	// indexes
	pageIndexes   object.Indexes
	objectIndexes object.Indexes

	// store
	systems 	  object.Store
	objects       object.Store
	pages         object.Store

	// template
	templates 	  *template.TemplateLoader

	// stats
	statistics 	  *stats.Stats

	// background workers
	background    util.Background
	exitc   	  chan int
}

func NewApplication(options *Options, configs *CONFIG, logger log.Logger) *Application {
	return &Application{
		Logger: logger,
		options: options,
		configs: configs,
		exitc: make(chan int),
	}
}

func (app *Application) Init() error {
	httpAddr, err := net.ResolveTCPAddr("tcp", app.options.HTTPAddress)
	if err != nil {
		return err
	}
	app.httpAddr = httpAddr

	httpListener, err := net.Listen("tcp", app.httpAddr.String())
	if  err != nil {
		return err
	}
	app.httpListener = httpListener
	
	app.configs.SetSection("h2object")

	cache_expire := app.configs.StringDefault("cache.expire","10m")
	duration_expire, err := time.ParseDuration(cache_expire)
	if err != nil {
		return err
	}

	cache_flush := app.configs.StringDefault("cache.flush", "10s")
	duration_flush, err := time.ParseDuration(cache_flush)
	if err != nil {
		return err
	}
	// cache init
	app.cache = object.NewMutexCache(duration_expire, duration_flush)

	// index init
	app.pageIndexes = object.NewBleveIndexes(app.options.IndexesRoot, "pages.idx")
	app.objectIndexes = object.NewBleveIndexes(app.options.IndexesRoot, "objects.idx")

	// store init
	systems := object.NewBoltStore(app.options.StorageRoot, "systems.dat", object.BoltCoder{})
	if err := systems.Load(); err != nil {
		return err
	}
	app.systems = systems

	objects := object.NewBoltStore(app.options.StorageRoot, "objects.dat", object.BoltCoder{})
	if err := objects.Load(); err != nil {
		return err
	}
	app.objects = objects

	pages := object.NewBoltStore(app.options.StorageRoot, "pages.dat", object.BoltCoder{})
	if err := pages.Load(); err != nil {
		return err
	}
	app.pages = pages
	// template
	paths := []string{app.options.TemplateRoot}
	delimiters := app.configs.StringDefault("template.delimiters","{{ }}")
	app.templates = template.NewTemplateLoader(delimiters, paths, app.Logger)

	// template load
	if err := app.templates.Load(); err != nil {
		return err
	}

	// context set
	index := app.configs.StringDefault("index", "")
	appid := app.configs.StringDefault("appid", "")
	secret := app.configs.StringDefault("secret", "")
	duration := app.configs.DurationDefault("markdown.cache", 10 * time.Minute)

	ctx := new_context(app)
	ctx.index = index

	ctx.signature = util.SignString(secret, appid)
	ctx.markdowns = app.configs.MultiStringDefault("markdown.suffix", ",", []string{"md", "markdown"})
	ctx.templates = app.configs.MultiStringDefault("template.suffix", ",", []string{"html", "htm", "tpl"})
	ctx.cache_duration = duration
	ctx.devmode = app.configs.BoolDefault("develope.mode", false)
	if err := ctx.init(); err != nil {
		return err
	}
	set_context(ctx)

	// service
	serv := new_service(ctx)
	if err := serv.init(); err != nil {
		return err
	}
	
	// stats
	stats_enable := app.configs.BoolDefault("stats.enable", true)
	if stats_enable {
		app.statistics = stats.New()
		app.service = app.statistics.Handler(serv)
	} else {
		app.service = serv
	}

	// init succeed
	return nil
}

func (app *Application) Main() {
	app.background.Work(func() { 
		ext.Serve(app.httpListener, app.service, "http", app.Logger) 
		app.Info("background serving worker exiting")
	})

	app.background.Work(func() { 
		c := time.Tick(app.options.RefreshInterval)
		for {
			select {
			case <- c:
				app.Refresh()
			case <- app.exitc:
				goto timeExit
			}	
		}
	timeExit:
		app.Info("background refresh worker exiting")
	})
}

func (app *Application) Refresh() {
	app.Info("application refresh ...")
	app.templates.Refresh()
}

func (app *Application) Exit() {
	if app.httpListener != nil {
		app.httpListener.Close()
	}

	app.Lock()
	// do something if needed
	app.Unlock()

	// notify app to exit
	close(app.exitc)

	// wait all backgroud workers
	app.background.Wait()
}

