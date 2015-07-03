package app

import (
	"fmt"
	"sort"
	"strings"
	"github.com/h2object/cast"
	"github.com/blevesearch/bleve"
	"github.com/h2object/h2object/page"
)

type OrderBy struct{
	Field string
	Ascend bool
}

var PagesOffset int64 = 0
var PagesSize int64 = 20

type Pages struct{
	uri 	string
	nested  bool
	ctx 	*context
	queries []bleve.Query
	order 	*OrderBy
	total 	int64
	offset  int64
	size 	int64
	need    bool
	results []*page.Page
}

func NewPages(uri string, ctx *context) *Pages {	
	return &Pages{
		uri: uri,
		nested: true,
		ctx: ctx,
		queries: []bleve.Query{},
		offset: PagesOffset,
		size: PagesSize,
		need: true,
	}
}

func (pages *Pages) clear() {
	pages.total = 0
	pages.results = []*page.Page{}
	pages.need = true
}

func (pages *Pages) Query(json string) *Pages {
	if q, err := bleve.ParseQuery([]byte(json)); err == nil {
		pages.queries = append(pages.queries, q)
	} else {
		pages.ctx.app.Warn("pages query parse failed:%s", err.Error())
	}
	pages.clear()
	return pages
}

func (pages *Pages) QueryDatetimeRange(field, start, end string) *Pages {
	return pages.QueryDatetimeRangeInclusive(field, start, end, false, false)
}
func (pages *Pages) QueryDatetimeRangeInclusive(field, start, end string, startInc bool, endInc bool) *Pages {
	q := bleve.NewDateRangeInclusiveQuery(&start, &end, &startInc, &endInc)
	q.SetField(field)
	pages.queries = append(pages.queries, q)
	pages.clear()
	return pages
}
func (pages *Pages) QueryNumberRange(field string, start, end interface{}) *Pages {
	return pages.QueryNumberRangeInclusive(field, start, end, false, false)
}

func (pages *Pages) QueryNumberRangeInclusive(field string, start, end interface{}, startInc bool, endInc bool) *Pages {
	var fstart *float64 = nil
	var fend *float64 = nil
	if start != nil {
		if f, err := cast.ToFloat64E(start); err == nil {
			fstart = &f
		}
	}
	if end != nil {
		if f, err := cast.ToFloat64E(end); err == nil {
			fend = &f
		}
	}
	q := bleve.NewNumericRangeInclusiveQuery(fstart, fend, &startInc, &endInc)
	q.SetField(field)
	pages.queries = append(pages.queries, q)
	pages.clear()
	return pages
}

func (pages *Pages) QueryRegexp(field string, regex string) *Pages {
	q := bleve.NewRegexpQuery(regex)
	q.SetField(field)
	pages.queries = append(pages.queries, q)
	pages.clear()
	return pages
}

func (pages *Pages) QueryString(field string, str string) *Pages {
	q := bleve.NewQueryStringQuery(str)
	q.SetField(field)
	pages.queries = append(pages.queries, q)	
	pages.clear()
	return pages
}

func (pages *Pages) QueryTerm(field string, term string) *Pages {	
	q := bleve.NewTermQuery(term)
	q.SetField(field)
	pages.queries = append(pages.queries, q)	
	pages.clear()
	return pages
}

func (pages *Pages) QueryTermPrefix(field string, prefix string) *Pages {	
	q := bleve.NewPrefixQuery(prefix)
	q.SetField(field)
	pages.queries = append(pages.queries, q)
	pages.clear()
	return pages
}

func (pages *Pages) QueryMatch(field string, match string) *Pages {	
	q := bleve.NewMatchQuery(match)
	q.SetField(field)
	pages.queries = append(pages.queries, q)	
	pages.clear()
	return pages
}
func (pages *Pages) QueryMatchPhrase(field string, phrase string) *Pages {	
	q := bleve.NewMatchPhraseQuery(phrase)
	q.SetField(field)
	pages.queries = append(pages.queries, q)
	pages.clear()
	return pages
}

func (pages *Pages) OrderBy(field string, asc bool) *Pages {
	pages.order = &OrderBy{
		Field: field,
		Ascend: asc,
	}
	pages.clear()
	return pages
}

func (pages *Pages) Nested(flag bool) *Pages {
	pages.nested = flag
	pages.clear()
	return pages
}

func (pages *Pages) retrieve() error {
	if !pages.need {
		return nil
	}
	pages.need = false
	var pgs []*page.Page
	if len(pages.queries) > 0 {

		query := bleve.NewConjunctionQuery(pages.queries)
		total, uris, err := pages.ctx.app.pageIndexes.Query(pages.uri, query, pages.offset, pages.size)
		if err != nil {
			return err
		}

		for _, uri := range uris {			
			pg := pages.ctx.get_page(uri)
			if pg != nil {
				pgs = append(pgs, pg)	
			}			
		}
		pages.total = total
	} else {
		datas := []interface{}{}
		for _, sfx := range pages.ctx.markdowns {
			vals, err := pages.ctx.app.pages.MultiGet(pages.uri, "." + sfx, pages.nested)
			if err != nil {
				return err
			}	
			datas = append(datas, vals...)
		}
		pages.ctx.Info("retrieve pages datas len:(%d)", len(datas))	

		for _, data := range datas {
			pg := page.NewPage("")
			if err := pg.SetData(data); err != nil {
				continue
			}
			pages.results = append(pages.results, pg)
		}
		pages.total = int64(len(pages.results))
	}

	if pages.order != nil {
		sorter := page.PageSort{
			Field: pages.order.Field,
			Ascend: pages.order.Ascend,
			Pages: pgs,
		}
		sort.Sort(sorter)
	}

	pages.results = append(pages.results, pgs...)
	
	pages.ctx.Info("retrieve pages results len:(%d)", len(pages.results))	
	return nil
}

type PaginationItem struct{
	PageNo 		int
	PageSize 	int
	PageLink 	string
}

func NewPaginationItem(no int, size int, url string) *PaginationItem {
	var u string
	if strings.Index(url, "?") > 0 {
		u = fmt.Sprintf("%s&page=%d&size=%d", url, no, size)
	} else {
		u = fmt.Sprintf("%s?page=%d&size=%d", url, no, size)
	}	
	return &PaginationItem{
		PageNo: no,
		PageSize: size,
		PageLink: u,
	}
}

func (pages *Pages) Pagination(url string, size int) []*PaginationItem {
	pages.retrieve()
	var items []*PaginationItem
	for i := 0; i <= int(pages.total)/size; i++ {
		item := NewPaginationItem(i, size, url)
		items = append(items, item)
	}
	return items
}

func (pages *Pages) range_size(offset, size int64) {
	if len(pages.queries) != 0 {
		if offset != pages.offset || size != pages.size {
			pages.clear()
		}
	}
	pages.offset = offset
	pages.size = size
}

func (pages *Pages) Limit(offset, size int64) []*page.Page {
	pages.range_size(offset, size)	
	pages.retrieve()
	
	pages.ctx.Info("pages (%d:%d) limit (%d:%d)", len(pages.results), pages.total, offset, size)
	// result only store the limit pages
	if len(pages.queries) > 0 {
		return pages.results
	}

	// result store all uri pages, without
	if offset >= pages.total {
		return []*page.Page{}
	}
	if offset + size >= pages.total {
		return pages.results[offset:]
	}
	return pages.results[offset:offset + size]
}

func (pages *Pages) All() []*page.Page {
	pages.retrieve()
	return pages.results
}

func (pages *Pages) Total() int {
	pages.retrieve()
	return int(pages.total)
}










