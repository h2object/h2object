package app

import (
	"os"
	"sync"
	"time"
	"strings"
	"github.com/h2object/config"
	"github.com/h2object/h2object/util"
)

var default_config *CONFIG

var default_comment string = `
h2object applicaiton settings

section [h2object] is for the applicaiton settings
section [theme]	is for the applicaiton theme push settings
section [deploy] is for the applicaiton deployment server settings
section [third] is for the applicaiton use third plugin settings, like qiniu cloud storage settings
section [logs] is for the file log settings
`

type CONFIG struct{
	sync.RWMutex
	filename string
	section string 
	config  *config.Config
}

func NewCONFIG() *CONFIG {
	return &CONFIG{
		filename: "",
		section: "",
		config: config.New(config.DEFAULT_COMMENT, config.DEFAULT_SEPARATOR, false, true),
	}
}

func DefaultCONFIG() *CONFIG {
	default_config = NewCONFIG()
	default_config.config.AddOption("theme", "provider", "h2object")
	default_config.config.AddOption("theme", "name", "demo")
	default_config.config.AddOption("theme", "catagory", "0")
	default_config.config.AddOption("theme", "version", "0.0.1")
	default_config.config.AddOption("theme", "description", "h2object demo site")
	
	default_config.config.AddOption("deploy", "host", "h2object.io")
	default_config.config.AddOption("deploy", "port", "80")
	default_config.config.AddOption("deploy", "appid", "")
	default_config.config.AddOption("deploy", "secret", "")

	appid, _ := util.AlphaStringRange(24, 32)
	secret, _ := util.AlphaStringRange(32, 36)
	default_config.config.AddOption("h2object", "appid", appid)
	default_config.config.AddOption("h2object", "secret", secret)
	default_config.config.AddOption("h2object", "host", "")
	default_config.config.AddOption("h2object", "index", "")

	default_config.config.AddOption("h2object", "markdown.cache", "10m")
	default_config.config.AddOption("h2object", "markdown.suffix", "md,markdown")
	default_config.config.AddOption("h2object", "template.suffix", "html,htm,tpl")
	default_config.config.AddOption("h2object", "develope.mode", "false")

	default_config.config.AddOption("third", "qiniu.enable", "false")
	default_config.config.AddOption("third", "qiniu.appid", "")
	default_config.config.AddOption("third", "qiniu.secret", "")
	default_config.config.AddOption("third", "qiniu.domain", "")
	default_config.config.AddOption("third", "qiniu.bucket", "")

	default_config.config.AddOption("logs", "file.enable", "false")
	default_config.config.AddOption("logs", "file.name", "h2o.log")
	default_config.config.AddOption("logs", "file.level", "info")
	default_config.config.AddOption("logs", "file.rotate_max_size", "")
	default_config.config.AddOption("logs", "file.rotate_max_line", "")
	default_config.config.AddOption("logs", "file.rotate_daily", "true")

	return default_config
}

func LoadCONFIG(fn string) (*CONFIG, error) {
	var err error
	conf, err := config.ReadDefault(fn)
	if err != nil {
		return nil, err
	}

	config :=  &CONFIG{
		filename: fn,
		section: "h2object",
		config: conf,
	}

	if err := config.Save(""); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *CONFIG) Raw() *config.Config {
	return c.config
}

func (c *CONFIG) SetSection(section string) {
	c.Lock()
	defer c.Unlock()
	c.section = section
}

func (c *CONFIG) SetOption(name, value string) {
	c.Lock()
	defer c.Unlock()
	c.config.AddOption(c.section, name, value)
}

func (c *CONFIG) Int(option string) (result int, found bool) {
	c.RLock()
	defer c.RUnlock()
	result, err := c.config.Int(c.section, option)
	if err == nil {
		return result, true
	}
	if _, ok := err.(config.OptionError); ok {
		return 0, false
	}

	return 0, false
}

func (c *CONFIG) IntDefault(option string, dfault int) int {
	if r, found := c.Int(option); found {
		return r
	}
	return dfault
}

func (c *CONFIG) Float(option string) (result float64, found bool) {
	c.RLock()
	defer c.RUnlock()
	result, err := c.config.Float(c.section, option)
	if err == nil {
		return result, true
	}
	if _, ok := err.(config.OptionError); ok {
		return 0, false
	}

	return 0, false
}

func (c *CONFIG) FloatDefault(option string, dfault float64) float64 {
	if r, found := c.Float(option); found {
		return r
	}
	return dfault
}


func (c *CONFIG) Bool(option string) (result, found bool) {
	c.RLock()
	defer c.RUnlock()
	result, err := c.config.Bool(c.section, option)
	if err == nil {
		return result, true
	}
	if _, ok := err.(config.OptionError); ok {
		return false, false
	}

	return false, false
}

func (c *CONFIG) BoolDefault(option string, dfault bool) bool {
	if r, found := c.Bool(option); found {
		return r
	}
	return dfault
}

func (c *CONFIG) String(option string) (result string, found bool) {
	c.RLock()
	defer c.RUnlock()
	if r, err := c.config.String(c.section, option); err == nil {
		return stripQuotes(r), true
	}
	return "", false
}

func (c *CONFIG) StringDefault(option, dfault string) string {
	if r, found := c.String(option); found {
		s := stripQuotes(r)
		if len(s) != 0 {
			return s
		}
	}
	return dfault
}

func (c *CONFIG) Duration(option string) (result time.Duration, found bool) {
	c.RLock()
	defer c.RUnlock()

	if r, err := c.config.String(c.section, option); err == nil {
		if d, err := time.ParseDuration(stripQuotes(r)); err == nil {
			return d, true
		}
	}
	return time.Duration(0), false
}

func (c *CONFIG) DurationDefault(option string, dfault time.Duration) time.Duration {
	if r, found := c.Duration(option); found {
		return r
	}
	return dfault
}

func (c *CONFIG) MultiString(option string, sep string) ([]string, bool) {
	c.RLock()
	defer c.RUnlock()
	if r, err := c.config.String(c.section, option); err == nil {
		r = stripQuotes(r)
		rs := strings.Split(r, sep)
		return rs, true
	}
	return []string{}, false
}

func (c *CONFIG) MultiStringDefault(option, sep string, dfault []string) []string {
	if r, found := c.MultiString(option, sep); found {
		return r
	}
	return dfault
}

func (c *CONFIG) HasSection(section string) bool {
	c.RLock()
	defer c.RUnlock()
	return c.config.HasSection(section)
}

// Options returns all configuration option keys.
// If a prefix is provided, then that is applied as a filter.
func (c *CONFIG) Options(prefix string) []string {
	c.RLock()
	defer c.RUnlock()
	var options []string
	keys, _ := c.config.Options(c.section)
	for _, key := range keys {
		if strings.HasPrefix(key, prefix) {
			options = append(options, key)
		}
	}
	return options
}

// Helpers
func stripQuotes(s string) string {
	if s == "" {
		return s
	}

	if s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}

	return s
}

func (c *CONFIG) Save(fn string) error {
	c.Lock()
	defer c.Unlock()
	savefile := fn
	if fn == "" {
		savefile = c.filename
	}

	tmp := savefile + ".tmp"
	if err := c.config.WriteFile(tmp, 0644, default_comment); err != nil {
		return err
	}
	os.Remove(savefile)
	return os.Rename(tmp, savefile)
}



