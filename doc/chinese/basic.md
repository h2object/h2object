基本概念
======

[返回目录](https://github.com/h2object/h2object/blob/master/doc/chinese/index.md) 

### 工作原理

**h2object** 通过利用

  * **markdown** 文件提供生成html页面的内容
  * **template** 文件提供生成html页面的风格

生成最终的html网页提供服务。

如果你要创建一个基于内容型的站点, h2object 是一个非常好的选择。因为, 它真的很快而且很方便。

### 工作目录

Before h2object start http service, you need to set the working directory firstly, as the following command:

在 启动 h2object 本地服务时, 必须指定其工作目录:

<kbd>h2object -w=/path/to/work http start</kbd>

服务启动完成后, 打开工作目录会看见如下工作子目录:

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

其中, **storage** 、 **indexes** 、 **logs** 、**.tmp** 属于系统文件夹, 请不要手动修改相关文件。

而以下三个目录，

-	目录**markdowns**

	存放 markdown 内容文件, 文件后缀参考[系统配置](https://github.com/h2object/h2object/blob/master/doc/chinese/configure.md)中的<kbd>markdown.suffix</kbd>配置。

	后缀不符合条件的文件将会存放到statics目录中。

-	目录**templates**

	存放 template 模板文件, 文件后缀参考[系统配置](https://github.com/h2object/h2object/blob/master/doc/chinese/configure.md)中的<kbd>template.suffix</kbd>配置。

	后缀不符合条件的文件将会存放到statics目录中。

-	目录**statics**

	存放所有项目中的静态资源文件。

你可以直接通过文件路径(不包括子目录名)访问以上三个目录中的任意文件。如:

<table class="table">
	<thead>
		<th>文件路径</th>
		<th>访问URI</th>
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

h2object 服务通过访问uri的后缀进行判断从何处找出相应的原始文件并最终生成html页面.




