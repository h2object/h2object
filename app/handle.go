package app

import (
	"os"
	"io"
	"fmt"
	"path"
	"time"
	"errors"
	"strings"
	"net/http"
	"github.com/h2object/feeds"
	"github.com/h2object/h2object/util"
	"github.com/h2object/h2object/httpext"
	"github.com/docker/docker/pkg/archive"
)

func do_rss_xml(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request
	dir,_ := path.Split(r.URI())

	pages := NewPages(dir, ctx)
	pgs := pages.All()


	now := time.Now()
	feed := &feeds.Feed{
	    Title:       ctx.site_name,
	    Link:        &feeds.Link{Href: path.Join(ctx.host, dir)},
	    Description: ctx.site_description,
	    Author:      &feeds.Author{ctx.site_author, ctx.site_contact},
	    Created:     now,
	}
	for _, page := range pgs {
		var item feeds.Item
		item.Title = page.Title()
		item.Link = &feeds.Link{Href:fmt.Sprintf("http://%s%s", ctx.host, page.URI())}
		item.Description = page.Description()
		item.Author = feed.Author
		if page.Author() != "" {
			item.Author.Name = page.Author()
		}
		if page.Contact() != "" {
			item.Author.Email = page.Contact()
		}
		item.Created = page.PublishedDatetime()
		feed.Add(&item)
	}

	rss := &feeds.Rss{feed}

	ctrl.Xml(rss)
	return true
}

func do_atom_xml(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request
	dir,_ := path.Split(r.URI())

	pages := NewPages(dir,ctx)
	pgs := pages.All()


	now := time.Now()
	feed := &feeds.Feed{
	    Title:       ctx.site_name,
	    Link:        &feeds.Link{Href: path.Join(ctx.host, dir)},
	    Description: ctx.site_description,
	    Author:      &feeds.Author{ctx.site_author, ctx.site_contact},
	    Created:     now,
	}
	for _, page := range pgs {
		var item feeds.Item
		item.Title = page.Title()
		item.Link = &feeds.Link{Href:fmt.Sprintf("http://%s%s", ctx.host, page.URI())}
		item.Description = page.Description()
		item.Author = feed.Author
		if page.Author() != "" {
			item.Author.Name = page.Author()
		}
		if page.Contact() != "" {
			item.Author.Email = page.Contact()
		}
		item.Created = page.PublishedDatetime()
		feed.Add(&item)
	}

	atom := &feeds.Atom{feed}

	ctrl.Xml(atom)
	return true
}


