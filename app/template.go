package app

import (	
	"time"
	"strconv"
	"reflect"
	"net/url"
	"html/template"
	"golang.org/x/net/html"
	"github.com/h2object/h2object/page"
	h2object "github.com/h2object/h2object/template"
)

//! page
func t_page(uri string) *page.Page {
	ctx := get_context()
	if ctx != nil {
		return ctx.get_page(uri)
	}
	return nil	
}

func t_title(pg *page.Page) string {
	return pg.Title()
}

func t_uri(pg *page.Page) string {
	if pg != nil {
		return pg.Meta("uri").(string)
	}
	return ""
}

func t_meta(key string, pg *page.Page) interface{} {
	if pg != nil {
		return pg.Meta(key)
	}
	return nil
}

func t_meta_default(key string, dft interface{}, pg *page.Page) interface{} {
	if pg != nil {
		if v := pg.Meta(key); v != nil {
			return v
		}
	}
	return dft
}

func t_modify_datetime(pg *page.Page) time.Time {
	if pg != nil {
		return pg.ModifyDatetime()
	}
	return time.Time{}
}

func t_markdown(pg *page.Page) template.HTML {
	if pg != nil {
		return template.HTML(pg.Markdown())
	}
	return template.HTML("")
}

func t_tag(name string, pg *page.Page) []*html.Node {
	if pg != nil {
		return pg.TagNameAll(name)
	}
	return nil
}

func t_tag_id(id string, pg *page.Page) []*html.Node {
	if pg != nil {
		return pg.TagIDAll(id)
	}
	return nil		
}

func t_tag_class(class string, pg *page.Page) []*html.Node {
	if pg != nil {
		return pg.TagClassAll(class)
	}
	return nil
}

func t_node_text(node *html.Node) string {
	if node != nil {
		return page.NodeText(node)
	}
	return ""
}

func t_node_attr(attr string, node *html.Node) string {
	if node != nil {
		return page.NodeAttribute(node, attr)
	}	
	return ""
}

func t_node_html(node *html.Node) template.HTML {
	if node != nil {
		return page.NodeHtml(node)
	}	
	return template.HTML("")
}

//! pages
func t_pages(uri string) *Pages {
	ctx := get_context()
	if ctx != nil {
		return ctx.get_pages(uri)
	}
	return nil	
}

func t_nested(flag bool, pgs *Pages) *Pages {
	if pgs != nil {
		pgs.Nested(flag)	
	}			
	return pgs
}

func t_query(json string, pgs *Pages) *Pages {
	if pgs != nil {
		pgs.Query(json)
	}
	return pgs
}

func t_query_regexp(field string, regexp string, pgs *Pages) *Pages {
	if pgs != nil {
		pgs.QueryRegexp(field, regexp)
	}
	return pgs
}
func t_query_string(field string, str string, pgs *Pages) *Pages {
	if pgs != nil {
		pgs.QueryString(field, str)
	}
	return pgs
}

func t_query_term(field string, term string, pgs *Pages) *Pages {
	if pgs != nil {
		pgs.QueryTerm(field, term)
	}
	return pgs
}

func t_query_term_prefix(field string, prefix string, pgs *Pages) *Pages {
	if pgs != nil {
		pgs.QueryTermPrefix(field, prefix)
	}
	return pgs
}

func t_query_match(field string, match string, pgs *Pages) *Pages {
	if pgs != nil {
		pgs.QueryMatch(field, match)
	}
	return pgs
}

func t_query_match_phrase(field string, phrase string, pgs *Pages) *Pages {
	if pgs != nil {
		pgs.QueryMatchPhrase(field, phrase)
	}
	return pgs
}

func t_query_num_range(field string, start, end interface{}, pgs *Pages) *Pages {
	if pgs != nil {
		pgs.QueryNumberRange(field, start, end)
	}
	return pgs
}

func t_query_num_range_inclusive(field string, start, end interface{}, startInc bool, endInc bool, pgs *Pages) *Pages {
	if pgs != nil {
		pgs.QueryNumberRangeInclusive(field, start, end, startInc, endInc)
	}
	return pgs
}

func t_query_date_range(field, start, end string, pgs *Pages) *Pages {
	if pgs != nil {
		pgs.QueryDatetimeRange(field, start, end)
	}
	return pgs
}

