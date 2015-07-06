template
========

[返回目录](https://github.com/h2object/h2object/blob/master/doc/chinese/index.md) 

在 h2object 服务中, 应用通过 markdown 提供内容; 通过 template 提供页面样式。

内容可以是 markdown 的头信息 或者 正文 甚至是 html 标签节点。所有数据源 均可以通过[模板函数](https://github.com/h2object/h2object/blob/master/doc/chinese/functions.md)获取。所以在编辑页面风格的同时需要考虑数据源的导入。

在当前版本的 h2object 服务中， template 可以通过模板函数或变量 获取以下数据源:

 * 1. **url** 当前页面 访问URL 信息

 	在 template 可以通过 {{.url}} 输出变量 直接获取。并可通过 模板函数获取 访问url 具体属性信息。如:

 	{{.url | path}} 获取 url 路径

 	{{.url | param "size" 10}} 获取 url 查询Query参数信息

 * 2. **page** 当前页面访问的 markdown 文件数据, 或者是 具体路径的某个 markdown 文件数据

 	如果直接通过 markdown 文件路径直接访问指定文件, 并且在 markdown 的文件头中 指定了 template 属性。则在对应的template页面中, 可以通过 {{.page}} 访问该markdown文件。同样, 可以通过 具体函数 获取到 markdown 文件中的具体信息。 如

 	{{.page | title}} 获取markdown 文件 title

 	{{.page | meta "field"}} 获取markdown 文件头中某个字段对应的值

 	当然, 当个页面数据源除了在markdown 文件中指定 template 属性, 并在对应 template 通过{{.page}} 输出变量 获取以外, 还可以通过模板函数 page 直接获取。 如

 	{{page "/demo.md" | title}}

 * 3. **pages** 指定路径下 多个 markdown 文件数据

 	除了当个 markdown 文件数据的 template 展示外, 常规网站还有列表页面的需求。此时可以创建列表页面的template, 通过直接访问template文件生成的列表页面。同样, 多个markdown文件数据源, h2object 提供了 pages 模板函数获取。如:

 	{{pages "/"}} 获取所有markdowns目录下的文件数据。

 	同时 pages 对象提供了 丰富的查询接口供开发者使用, 具体函数参考[模板函数](https://github.com/h2object/h2object/blob/master/doc/chinese/functions.md).

样例1, 单 markdown 文件对应 template 文件的样例:

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

demo.md 文件通过设置 template 属性设置对应的模板文件. 用户可以通过 /demo.md 访问最终生成的页面效果。

样例2, 列表页面(list.html)

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

以上两个例子介绍了 模板文件的工作原理与机制。 h2object 文件模板主要采用了 go语言 的 html/template 开发包。所以很多功能是与之对应的。同时 h2object 提供了自有的模板函数, 请参考下节。