func do_xml(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request
	if r.MethodToLower() != "get" {
		ctrl.JsonError(http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return true
	}
	
	_ ,file := path.Split(r.URI())
	switch strings.ToLower(file) {
	case "atom.xml":
		return do_atom_xml(ctx, ctrl)
	case "rss.xml":
		return do_rss_xml(ctx, ctrl)
	} 
	ctrl.JsonError(http.StatusMethodNotAllowed, errors.New("request xml not allowed"))
	return true
}

func do_markdown(ctx *context, ctrl *ext.Controller) bool {
	switch ctrl.Request.MethodToLower() {
	case "get":
		return do_markdown_get(ctx, ctrl)
	case "put":
		return do_markdown_put(ctx, ctrl)
	case "delete":
		return do_markdown_delete(ctx, ctrl)
	}
	ctrl.JsonError(http.StatusMethodNotAllowed, errors.New("method not allowed"))
	return true
}

func do_markdown_get(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request

	if !util.Exist(path.Join(ctx.app.Options.MarkdownRoot, r.URI())) {
		return false
	}

	if pg := ctx.get_page(r.URI()); pg != nil {
		data := map[string]interface{}{
			"url": r.URL,
			"page": pg,
		}

		tname := r.Param("template")
		if tname == "" {
			tname = pg.Template()
		}
		if tname == "" {
			ctx.app.Configs.SetSection("h2object")
			tname = ctx.app.Configs.StringDefault("markdown.template", "")
		}

		if tname != "" {
			if t, err := ctx.app.templates.Template(tname); err == nil {
				ctrl.Template(t, data)
				return true
			}
		}
		ctrl.Html(pg.Markdown())
		return true
	}
	return false
}

func do_markdown_delete(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request

	if !util.Exist(path.Join(ctx.app.Options.MarkdownRoot, r.URI())) {
		return false
	}

	os.Remove(path.Join(ctx.app.Options.MarkdownRoot, r.URI()))
	ctx.del_page(r.URI())

	ctrl.Json(map[string]interface{}{
		r.URI(): "deleted",
	})
	return true
}
func do_markdown_put(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request
	dir, file := path.Split(r.URI())
	
	// create markdown dir if not exist
	realDir := path.Join(ctx.app.Options.MarkdownRoot, dir)
	if err := os.MkdirAll(realDir, os.ModePerm); err != nil {
		ctrl.JsonError(http.StatusInternalServerError, err)
		return true
	}

	// local filesystem static file processing
	fn := path.Join(ctx.app.Options.MarkdownRoot, r.URI())
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

	// for cache
	defer func() {
		ctx.del_page(r.URI())
		ctx.get_page(r.URI())
	}()

	ctrl.Json(map[string]interface{}{
			"markdown": file,
	})
	return true
}

func do_template(ctx *context, ctrl *ext.Controller) bool {
	switch ctrl.Request.MethodToLower() {
	case "get":
		return do_template_get(ctx, ctrl)
	case "put":
		return do_template_put(ctx, ctrl)
	case "delete":
		return do_template_delete(ctx, ctrl)
	}
	ctrl.JsonError(http.StatusMethodNotAllowed, errors.New("method not allowed"))
	return true
}

func do_template_get(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request
	if !util.Exist(path.Join(ctx.app.Options.TemplateRoot, r.URI())) {
		return false
	}

	data := map[string]interface{}{
		"url": r.URL,
	}

	if t, err := ctx.app.templates.Template(r.URI()); err == nil {
		ctx.Info("template get (%s) ok", r.URI())
		ctrl.Template(t, data)	
	} else {
		ctrl.JsonError(http.StatusInternalServerError, err)
	}	
	return true
}
func do_template_delete(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request

	if !util.Exist(path.Join(ctx.app.Options.TemplateRoot, r.URI())) {
		return false
	}

	os.Remove(path.Join(ctx.app.Options.TemplateRoot, r.URI()))
	defer ctx.app.templates.Refresh()
	
	ctrl.Json(map[string]interface{}{
		r.URI(): "deleted",
	})
	return true
}
func do_template_put(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request
	dir, file := path.Split(r.URI())
	
	// create markdown dir if not exist
	realDir := path.Join(ctx.app.Options.TemplateRoot, dir)
	if err := os.MkdirAll(realDir, os.ModePerm); err != nil {
		ctrl.JsonError(http.StatusInternalServerError, err)
		return true
	}

	// local filesystem static file processing
	fn := path.Join(ctx.app.Options.TemplateRoot, r.URI())
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

	defer ctx.app.templates.Refresh()
	
	ctrl.Json(map[string]interface{}{
			"template": file,
	})
	return true
}

func do_configure(ctx *context, ctrl *ext.Controller) bool {
	switch ctrl.Request.MethodToLower() {
	case "get":
		return do_configure_get(ctx, ctrl)
	case "put":
		return do_configure_put(ctx, ctrl)
	}
	ctrl.JsonError(http.StatusMethodNotAllowed, errors.New("method not allowed"))
	return true
}

func do_configure_get(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request
	section := strings.Trim(r.TrimSuffixURI(".conf"), "/")	
	ctx.Trace("conf get section (%s)", section)

	if ctx.app.Configs.HasSection(section) == false {		
		ctrl.JsonError(http.StatusNotImplemented, errors.New("section not exist:" + section))
		return true
	}

	ctx.app.Configs.SetSection(section)

	data := map[string]string{}
	fields := r.Params("field")
	if len(fields) == 0 {
		options := ctx.app.Configs.Options("")
		for _, opt := range options {
			data[opt] = ctx.app.Configs.StringDefault(opt, "")
		}
	}

	for _, field := range fields {
		data[field] = ctx.app.Configs.StringDefault(field, "")
	}

	ctrl.Json(data)
	return true
}

func do_configure_put(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request
	section := strings.Trim(r.TrimSuffixURI(".conf"), "/")	
	ctx.Trace("conf put section (%s)", section)

	if ctx.app.Configs.HasSection(section) == false {		
		ctrl.JsonError(http.StatusNotImplemented, errors.New("section not exist:" + section))
		return true
	}

	ctx.app.Configs.SetSection(section)

	data := map[string]interface{}{}
	if err := r.JsonData(&data); err != nil {
		ctrl.JsonError(http.StatusNotImplemented, err)
		return true
	}

	for k, v := range data {
		if str, ok := v.(string); ok {
			ctx.app.Configs.SetOption(k, str)
		} else {
			ctrl.JsonError(http.StatusBadRequest, errors.New("conf value must be string type"))
			return true
		}
	}

	defer ctx.load()
	defer ctx.app.Configs.Save("")

	ctrl.Json(map[string]interface{}{
		"section": section,
	})
	return true
}

func do_export(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request
	if r.MethodToLower() != "get" {
		ctrl.JsonError(http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return true
	}

	directory := path.Join(ctx.app.Options.Root, r.TrimSuffixURI(".export"))
	var tarOpt archive.TarOptions
	tarOpt.ExcludePatterns = append(tarOpt.ExcludePatterns, 
		".tmp", ".h2object", "h2object.pid", "h2object.conf")
	tarOpt.Compression = archive.Gzip
	rd, err := archive.TarWithOptions(directory, &tarOpt)
	if err != nil {
		ctrl.JsonError(http.StatusInternalServerError, err)
		return true		
	}
	defer rd.Close()

	_, fname := path.Split(r.TrimSuffixURI(".export"))	
	if fname == "" {
		fname = "h2object"
	}
	fname = fname + ".tar.gz"
	fn := path.Join(ctx.app.Options.TempRoot, fname)

	fd, err := os.Create(fn)
	if err != nil {
		ctrl.JsonError(http.StatusInternalServerError, err)
		return true		
	}
	if _, err := io.Copy(fd, rd); err != nil {
		ctrl.JsonError(http.StatusInternalServerError, err)
		return true		
	}
	fd.Close()

	fout, err := os.Open(fn)
	if err != nil {
		ctrl.JsonError(http.StatusInternalServerError, err)
		return true		
	}
	fstat, err := fout.Stat()
	if err != nil {
		ctrl.JsonError(http.StatusInternalServerError, err)
		return true		
	}
	
	ctrl.Binary(fname, fout, fstat.Size())
	return true
}

func do_page(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request
	if r.MethodToLower() != "get" {
		ctrl.JsonError(http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return true
	}

	val, err := ctx.app.pages.Get(r.TrimSuffixURI(".page"), true)
	if err != nil {
		ctrl.JsonError(http.StatusInternalServerError, err)
		return true	
	}

	ctrl.Json(val)
	return true
}

func do_click(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request
	if r.MethodToLower() != "get" {
		ctrl.JsonError(http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return true
	}

	val, err := ctx.app.clicks.Get(r.TrimSuffixURI(".click"), true)
	if err != nil {
		ctrl.JsonError(http.StatusInternalServerError, err)
		return true	
	}

	ctrl.Json(val)
	return true
}

func do_system(ctx *context, ctrl *ext.Controller) bool {
	r := ctrl.Request
	if r.MethodToLower() != "get" {
		ctrl.JsonError(http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return true
	}

	val, err := ctx.app.systems.Get(r.TrimSuffixURI(".system"), true)
	if err != nil {
		ctrl.JsonError(http.StatusInternalServerError, err)
		return true	
	}

	ctrl.Json(val)
	return true
}

func do_json(ctx *context, ctrl *ext.Controller) bool {
	return true
}

func do_size(ctx *context, ctrl *ext.Controller) bool {
	return true
}

func do_event(ctx *context, ctrl *ext.Controller) bool {
	return true
}
