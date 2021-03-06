template
========

[Back](https://github.com/h2object/h2object/blob/master/doc/english/index.md) 

In h2object service, the application use markdown to provide the content for the html pages; and use **template** to provide the style of the html pages. That's why the template is very important.

As we known already, we can access the page by the uri with the markdown suffix and tempalte suffix directly.

 * 1.**access html page generated by markdown file with markdown suffix**

 * 2.**access html page generated by template file with template suffix**

If you want to access markdown file with suffix directly, then if the markdown file's metadata has the <kbd>template</kbd> field, the html page generated will will the template style, other wise the html page only out put the markdown generate html string without any template style.

If the markdown file's metadata has <kbd>template</kbd> field, then generate the html page's template will has the <kbd>.page</kbd> render data parameter represent the markdown file.

There is an example to show you how to get the markdown file content at template.

markdown demo file:

````
{
	"template":"demo_template.html",
	"title":"demo",
	"keywords":"demo,test,h2object",
	"summary":"this is the summary of the demo markdown file",
	"published_date":"2015-07-06",
	"weight":2
}

demo
====

this is the demo content.
````

template demo file:
````
<html>
	<head>
		<title>{{.page | title}}</title>	
		<meta keywords='{{.page | meta "keywords"}}'>
	</head>
	<body>
		<h1>{{.page | title}}</h1>
		<h4>{{.page | meta "published_date" | datetime}}</h4>
		<p>
		{{.page | markdown}}
		</p>
	</body>
<html>

````
This is access the markdown file with a template field in it's metadata directly.

There is another way to generate the html pages, is access the template file directly. If we access the template file directly generated pages, the template can't use <kbd>.page</kbd> render page parameter.

If the template can't use the <kbd>.page</kbd> get the content from the markdown file, there are two another way to get the markdown file content. The h2object service provide two template's functions to get the markdown file content. The functions are:

 *  **page** [uri]

 	page function, with the param with the markdown file uri, will return the markdown file content as the result.

 *  **pages** [uri]

 	pages function, with the param with the markdown file uri or folder, will return the markdown files array as the result.

there is an example, template file name is **list.html**:

````
<html>
	<head>
		<title>list</title>	
	</head>
	<body>
		<h1>list</h1>
		<h4></h4>
		{{with pages "/" | query_date_range_inclusive "2015-06-01" "2015-06-30" true true}}
			<h1>pages list</h1>
			{{range . | all }}
			<li><a href="{{. | uri}}">{{. | title}}</a></li>
			{{end}}
			<h5>pagination</h5>
			<ul>
			{{range . | pagination "/list.html" ($.url | param "size" 5)}}
			<li><a href="{{. | page_link}}">{{. | page_no }}({{. | page_size}})</a></li>	
			{{end}}
			</ul>
		{{end}}
	</body>
<html>
````

Of course, the article is only the template working mechanism's description. If you want know more about the template's functions, please deep into the template's [functions](/functions.md), and you should better know about **go html/template package** if you can.

