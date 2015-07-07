System Configure
======

[Back](https://github.com/h2object/h2object/blob/master/doc/english/index.md) 

In the h2object service starting process, the application will generate a default **h2object.conf** file at the working directory. The configure has four section parts, as the following:

  *  [h2object](https://github.com/h2object/h2object/blob/master/doc/english/configure.md#h2object)
  *  [third](https://github.com/h2object/h2object/blob/master/doc/english/configure.md#third)
  *  [deploy](https://github.com/h2object/h2object/blob/master/doc/english/configure.md#deploy)
  *  [theme](https://github.com/h2object/h2object/blob/master/doc/english/configure.md#theme)

<a name="h2object"></a>
####	1. h2object section

this section is major section for h2object application.

````
[h2object]
# service host format(domain or IP:PORT), if empty, the value will be the listenning ip:port
host=
# default uri index page name
index= 
# markdown page's cache duration, and the unit should in (1ns, 1us, 1ms, 1s, 1m, 1h)
markdown.cache= 10m
# markdown files' suffixes, the delemiters to be: ,
markdown.suffix= md,markdown
# template files' suffixes, the delemiters to be: ,
template.suffix= html,htm,tpl
# running mode: true, develope mode; false, production mode;
develope.mode= false
# application id
appid= rUbZ2GIwa7JTgJLimw2EOs3vkMjEez0
# application secret
secret= OcwRiV3zY1JJ9MTvc3ZAjRskMcufQOZ2zSA
````

the **appid**, **secret** items are the system auto generates configures, recommend not to be changed.

If the application is running at develope mode, the applicaiton will append two suffixes for developers to debug the data. And the suffixes are <kbd>.page</kbd> 和 <kbd>.system</kbd> , but only can used will administrator rights.

<a name="third"></a>
####	2. third section

Till now, h2object service only support qiniu.com cloud storage service(like CDN service). And the configure is as following:

````
[third]
qiniu.enable= false
qiniu.appid= 
qiniu.secret= 
qiniu.domain= 
qiniu.bucket= 
````
<kbd>Attention:</kbd>If h2object application running at local, the statics files will not transfer to qiniu cloud. Only when you do deployment, push the local static files to remote server with the remote application's qiniu.enable equal true, the static files will be stored at qiniu cloud service. 

<a name="deploy"></a>
####	3. deploy section

The h2object service's deployment is vary convinent, you just configure the h2object.conf's deploy section items, then you can deploy the local to the remote by one command, or even pull the remote site to local.

````
[deploy]
# remote h2object application's listenning address or domain name
host= h2object.io
# remote h2object application's listenning port
port= 80
# remote h2object application's id
appid= 
# remote h2object application's secret
secret= 
````

After configured this section, your cant use the command:

<kbd>h2object -w=/path/to/work deploy push</kbd>

to deploy the local site to the remote site.

<a name="theme"></a>
####	4. theme section

If you friend like your h2object site, you can sharing it by push the site theme to the h2object.io platform.

Before you push, you need configure the theme parameters firstly.


````
[theme]
# theme provider's name
provider= h2object
# theme's name
name= demo
# theme's catagory (0: free; 1: member free;)
catagory= 0
# theme's version
version= 0.0.1
# theme's description
description= h2object demo site

````

After configured, use the following command to push the theme to the platform.

<kbd>h2object -w=/path/to/work theme push </kbd>

**Attention**: Before push the theme to the platform, you need authorized at the directory on the h2object.io platform. Any Questions About [Authorization](/command.md#auth), please read the doc.







