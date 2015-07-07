系统配置
======

[返回目录](https://github.com/h2object/h2object/blob/master/doc/chinese/index.md) 

初始情况下, h2object 启动过程中会在工作目录中自动创建 **h2object.conf** 配置文件。

h2object 配置非常简单, 主要包括四个**section**, 分别如下:

  *  应用主配置 [h2object](https://github.com/h2object/h2object/blob/master/doc/chinese/configure.md#h2object)
  *  第三方配置 [third](https://github.com/h2object/h2object/blob/master/doc/chinese/configure.md#third)
  *  部署配置  [deploy](https://github.com/h2object/h2object/blob/master/doc/chinese/configure.md#deploy)
  *  主题配置  [theme](https://github.com/h2object/h2object/blob/master/doc/chinese/configure.md#theme)

<a name="h2object"></a>
####	1. 应用配置(h2object)

以下是默认站点工作路径下自动创建的 应用配置。

````
[h2object]
# 应用服务地址 格式(域名 或 IP:PORT), 主要用于 QRCODE 路径
host= 
# 资源默认访问页
index= 
# markdown 页面缓存时间
markdown.cache= 10m
# markdown 文件后缀
markdown.suffix= md,markdown
# template 文件后缀
template.suffix= html,htm,tpl
# 运行模式: true, 开发; false, 生产;
develope.mode= false
# 应用ID
appid= rUbZ2GIwa7JTgJLimw2EOs3vkMjEez0
# 应用密钥
secret= OcwRiV3zY1JJ9MTvc3ZAjRskMcufQOZ2zSA
````

其中**appid**, **secret**两个配置属于系统配置是系统启动过程中自动创建的随机串。用户可以根据个人喜欢修改, 常规情况下不建议修改。

-	index 配置项

	默认情况下, 如果 URL 标识的是目录, 默认会寻找 该目录下的 index 页面。

-	markdown.cache 配置项

	markdown 文件生成 html 页面后，在系统缓存中的缓存时间。有效单位为: 1ns, 1us, 1ms, 1s, 1m, 1h.

-	markdown.suffix 配置项

	markdown文件后缀, 对于多后缀, 请使用 <kbd>,</kbd> 进行分割。 

-	template.suffix 配置项

	template文件后缀, 对于多后缀, 请使用 <kbd>,</kbd> 进行分割。

-	develope.mode 配置项

	开发模式, 默认情况是关闭状态。打开开发模式后, 系统提供 <kbd>.page</kbd> 和 <kbd>.system</kbd> 后缀用于查询系统对象数据。不过此类操作需要具备管理员权限。

<a name="third"></a>
####	2. 第三方配置(third)

目前 h2object 程序主要提供的第三方功能主要是七牛云存储功能, 用户通过配置该section, 打开七牛云存储, 实现静态资源的云端存储加速访问。

具体配置如下, 不再一一解释。

````
[third]
# 是否启用七牛云存储
qiniu.enable= false
# 七牛云存储提供的 access key
qiniu.appid= 
# 七牛云存储提供的 secret key
qiniu.secret= 
# 七牛云存储提供的 访问域名
qiniu.domain= 
# 七牛云存储提供的 资源存放空间名
qiniu.bucket= 
````
<kbd>注意:</kbd>本地运行时, 本地静态资源是不会存放到云端。该功能只在进行 本地=》远端 的发布部署时，在远端生效。

<a name="deploy"></a>
####	3. 部署配置(deploy)

h2object的部署过程非常简单, 只需要:

  *  远端: 在远端服务器安装 h2object 程序, 并且运行该程序;

  *  本地: 在本地工作路径配置好远端部署程序的访问域名(或IP)与端口, 同时, 设置好远端程序的应用ID与密钥即可。

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

完成该section部署后, 用户就可以在本地发起

<kbd>h2object -w=/path/to/work deploy push</kbd>

命令, 一键将本地站点进行远端发布了。

<a name="theme"></a>
####	4. 主题配置(theme)

h2object.io 作为 h2object 开发服务商, 通过 

<kbd>h2object -w=/path/to/work theme push </kbd>

命令进行主题的上传、分享功能。运行该命令需要进行用户登陆, 登陆操作请参考[用户认证](/authorize.md)。

一旦用户需要进行主题的上传分享, 就需要通过该section进行主题的配置说明。具体配置项如下:

````
[theme]
# 主题提供商名称
provider= h2object
# 主题名称
name= demo
# 主题类型 (0: 所有人 1: 会员)
catagory= 0
# 主题版本号
version= 0.0.1
# 主题描述
description= h2object demo site

````







