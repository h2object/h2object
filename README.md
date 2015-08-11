H2OBJECT
========

[English](https://github.com/h2object/h2object/blob/master/doc/english/index.md)

## H2OBJECT 是什么?

H2OBJECT 同 hexo, hugo 一样是一个通过 markdown 文件快速创建内容型站点的工具。不同的是, 
H2OBJECT 参考了 Docker Pull/Push 的方式实现快速本地到线上站点的一键发布。

![deploy.png](https://github.com/h2object/h2object/blob/master/doc/img/deploy.png)

同时, 用户可以申请[h2object.io](http://h2object.io)平台提供的仅运行h2obect应用的docker容器运行线上站点.

## 快速开始

### 二进制安装

如果你无法直接源码安装或者不想被GFW扰乱心绪, 就直接下载可执行程序吧:

[h2object-darwin-amd64.tar.gz](http://dl.h2object.io/h2object/macosx/1.0.3.tar.gz)

[h2object-linux-amd64.tar.gz](http://dl.h2object.io/h2object/linux/1.0.3.tar.gz)

[h2object-windows-amd64.tar.gz](http://dl.h2object.io/h2object/windows/1.0.3.tar.gz)

解压后,将 h2object 放入系统执行路径中。

### 源码安装

````
	$: go get github.com/h2object/h2object
````

#### 国内安装吐槽(VPN 用户跳过)

[取经求助](http://tangseng99.com)

虽然已经竭尽全力减少对墙外包的依赖, 国内安装还是得提前做些准备工作:

由于项目中使用了以下两个国内绝对不能直接go get 的依赖包。
	
 * golang.org/x/net
 * golang.org/x/text
 * golang.org/x/image

请在 GOPATH 目录下创建相应目录:

	mkdir -p $GOPATH/src/golang.org/x
	cd $GOPATH/src/golang.org/x
	git clone https://github.com/golang/net.git
	git clone https://github.com/golang/text.git
	git clone https://github.com/golang/image.git

### 本地运行

````
$: h2object -w=/path/to/workdir http start
````

### 站点主题

#### 主题查询

````
$: h2object theme search
````

#### 主题下载

下载他人分享的站点主题

````
$: h2object  -w=/path/to/workdir theme pull [provider/name:version]
````

#### 主题发布

将个人站点主题分享给其他用户，请先在配置文件中配置好[theme]项

````
$: h2object  -w=/path/to/workdir theme push
````

### 容器申请

* 创建容器

容器创建前必须到[h2object.io](http://h2object.io)平台申请容器邀请码。

````
$: h2object -w=/path/to/workdir container create [邀请码]
````

* 运行容器

````
$: h2object -w=/path/to/workdir container start [container id]
````

### 站点发布

将容器提供的appid,secret,host,port配置到本地发布站点的[deploy]配置项中
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
再通过以下命令一键发布站点

````
$: h2object -w=/path/to/workdir deploy push
````

### 加入QQ群讨论(159823022)

### 开发指南

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



