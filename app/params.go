package app

import (
	"path"
	"strconv"
	"net/http"
	"github.com/h2object/h2object/util"
	"github.com/h2object/h2object/httpext"
	"github.com/disintegration/imaging"
)


func params_filter(ctx *context, c *ext.Controller, filters []filter) {
	if done := do_params(ctx, c); done {
		ctx.Info("request (%s) (%s) done by params", c.Request.MethodToLower(), c.Request.URI())
		return
	}
	filters[0](ctx, c, filters[1:])
}

func do_params(ctx *context, ctrl *ext.Controller) bool {
	qr := ctrl.Request.Param("qrcode")
	if qr != "" {
		if size, err := strconv.ParseInt(qr, 10, 64); err == nil {
			return do_qrcode(ctx, ctrl, int(size))
		}
	}

	resize := ctrl.Request.Param("resize")
	if resize != "" {	
		fn := path.Join(ctx.app.Options.StaticRoot, "resize", util.ResizeKey(ctrl.Request.URI(), resize))
		if util.Exist(fn) {
			ctrl.File(fn)
			return true
		}
		if rect, err := util.ParseRect(resize, "x"); err == nil {
			f1 := path.Join(ctx.app.Options.StaticRoot, ctrl.Request.URI())
			ctx.Info("f1 :%s", f1)
			if src, err := imaging.Open(f1); err == nil {
				dst := imaging.Resize(src, rect.Width, rect.Height, imaging.Lanczos)
				if err := imaging.Save(dst, fn); err == nil {
					ctrl.File(fn)
					return true
				} 				
			}
		}
	}

	fit := ctrl.Request.Param("fit")
	if fit != "" {
		fn := path.Join(ctx.app.Options.StaticRoot, "fit", util.ResizeKey(ctrl.Request.URI(), resize))
		if util.Exist(fn) {
			ctrl.File(fn)
			return true
		}
		if rect, err := util.ParseRect(resize, "x"); err == nil {
			f1 := path.Join(ctx.app.Options.StaticRoot, ctrl.Request.URI())
			if src, err := imaging.Open(f1); err == nil {
				dst := imaging.Fit(src, rect.Width, rect.Height, imaging.Lanczos)
				if err := imaging.Save(dst, fn); err == nil {
					ctrl.File(fn)
					return true
				}
			}
		}
	}

	thumbnail := ctrl.Request.Param("thumbnail")
	if thumbnail != "" {
		fn := path.Join(ctx.app.Options.StaticRoot, "thumbnail", util.ResizeKey(ctrl.Request.URI(), resize))
		if util.Exist(fn) {
			ctrl.File(fn)
			return true
		}
		if rect, err := util.ParseRect(thumbnail, "x"); err == nil {
			f1 := path.Join(ctx.app.Options.StaticRoot, ctrl.Request.URI())
			if src, err := imaging.Open(f1); err == nil {
				dst := imaging.Thumbnail(src, rect.Width, rect.Height, imaging.Lanczos)
				if err := imaging.Save(dst, fn); err == nil {
					ctrl.File(fn)
					return true
				}
			}
		}
	}

	return false
}


func do_qrcode(ctx *context, ctrl *ext.Controller, size int) bool {
	file, err := ctx.qrcode_path(ctrl.Request.URL, size)
	if err != nil {
		ctrl.JsonError(http.StatusInternalServerError, err)
		return true
	}
	ctrl.File(file)
	return true
}