func t_query_date_range_inclusive(field, start, end string, startInc bool, endInc bool, pgs *Pages) *Pages {
	if pgs != nil {
		pgs.QueryDatetimeRangeInclusive(field, start, end, startInc, endInc)
	}
	return pgs
}

func t_order_by(field string, asc bool, pgs *Pages) *Pages {
	if pgs != nil {
		pgs.OrderBy(field, asc)
	}	
	return pgs
}

func t_pagination(url string, size int, pgs *Pages) []*PaginationItem{
	if pgs != nil {
		return pgs.Pagination(url, size)
	}
	return nil
}

func t_total(pgs *Pages) int {
	if pgs != nil {
		return pgs.Total()
	}	
	return 0
}

func t_limit(offset int, size int, pgs *Pages) []*page.Page {
	if pgs != nil {
		return pgs.Limit(int64(offset), int64(size))
	}
	return nil
}

func t_limit_by_page(page int, size int, pgs *Pages) []*page.Page {
	if pgs != nil {
		return pgs.Limit(int64(page*size), int64(size))
	}
	return nil
}

func t_all(pgs *Pages) []*page.Page {
	if pgs != nil {
		return pgs.All()
	}
	return nil
}

func t_page_no(item *PaginationItem) int {
	if item != nil {
		return item.PageNo
	}
	return 0
}

func t_page_size(item *PaginationItem) int {
	if item != nil {
		return item.PageSize
	}
	return 0
}

func t_page_link(item *PaginationItem) string {
	if item != nil {
		return item.PageLink
	}
	return ""
}

func t_path(u *url.URL) string {
	if u != nil {
		return u.Path
	}
	return ""
}

func t_param(key string, dft interface{}, u *url.URL) interface{} {
	if u != nil {
		val := u.Query().Get(key)
		switch reflect.TypeOf(dft).Kind() {
		case reflect.Bool:
			if b, err := strconv.ParseBool(val); err == nil {
				return b
			}
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				return int(i)
			}
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			if i, err := strconv.ParseUint(val, 10, 64); err == nil {
				return int(i)
			}
		case reflect.Float32:
			if f, err := strconv.ParseFloat(val, 32); err == nil {
				return f
			}
		case reflect.Float64:
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				return f
			}
		case reflect.String:
			return val
		}
	}
	return dft
}

func init() {
	// url
	h2object.Function("path", t_path)
	h2object.Function("param", t_param)

	// single page template method
	h2object.Function("page", t_page)
	h2object.Function("uri", t_uri)
	h2object.Function("title", t_title)
	h2object.Function("meta", t_meta)
	h2object.Function("meta_default", t_meta_default)
	h2object.Function("modify_datetime", t_modify_datetime)
	h2object.Function("markdown", t_markdown)
	h2object.Function("tag", t_tag)
	h2object.Function("tag_id", t_tag_id)
	h2object.Function("tag_class", t_tag_class)
	h2object.Function("node_text", t_node_text)
	h2object.Function("node_attr", t_node_attr)
	h2object.Function("node_html", t_node_html)

	// multi pages template method
	h2object.Function("pages", t_pages)
	h2object.Function("nested", t_nested)
	h2object.Function("query", t_query)
	h2object.Function("query_regexp", t_query_regexp)
	h2object.Function("query_string", t_query_string)
	h2object.Function("query_term", t_query_term)
	h2object.Function("query_term_prefix", t_query_term_prefix)
	h2object.Function("query_match", t_query_match)
	h2object.Function("query_match_phrase", t_query_match_phrase)
	h2object.Function("query_num_range", t_query_num_range)
	h2object.Function("query_num_range_inclusive", t_query_num_range_inclusive)
	h2object.Function("query_date_range", t_query_date_range)
	h2object.Function("query_date_range_inclusive", t_query_date_range_inclusive)
	h2object.Function("order_by", t_order_by)
	h2object.Function("pagination", t_pagination)
	h2object.Function("total", t_total)
	h2object.Function("limit", t_limit)
	h2object.Function("limit_by_page", t_limit_by_page)
	h2object.Function("all", t_all)

	// pagination
	h2object.Function("page_no", t_page_no)
	h2object.Function("page_size", t_page_size)
	h2object.Function("page_link", t_page_link)	
	
}

		
		