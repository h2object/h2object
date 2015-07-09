H2OBJECT
========

[English](https://github.com/h2object/h2object/blob/master/doc/english/index.md)

## H2OBJECT 是什么?

**又一个Markdown建站工具, 但不只是工具!**

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
-	二维码分享文章
-	自由灵活的网站风格, 不仅仅是博客也可以是其它基于纯内容的站点
-	全文检索支持, 不需要引入第三方插件进行页面的检索
-	云存储支持, 实现静态资源的自动转储云端
-	自定义主题发布、分享, 同样一键操作

## 快速开始

### 程序下载

如果你无法直接源码安装或者不想被GFW扰乱心绪, 就直接下载可执行程序吧:

-	mac osx

	[h2object-darwin-amd64.tar.gz](http://dl.h2object.io/h2object/macosx/1.0.1.tar.gz)

-	linux

	[h2object-linux-amd64.tar.gz](http://dl.h2object.io/h2object/linux/1.0.1.tar.gz)

-	windows

	[h2object-windows-amd64.tar.gz](http://dl.h2object.io/h2object/windows/1.0.1.tar.gz)

解压安装,将 h2object 放入系统执行路径中。

### 源码安装

````
	$: go get github.com/h2object/h2object
````

#### 国内安装吐槽(VPN 用户跳过)

[取经求助](http://tangseng99.com)

TMD!!! 虽然已经竭尽全力减少对墙外包的依赖, 国内安装还是得提前做些准备工作:

由于项目中使用了以下两个国内绝对不能直接go get 的依赖包。
	
 * golang.org/x/net
 * golang.org/x/text

请在 GOPATH 目录下创建相应目录:

	mkdir -p $GOPATH/src/golang.org/x
	cd $GOPATH/src/golang.org/x
	git clone https://github.com/golang/net.git
	git clone https://github.com/golang/text.git


### 离线指南

安装完成后, 可以直接下载 指南主题 到本地目录, 快速开始一个本地的 h2object 指南站点.
````
	$: h2object -w=/path/to/work theme pull h2object/tutorial:1.0.0
	$: h2object -w=/path/to/work http start
````

### 在线指南

[tutorial.h2object.io](http://tutorial.h2object.io)

![tutorial](http://tutorial.h2object.io?qrcode=256)

### 博客站点

通过 h2object 快速创建博客站点:

````
	$: h2object -w=/path/to/blog theme pull h2object/blog:1.0.0
	$: h2object -w=/path/to/blog http start
````
更多主题查询:

````
	$: h2object theme search
````
等待您来发布更多主题。


### 一键发布

![deploy command](https://github.com/h2object/h2object/blob/master/doc/img/deploy.png)

在本地[deploy]配置项设置好, 远端h2object应用的相关参数即可以一键进行发布。具体配置包括：

````
[deploy]
# 远端部署服务 域名 或 地址
host= h2object.io
# 远端部署服务 端口
port= 80
# 远端应用ID
appid= 
# 远端应用密钥
secret= 
````
### 二维码分享

所有 h2object 服务提供的资源URI, 均可以通过 url 参数<kbd>qrcode=256</kbd>, 设置好相应的二维码像素值。即可产生二维码图片。

### 加入QQ群讨论(159823022)

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



