h2object 命令
======

[返回目录](https://github.com/h2object/h2object/blob/master/doc/chinese/index.md) 

h2object 以下几个方面的操作命令:

  * 服务启停命令 [http](https://github.com/h2object/h2object/blob/master/doc/chinese/command.md#http) 
  * 主题相关命令 [theme](https://github.com/h2object/h2object/blob/master/doc/chinese/command.md#theme) 
  * 用户认证命令 [auth](https://github.com/h2object/h2object/blob/master/doc/chinese/command.md#auth) 
  * 站点部署命令 [deploy](https://github.com/h2object/h2object/blob/master/doc/chinese/command.md#deploy) 

h2object 成功安装后, 可以通过 <kbd>-h</kbd> 选项查询具体命令的操作帮助。

````
$: h2object -h
$: h2object http -h
$: h2object theme -h
$: h2object auth -h
$: h2object deploy -h
````

<a name="http"></a>
####  1. 服务命令 (http)

  * **服务启动**

  ````
  $: h2object -w=/path/to/workdir -d http start
  ````
通过<kbd>-d</kbd>选项, 后台启动h2object服务.windows环境不支持该选项功能。

  * **服务停止**

  ````
  $: h2object -w=/path/to/workdir http stop
  ````

  * **服务重载**

  ````
  $: h2object -w=/path/to/workdir http reload
  ````

<a name="theme"></a>
####  2. 主题命令 (theme)

  * **主题查询**

  ````
  $: h2object theme search <keyword>
  ````
通过设置查询关键字查询 h2object.io 平台所有发布主题, 

  * **主题下载**

  ````
  $: h2object -w=/path/to/workdir theme pull <provider/name:version>
  ````
通过主题拉取命令, 可以直接将平台已发布的主题下载到指定的工作目录。

  * **主题发布**

  ````
  $: h2object -w=/path/to/workdir theme push
  ````

主题发布的前提是, 用户必须在 **h2object.conf** 主题配置项中进行相应的主题设置. 

同时, 发布主题必须首先进行用户的登录, 具体用户登录参考[用户认证](/authorize.md)。


<a name="deploy"></a>
####  3. 部署命令 (deploy)

部署命令的操作前提是, 用户必须在 **h2object.conf** 已经对 deploy section 进行正确设置, 才能进行以下部署操作。

  * **站点发布**

  ````
  $: h2object -w=/path/to/workdir deploy push
  ````

除了可以整体站点一键发布外, 还可以发布具体的路径或者文件。如, 

发布新增 markdow 文件:

  ````
  $: h2object -w=/path/to/workdir deploy push markdowns/newfile.md
  ````

发布具体 template 某个子目录:

  ````
  $: h2object -w=/path/to/workdir deploy push templates/some/dir/
  ````

  * **站点备份**

通过以下命令可以对配置项中的站点进行整体下载备份到指定目录。

  ````
  $: h2object -w=/path/to/workdir deploy pull 
  ````

也可以通过设置具体备份目录, 将具体目录备份到本地。如下, 备份站点所有markdown文件。

  ````
  $: h2object -w=/path/to/workdir deploy pull markdowns/
  ````

<a name="auth"></a>
####  4. 认证命令 (authorize)

认证命令是针对 portal.h2object.io 平台的. 用户可以直接通过命令进行该平台的注册、登录、登出操作。认证命令目前主要是针对主题操作。

  * **用户注册**

新用户注册, 提供邮箱和密码即可注册新用户. 注册完成后, 需要登录邮箱完成账户的激活操作。

  ````
  $: h2object auth new
  ````  

  * **用户登录**

用户登录操作是针对工作路径的, 这样下次用户在对应的工作路径就可以免登录或简单重登录即可。
  
  ````
  $: h2object -w=/path/to/workdir auth login
  ````  

  * **用户状态**

查询当前工作路径中，是否有用户登录信息以及具体用户信息。
  
  ````
  $: h2object -w=/path/to/workdir auth status
  ````  

  * **用户登出**

用户从当前工作路径中登出。
  
  ````
  $: h2object -w=/path/to/workdir auth logout
  ````



