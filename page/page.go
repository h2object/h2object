package page

import (
	"os"
	"path"
	"bytes"
	"time"
	"html/template"
	"github.com/h2object/cast"
	"github.com/h2object/markdown-pager"
	"github.com/russross/blackfriday"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"github.com/yhat/scrape"
)

type Page struct{
	uri 		string
	path 		string
	meta    	map[string]interface{}
	content 	string
	node  		*html.Node	
	modify_datetime time.Time	
}

func NewPage(uri string) *Page {
	return &Page{
		uri: uri,
		meta: make(map[string]interface{}),
		modify_datetime: time.Now(),
	}
} 

func NewPageWithData(uri string, data interface{}) (*Page, error) {
	page := NewPage(uri)
	err := page.SetData(data)
	return page, err
}

func (page *Page) URI() string {
	return page.uri
}

func (page *Page) GetData() interface{} {
	return page.meta
}

func (page *Page) Meta(key string) interface{} {
	if v, ok := page.meta[key]; ok {
		return v
	}
	return nil
}

func (page *Page) SetData(data interface{}) error {
	if data != nil {
		meta, err := cast.ToStringMapE(data)
		if err != nil {
			return err
		}

		for k, v := range meta {
			if k == "uri" {
				page.uri = v.(string)
				page.meta[k] = v
				continue
			}
			if k == "path" {
				page.path = v.(string)
				page.meta[k] = v
				continue
			}
			if k == "content" {
				page.content = v.(string)
				page.meta[k] = v
				continue
			}
			if k == "modify_datetime" {
				page.modify_datetime = cast.ToTime(v)
				page.meta[k] = v
				continue
			}

			if t, err := cast.ToTimeE(v); err == nil {
				page.meta[k] = t
				continue
			}
			if f, err := cast.ToFloat64E(v); err == nil {
				page.meta[k] = f
				continue
			}
			page.meta[k] = v
		}	
	}
	return nil
}

func (page *Page) Load(path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	page.path = path

	p, err := pager.ReadFrom(fd)
	if err != nil {
		return err
	}

	metadata, err := p.Metadata()
	if err != nil {
		return err
	}

	if metadata != nil {
		meta, err := cast.ToStringMapE(metadata)
		if err != nil {
			return err
		}
		for k, v := range meta {
			if t, err := cast.ToTimeE(v); err == nil {
				page.meta[k] = t
				continue
			}
			if f, err := cast.ToFloat64E(v); err == nil {
				page.meta[k] = f
				continue
			}
			page.meta[k] = v
		}		
	}
	page.content = string(p.Content())
	page.meta["uri"] = page.uri
	page.meta["path"] = page.path
	page.meta["content"] = page.content
	page.meta["modify_datetime"] = page.modify_datetime
	return nil
}

func (page *Page) Title() string {
	if title, ok := page.meta["title"]; ok {
		return title.(string)
	}

	if n := page.TagName("h1"); n != nil {
		return NodeText(n)	
	}

	if n := page.TagName("h2"); n != nil {
		return NodeText(n)	
	}

	if n := page.TagName("h3"); n != nil {
		return NodeText(n)	
	}

	if n := page.TagName("h4"); n != nil {
		return NodeText(n)	
	}

	if n := page.TagName("h5"); n != nil {
		return NodeText(n)	
	}

	_, file := path.Split(page.uri)
	return file
}

func (page *Page) PublishedDatetime() time.Time {
	if published, ok := page.meta["published_datetime"]; ok {
		if t, err := cast.ToTimeE(published); err == nil {
			return t
		}
	}
	return page.modify_datetime
}

func (page *Page) Summary(max int) string {
	if summary, ok := page.meta["summary"]; ok {
		return summary.(string)
	}

	left := max
	result := ""
	for _, node := range page.TagNameAll("p") {
		text := NodeText(node)
		for _, ch := range text {
			if left > 0 {
				left --
				result = result + string(ch)
			}
		}
	}
	return result
}

func (page *Page) Template() string {
	if tmpl, ok := page.meta["template"]; ok {
		return tmpl.(string)
	}
	return ""
}

func (page *Page) ModifyDatetime() time.Time {
	return page.modify_datetime
}

func (page *Page) Body() *html.Node {
	if page.node == nil {
		buffer := bytes.NewBuffer([]byte(page.Markdown()))
		if body, err := html.Parse(buffer); err == nil {
			page.node = body
			return body
		}
	}

	return page.node
}

func (page *Page) Markdown() string {
	return string(blackfriday.MarkdownCommon([]byte(page.content)))	
}

func (page *Page) TagName(name string) *html.Node {
	tag := atom.Lookup([]byte(name))
	if tag == 0 {
		return nil
	}

	title, ok := scrape.Find(page.Body(), scrape.ByTag(tag))
	if ok {
		return title
	}
	return nil
}

func (page *Page) TagNameAll(name string) []*html.Node {
	var result []*html.Node
	tag := atom.Lookup([]byte(name))
	if tag == 0 {
		return result
	}

	return scrape.FindAll(page.Body(), scrape.ByTag(tag))
}

func (page *Page) TagID(id string) *html.Node {
	node, ok := scrape.Find(page.Body(), scrape.ById(id))
	if ok {
		return node
	}
	return nil
}

func (page *Page) TagIDAll(id string) []*html.Node {
	return scrape.FindAll(page.Body(), scrape.ById(id))
}

func (page *Page) TagClass(class string) *html.Node {
	node, ok := scrape.Find(page.Body(), scrape.ByClass(class))
	if ok {
		return node
	}
	return nil
}

func (page *Page) TagClassAll(class string) []*html.Node {
	return scrape.FindAll(page.Body(), scrape.ByClass(class))
}

func NodeText(node *html.Node) string {
	return scrape.Text(node)
}

func NodeAttribute(node *html.Node, attr string) string {
	return scrape.Attr(node, attr)
}

func NodeHtml(node *html.Node) template.HTML {
	buffer := bytes.NewBuffer([]byte(""))	
	if err := html.Render(buffer, node); err == nil {
		return template.HTML(buffer.String())
	}
	return template.HTML("")
}

