Conception
======

[Back](https://github.com/h2object/h2object/blob/master/doc/english/index.md) 

### How does it work?

**h2object** generate html pages from the follow files:

  * **markdown** as the content of the page
  * **template** as the style of the page 

if you want to create a content website without a database, h2object is a good option. And it's very fast. 

### Working Directory

Before h2object start http service, you need to set the working directory firstly, as the following command:

<kbd>h2object -w=/path/to/work http start</kbd>

In the h2object starting process, the application will create the working directory if it not exists.

After the directory created, the tree view of the folder will like:

````
\working directory
|-- markdowns (folder) 
|-- templates (folder)
|-- statics (folder)
|-- storage (folder)
|-- indexes (folder)
|-- logs (folder)
|-- .tmp (folder)
|-- h2object.conf
\-- h2object.pid

````

the folders named  **storage** 、 **indexes** 、 **logs** 、**.tmp** are system folders, you shouldn't change their content mannully.

And the **markdowns**, **templates**, **statics** folders will be the folder where you will put according files into.

If you want to access the files from h2object http service. The files' URI will be :

<table class="table">
	<thead>
		<th>file path</th>
		<th>http uri</th>
	</thead>
	<tr>
		<td>markdowns/index.md</td>
		<td>/index.md</td>
	</tr>
	<tr>
		<td>markdowns/docs/index.md</td>
		<td>/docs/index.md</td>
	</tr>
	<tr>
		<td>templates/index.html</td>
		<td>/index.html</td>
	</tr>
	<tr>
		<td>templates/guide/index.html</td>
		<td>/guide/index.html</td>
	</tr>
	<tr>
		<td>statics/img/logo.png</td>
		<td>/img/logo.png</td>
	</tr>
</table>

h2object service just generate the html page by the http uris' suffix, which the h2object configures.

In the h2object.conf file, you will find <kbd>markdown.suffix</kbd>、 <kbd>templates.suffix</kbd> items.

As an example:

If the markdown.suffix is <kbd>md</kbd>, then if you want to access uri **/demo.md**, h2object will find the **markdowns/demo.md** at the working directory. If the file exists, it will generate the according page for you.Otherwise, the application will out put **Not Found** error.




