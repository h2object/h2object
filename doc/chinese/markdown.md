markdown
========

[返回目录](https://github.com/h2object/h2object/blob/master/doc/chinese/index.md) 

在 h2object 服务中, 通过 markdown 文件内容提供生成最终html页面的正文。常规的 markdown 文件, 可以在 h2object 中直接支持并生成相应的html页面。 在 h2object 中, 对可以识别的markdown文件增加了一项功能, 即文件的头信息的支持。

下面是 h2object 对 markdown 文件格式的支持语法:

````
{
	#json string as the metadata part.
}

this is the markdown content part.

````

例子:

````
{
	"template":"default.html",
	"title":"this is article title",
	"publish":"2015-07-05",
	"weight": 2
}

## title
=====

### item 1

this is item 1 for the article

code example:

	````
		some code 
	````

### item 2

this is item 2 for the article

````

从样例可知编写一篇 h2object 可识别的 markdown 文件非常简单。即

````
	h2object markdown = json 格式的文件头信息 + markdown 文件正文
````
如果文件中没有 json 格式的文件头信息, h2object 照样可以识别。

下一步, 就是如何将 markdown 文件正文 输出到模板页中, 请参考 [template](https://github.com/h2object/h2object/blob/master/doc/chinese/template.md) 文档。

