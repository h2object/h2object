package app

import (
	"os"
	"io"
	"fmt"
	"path"
	"errors"
	"net/http"
	"github.com/h2object/h2object/util"
	"github.com/h2object/h2object/third"
	"github.com/h2object/h2object/httpext"
)

func static_filter(ctx *context, c *ext.Controller, filters []filter) {
	if done := do_static(ctx, c); done {
		ctx.Info("request (%s) (%s) done by static", c.Request.MethodToLower(), c.Request.URI())
		return
	}
	filters[0](ctx, c, filters[1:])
}

func do_static(ctx *context, ctrl *ext.Controller) bool {
	switch ctrl.Request.MethodToLower() {
	case "get":
		return do_static_get(ctx, ctrl)
	case "put":
		return do_static_put(ctx, ctrl)
	}
	ctrl.JsonError(http.StatusMethodNotAllowed, errors.New("method not allowed"))
	return true
}

func do_static_get(ctx *context, ctrl *ext.Controller) bool {	
	r := ctrl.Request
	
	if !util.Exist(path.Join(ctx.app.Options.StaticRoot, r.URI())) {
		return false
	}

	ctx.app.Configs.SetSection("third")
	qiniu_enable := ctx.app.Configs.BoolDefault("qiniu.enable", false)
	qiniu_domain := ctx.app.Configs.StringDefault("qiniu.domain", "")

	if qiniu_enable {
		if val, err := ctx.app.systems.Get(path.Join("/qiniu", r.URI()), true); err == nil {
			var key string
			if err := util.Convert(val, &key); err == nil {
				uri := fmt.Sprintf("http://%s", path.Join(qiniu_domain, key))
				ctrl.Redirect(uri, r.URL.RawQuery)
				return true
			}
		}
	}
	ctrl.File(path.Join(ctx.app.Options.StaticRoot, r.URI()))
	return true
}

func do_static_put(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request
	
	dir, file := path.Split(r.URI())

	// create static path if not exist
	realDir := path.Join(ctx.app.Options.StaticRoot, dir)
	if err := os.MkdirAll(realDir, os.ModePerm); err != nil {
		ctrl.JsonError(http.StatusInternalServerError, err)
		return true
	}
	fn := path.Join(ctx.app.Options.StaticRoot, r.URI())

	ctx.app.Configs.SetSection("third")
	qiniu_enable := ctx.app.Configs.BoolDefault("qiniu.enable", false)
	if qiniu_enable {
		
		qiniu_appid := ctx.app.Configs.StringDefault("qiniu.appid", "")
		qiniu_secret := ctx.app.Configs.StringDefault("qiniu.secret", "")
		qiniu_bucket := ctx.app.Configs.StringDefault("qiniu.bucket", "")	
		
		defer func(key string, file string) {
		
			helper := third.NewQiniuHelper(qiniu_appid, qiniu_secret, ctx.app.cache)		
			key, err := helper.PutFile(qiniu_bucket, key, file)
			if err != nil {
				ctx.Warn("request static qiniu put failed: (%s)", err.Error())
				return
			}

			if err := ctx.app.systems.Put(path.Join("/qiniu", r.URI()), key); err != nil {
				ctx.Warn("request static qiniu save failed: (%s)", err.Error())
				return
			}
			ctx.Info("request static qiniu put succeed with key: (%s)", key)
		
		}(third.QiniuKey(r.URI()), fn)
	}

	tmp := fn + ".tmp"
	fd, err := os.Create(tmp)
	if err != nil {
		ctrl.JsonError(http.StatusInternalServerError, err)
		return true
	}	
	defer fd.Close()

	if _, err := io.Copy(fd, r.Body); err != nil {
		ctrl.JsonError(http.StatusInternalServerError, err)
		return true
	}

	if err := os.Rename(tmp, fn); err != nil {
		ctrl.JsonError(http.StatusInternalServerError, err)
		return true
	}

	ctrl.Json(map[string]interface{}{
		"file": file,
	})
	return true
}

