markdown
========

[Back](https://github.com/h2object/h2object/blob/master/doc/english/index.md) 

h2object service use markdown file's content to generate the html page's content. So how to write the markdown file? It's simple, h2object can recognize the normal markdown files. And with the improvement, h2object can parse file metadata.

Here is the markdown file format:

````
{
	#json string as the metadata part.
}

this is the markdown content part.

````

Example:

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

Really easy, if you want to get more examples just check the tutorial site markdown files.

Now next step is how to get the markdown content to the html pages? That's the reason for [template](/template.md) Existence. 


