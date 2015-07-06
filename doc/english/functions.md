functions
=========

[Back](https://github.com/h2object/h2object/blob/master/doc/english/index.md) 

As we known, if you want generate the markdown meta or content to the template, you need use the template functions to control the content output. the template can use three catagory data from the h2object service.
They are the following:

 * **render variables**
 * **page function to get a single page**
 * **pages function to get the pages' array**

#### Render Variables

<table class="table">
	<thead>
		<th>variable name</th>
		<th>description</th>
	</thead>
	<tr>
		<td>.url</td>
		<td>current page uri</td>
	</tr>
	<tr>
		<td>.page</td>
		<td>current page, need be access by markdown suffix directly</td>
	</tr>
</table>

#### Page URL functions

<table class="table">
	<thead>
		<th>function</th>
		<th>description</th>
		<th>example</th>
	</thead>
	<tr>
		<td>path(*url)</td>
		<td>get the url path</td>
		<td>{{ .url | path}}</td>
	</tr>
	<tr>
		<td>param(name, default_value, *url)</td>
		<td>get the url param</td>
		<td>{{ .url | param "size" 20}}</td>
	</tr>
</table>


#### Single Page functions

<table class="table">
	<thead>
		<th>function</th>
		<th>description</th>
		<th>example</th>
	</thead>
	<tr>
		<td>func page(uri) *page</td>
		<td>get page by uri</td>
		<td>{{page "/demo.md"}}</td>
	</tr>
	<tr>
		<td>func title(*page)</td>
		<td>get title</td>
		<td>{{page "/demo.md" | title}}</td>
	</tr>
	<tr>
		<td>func meta(fieldname, *page) </td>
		<td>get meta</td>
		<td>{{page "/demo.md" | meta "author"}}</td>
	</tr>
	<tr>
		<td>func meta_default(fieldname, default_value, *page) </td>
		<td>get meta with default</td>
		<td>{{page "/demo.md" | meta "author" "james"}}</td>
	</tr>
	<tr>
		<td>func markdown(*page) </td>
		<td>get markdown content</td>
		<td>{{page "/demo.md" | markdown}}</td>
	</tr>
	<tr>
		<td>func tag(name, *page) []*html.node </td>
		<td>get name tag array</td>
		<td>{{page "/demo.md" | tag "div"}}</td>
	</tr>
	<tr>
		<td>func tag_id(id, *page) []*html.node </td>
		<td>get id tag array</td>
		<td>{{page "/demo.md" | tag_id "field_id"}}</td>
	</tr>
	<tr>
		<td>func tag_class(id, *page) []*html.node </td>
		<td>get class tag array</td>
		<td>{{page "/demo.md" | tag_id "css_class"}}</td>
	</tr>
	<tr>
		<td>func node_text(*html.node) </td>
		<td>get text of tag node</td>
		<td>{{.page | tag_id "div"}}<br>
				{{. | node_text}}<br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>func node_html(*html.node) </td>
		<td>get html of tag node</td>
		<td>{{.page | tag_id "div"}}<br>
				{{. | node_html}}<br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>func node_attr(*html.node) </td>
		<td>get html of tag node</td>
		<td>{{.page | tag_id "div"}}<br>
				{{. | node_attr "class"}}<br>
			{{end}}<br>
		</td>
	</tr>
</table>

#### Muiltiple Pages functions

<table class="table">
	<thead>
		<th>function</th>
		<th>description</th>
		<th>example</th>
	</thead>
	<tr>
		<td>pages(uri) *pages</td>
		<td>get pages with prefix uri</td>
		<td>{{with pages "/"}}<br>
			doing somthing ... <br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>nested(bool, *pages) *pages</td>
		<td>set pages nested flag</td>
		<td>{{with pages "/" | nested true}}<br>
			doing somthing ... <br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>query_regexp(regexp, *pages) *pages</td>
		<td>query pages by regexp</td>
		<td>{{with pages "/" | query_regexp "*.md"}}<br>
			doing somthing ... <br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>query_string(regexp, *pages) *pages</td>
		<td>query pages by string</td>
		<td>{{with pages "/" | query_string "keyword"}}<br>
			doing somthing ... <br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>query_term(regexp, *pages) *pages</td>
		<td>query pages by term</td>
		<td>{{with pages "/" | query_term "keyword"}}<br>
			doing somthing ... <br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>query_term_prefix(regexp, *pages) *pages</td>
		<td>query pages by term prefix</td>
		<td>{{with pages "/" | query_term_prefix "prefix"}}<br>
			doing somthing ... <br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>query_match(match, *pages) *pages</td>
		<td>query pages by match</td>
		<td>{{with pages "/" | query_match "keyword"}}<br>
			doing somthing ... <br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>query_match_phrase(match, *pages) *pages</td>
		<td>query pages by match phrase</td>
		<td>{{with pages "/" | query_match_phrase "keyword"}}<br>
			doing somthing ... <br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>query_num_range(field, start, end, *pages) *pages</td>
		<td>query pages by num range</td>
		<td>{{with pages "/" | query_num_range "weight" 1 10}}<br>
			doing somthing ... <br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>query_num_range_inclusive(field, start, end, inclusive_start, inclusive_end, *pages) *pages</td>
		<td>query pages by num range inclusive</td>
		<td>{{with pages "/" | query_num_range "weight" 1 10 true true}}<br>
			doing somthing ... <br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>query_date_range(field, start, end, regexp, *pages) *pages</td>
		<td>query pages by date range</td>
		<td>{{with pages "/" | query_date_range "publish" "2015-07-01" "2015-07-30"}}<br>
			doing somthing ... <br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>query_date_range_inclusive(field, start, end, inclusive_start, inclusive_end, *pages) *pages</td>
		<td>query pages by date range inclusive</td>
		<td>{{with pages "/" | query_date_range_inclusive "publish" "2015-07-01" "2015-07-30" true true}}<br>
			doing somthing ... <br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>order_by(field, ascend, *pages) *pages</td>
		<td>query pages order by</td>
		<td>{{with pages "/" | order_by "publish" false}}<br>
			doing somthing ... <br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>all(*pages) []*page</td>
		<td>get all pages' array from pages object</td>
		<td>{{with pages "/" | order_by "publish" false}}<br>
				{{range .| all}}<br>
					{{.| title}}<br>
				{{end}}<br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>limit(offset, size, *pages) []*page</td>
		<td>get limit pages' array from pages object</td>
		<td>{{with pages "/" | order_by "publish" false}}<br>
				{{range .| limit 0 20}}<br>
					{{.| title}}<br>
				{{end}}<br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>limit_by_page(page_no, page_size, *pages) []*page</td>
		<td>get limit pages' array by pagination from pages object</td>
		<td>{{with pages "/" | order_by "publish" false}}<br>
				{{range .| limit_by_page 0 20}}<br>
					{{.| title}}<br>
				{{end}}<br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>pagination(url_base, size, *pages) []*pagination_item</td>
		<td>get pagination items' array from pages object</td>
		<td>{{with pages "/" | order_by "publish" false}}<br>
				{{range .| pagination "/list.html"  20}}<br>
					 doing pagination item ...<br>
				{{end}}<br>
			{{end}}<br>
		</td>
	</tr>
	<tr>
		<td>page_link(*pagination_item)</td>
		<td>get pagination item's link</td>
		<td></td>
	</tr>
	<tr>
		<td>page_no(*pagination_item)</td>
		<td>get pagination item's no</td>
		<td></td>
	</tr>
	<tr>
		<td>page_size(*pagination_item)</td>
		<td>get pagination item's size</td>
		<td></td>
	</tr>
	<tr>
		<td>total(*pages)</td>
		<td>get pages total count</td>
		<td>{{pages "/" | nested true | total}}</td>
	</tr>
</table>