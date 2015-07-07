package app

import (
	"os"
	"time"
	"fmt"
	"path"
	"path/filepath"
)

type Options struct{
	HTTPAddress 			string
	HTTPSAddress    		string
	StaticRoot 		string
	MarkdownRoot	string
	TemplateRoot 	string
	LogsRoot 				string
	StorageRoot				string
	IndexesRoot 			string
	TempRoot 				string
	Root 					string
	RefreshInterval 		time.Duration
}

func NewOptions(host string, port int) *Options {
	return &Options{
		HTTPAddress: fmt.Sprintf("%s:%d", host, port),
	}
}

func (opt *Options) Prepare(workdir string) error {
	directory, err := filepath.Abs(workdir)
	if err != nil {
		return err
	}
	opt.Root = directory
	opt.StaticRoot = path.Join(directory, "statics")
	opt.MarkdownRoot = path.Join(directory, "markdowns")
	opt.TemplateRoot = path.Join(directory, "templates")
	opt.StorageRoot = path.Join(directory, "storage")
	opt.IndexesRoot = path.Join(directory, "indexes")
	opt.TempRoot = path.Join(directory, ".tmp")
	opt.LogsRoot = path.Join(directory, "logs")
	if err := os.MkdirAll(opt.Root, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(opt.StaticRoot, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(path.Join(opt.StaticRoot, "qrcode"), os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(opt.MarkdownRoot, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(opt.TemplateRoot, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(opt.StorageRoot, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(opt.IndexesRoot, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(opt.LogsRoot, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(opt.TempRoot, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (opt *Options) SetRefreshDefault(s string, default_refresh time.Duration) {
	if d, err := time.ParseDuration(s); err == nil {
		opt.RefreshInterval = d
	} else {
		opt.RefreshInterval = default_refresh
	}
}