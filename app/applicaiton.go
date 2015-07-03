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
	Options		 *Options
	Configs		 *CONFIG

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
		Options: options,
		Configs: configs,
		exitc: make(chan int),
	}
}

func (app *Application) Init() error {
	httpAddr, err := net.ResolveTCPAddr("tcp", app.Options.HTTPAddress)
	if err != nil {
		return err
	}
	app.httpAddr = httpAddr

	httpListener, err := net.Listen("tcp", app.httpAddr.String())
	if  err != nil {
		return err
	}
	app.httpListener = httpListener
	
	app.Configs.SetSection("h2object")

	cache_expire := app.Configs.StringDefault("cache.expire","10m")
	duration_expire, err := time.ParseDuration(cache_expire)
	if err != nil {
		return err
	}

	cache_flush := app.Configs.StringDefault("cache.flush", "10s")
	duration_flush, err := time.ParseDuration(cache_flush)
	if err != nil {
		return err
	}
	// cache init
	app.cache = object.NewMutexCache(duration_expire, duration_flush)

	// index init
	app.pageIndexes = object.NewBleveIndexes(app.Options.IndexesRoot, "pages.idx")
	app.objectIndexes = object.NewBleveIndexes(app.Options.IndexesRoot, "objects.idx")

	// store init
	systems := object.NewBoltStore(app.Options.StorageRoot, "systems.dat", object.BoltCoder{})
	if err := systems.Load(); err != nil {
		return err
	}
	app.systems = systems

	objects := object.NewBoltStore(app.Options.StorageRoot, "objects.dat", object.BoltCoder{})
	if err := objects.Load(); err != nil {
		return err
	}
	app.objects = objects

	pages := object.NewBoltStore(app.Options.StorageRoot, "pages.dat", object.BoltCoder{})
	if err := pages.Load(); err != nil {
		return err
	}
	app.pages = pages
	// template
	paths := []string{app.Options.TemplateRoot}
	delimiters := app.Configs.StringDefault("template.delimiters","{{ }}")
	app.templates = template.NewTemplateLoader(delimiters, paths, app.Logger)

	// template load
	if err := app.templates.Load(); err != nil {
		return err
	}

	// context set
	index := app.Configs.StringDefault("index", "")
	appid := app.Configs.StringDefault("appid", "")
	secret := app.Configs.StringDefault("secret", "")
	duration := app.Configs.DurationDefault("markdown.cache", 10 * time.Minute)

	ctx := new_context(app)
	ctx.index = index

	ctx.signature = util.SignString(secret, appid)
	ctx.markdowns = app.Configs.MultiStringDefault("markdown.suffix", ",", []string{"md", "markdown"})
	ctx.templates = app.Configs.MultiStringDefault("template.suffix", ",", []string{"html", "htm", "tpl"})
	ctx.cache_duration = duration
	ctx.devmode = app.Configs.BoolDefault("develope.mode", false)
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
	stats_enable := app.Configs.BoolDefault("stats.enable", true)
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
		c := time.Tick(app.Options.RefreshInterval)
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

