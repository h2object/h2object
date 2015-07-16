package app

import (
	"sort"
	"github.com/h2object/h2object/util"
)

type Click struct{
	URI 	string `object:"uri"`
	Count 	int64  `object:"count"`
}

func (ctx *context) inc_click(uri string) error {
	ctx.Lock()
	defer ctx.Unlock()

	var click Click
	if val, err := ctx.app.clicks.Get(uri, true); err == nil {
		
		if err := util.Convert(val, &click); err != nil {
			return err
		}
	}
	click.URI = uri
	click.Count++
	if err := ctx.app.clicks.Put(uri, click); err != nil {
		return err
	}

	return nil
}

func (ctx *context) get_click(uri string) (int64, error) {
	ctx.RLock()
	defer ctx.RUnlock()

	val, err := ctx.app.clicks.Get(uri, true)
	if err != nil {
		return 0, err
	}

	var click Click
	if err := util.Convert(val, &click); err != nil {
		return 0, err
	}

	return click.Count, nil
}

type ByAsc []Click

func (a ByAsc) Len() int           { return len(a) }
func (a ByAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAsc) Less(i, j int) bool { return a[i].Count < a[j].Count }


func (ctx *context) clicks_asc(namespace string, n int, suffix string) ([]Click, error) {
	ctx.RLock()
	defer ctx.RUnlock()

	var clicks []Click

	vals, err := ctx.app.clicks.MultiGet(namespace, suffix, true)
	if err != nil {
		return nil, err
	}

	if err := util.Convert(vals, &clicks); err != nil {
		return nil, err
	}

	sort.Sort(ByAsc(clicks))

	if n >= len(clicks) {
		return clicks, nil	
	}
	return clicks[0:n], nil	
}

type ByDesc []Click

func (a ByDesc) Len() int           { return len(a) }
func (a ByDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDesc) Less(i, j int) bool { return a[i].Count > a[j].Count }


func (ctx *context) clicks_desc(namespace string, n int, suffix string) ([]Click, error) {
	ctx.RLock()
	defer ctx.RUnlock()

	var clicks []Click

	vals, err := ctx.app.clicks.MultiGet(namespace, suffix, true)
	if err != nil {
		return nil, err
	}

	if err := util.Convert(vals, &clicks); err != nil {
		return nil, err
	}

	sort.Sort(ByDesc(clicks))

	if n >= len(clicks) {
		return clicks, nil	
	}
	return clicks[0:n], nil	
}
