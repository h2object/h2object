H2OBJECT
========

## H2OBJECT 是什么?

H2OBJECT 项目, 最开始开发的目的如同其名称一样(HTTP to OBJECT), 主要提供基于对象的HTTP接口服务。

如今, 项目演变成为提供快速建站与站点发布的工具,主要通过:

 * **markdown**文件提供页面内容

 * **template**文件提供页面风格
	
快速生成网页服务。

同时, 考虑到其它类似项目(Hugo)等, 在站点发布上的繁琐步骤, H2OBJECT 参考了 Docker 的 PUSH/PULL 机制,

可以通过 Deploy 命令进行, 快速的 **本地** 与 **远程** 的 PUSH/PULL 操作实现一键发布。

现在, 如果你要搭建一个内容型的站点, 从开始到发布通过H2OBJECT只需要几个命令就可以完成。

### 设计原则

-	一键安装、一键发布
-	灵活的网站风格, 不仅仅是博客也可以是其它基于纯内容的站点
-	全文检索支持, 不需要引入第三方插件进行页面的检索
-	云存储支持, 实现静态资源的自动转储云端
-	自定义主题发布、分享, 同样一键操作

## 快速开始

### 源码安装

````
	$: go get github.com/h2object/h2object
	$: go build
	$: cp h2object /user/local/bin/
````
### 一键发布

![deploy command](https://github.com/h2object/h2object/blob/master/doc/img/deploy.png)

### 快速指南

安装完成后, 可以直接下载 指南主题 到本地目录, 快速开始一个本地的 h2object 指南站点.
````
	$: h2object -w=/path/to/work theme pull h2object/tutorial.ch:1.0.0
	$: h2object -w=/path/to/work http start
````

### 使用指南

-	[基本概念](https://github.com/h2object/h2object/blob/master/doc/chinese/basic.md)
-	[系统配置](https://github.com/h2object/h2object/blob/master/doc/chinese/configure.md)
-	[服务命令](https://github.com/h2object/h2object/blob/master/doc/chinese/command.md)
-	[markdown](https://github.com/h2object/h2object/blob/master/doc/chinese/markdown.md)
-	[template](https://github.com/h2object/h2object/blob/master/doc/chinese/template.md)
-	[模板函数](https://github.com/h2object/h2object/blob/master/doc/chinese/functions.md)

### 参考&使用的项目

-	[revel](https://github.com/revel/revel)
-	[bleve](https://github.com/blevesearch/bleve)
-	[blotdb](https://github.com/boltdb/bolt)
-	[hugo](https://github.com/spf3/hugo)
-	[docker](https://github.com/docker/docker)

非常欢迎您使用并推荐 H2OBJECT 项目。
