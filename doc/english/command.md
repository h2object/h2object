h2object commands
======

[Back](https://github.com/h2object/h2object/blob/master/doc/english/index.md) 

h2object application has the following commands:

  * http command [http](https://github.com/h2object/h2object/blob/master/doc/english/command.md#http) 
  * theme command [theme](https://github.com/h2object/h2object/blob/master/doc/english/command.md#theme) 
  * deploy command [deploy](https://github.com/h2object/h2object/blob/master/doc/english/command.md#deploy) 
  * auth command [auth](https://github.com/h2object/h2object/blob/master/doc/english/command.md#auth)   

After Installed h2object, you can use <kbd>-h</kbd> option to get the help information.

````
$: h2object -h
$: h2object http -h
$: h2object theme -h
$: h2object auth -h
$: h2object deploy -h
````

<a name="http"></a>
####	1. http command

  * **start http service**

  ````
  $: h2object -w=/path/to/workdir -d http start
  ````
You cant add <kbd>-d</kbd>option to starting the service at daemon mode. This option has no effect at windows platform.

  * **stop http service**

  ````
  $: h2object -w=/path/to/workdir http stop
  ````

  * **reload http service**

  ````
  $: h2object -w=/path/to/workdir http reload
  ````

<a name="theme"></a>
####	2. theme command

  * **search themes**

  ````
  $: h2object theme search <keyword>
  ````
if <keyword> is empty, will search all themes at h2object.io platform.

  * **download theme**

  ````
  $: h2object -w=/path/to/workdir theme pull <provider/name:version>
  ````
this command will directly pull the theme to the dest working directory. After it done, start the service, you will find the results at once.

  * **upload theme**

  ````
  $: h2object -w=/path/to/workdir theme push
  ````

this command will push the current working directory site theme to the h2object.io platform with attribute's private. If you want to publish the theme, you need login the portal.h2object.io to publish the private theme mannully.

<a name="deploy"></a>
####	3. deploy command

Before deploy the local h2object site, you need configure the [deploy] section first.

  * **deployment** from local to remote

  ````
  $: h2object -w=/path/to/workdir deploy push 
  ````

You can use the command deploy the whole site. Or you can use the following command to deploy files or directorys.

  ````
  $: h2object -w=/path/to/workdir deploy push markdowns/add.md statics/
  ```` 
this example command will deploy add.md file and the whole statics directory files.

  * **deployment** from remote to local

if you want backup remote files to local, just use the following command:

  ````
  $: h2object -w=/path/to/workdir deploy pull 
  ````

this command will pull the whole site to the local;

  ````
  $: h2object -w=/path/to/workdir deploy pull markdowns/
  ````

this command only pull the markdowns directory files to the local working directory.

<a name="auth"></a>
####	4. auth command

h2object application authorization is base on the portal.h2object.io platform. And we can use the h2object auth command to sign up, sign in and sign off directly. h2object appliaction provide auth command , because when user want to share his site theme, there must have a place to store and index it.

  * **sign up**

  ````
  $: h2object auth new
  ````  

  * **sign in**

h2object sign in must set the working directory , as the following :
	
  ````
  $: h2object -w=/path/to/workdir auth login
  ````  

  * **sign status**

check the working directory authorized user infomation.
	
  ````
  $: h2object -w=/path/to/workdir auth status
  ````  

  * **sign off**

user sign off from the working directory.

  ````
  $: h2object -w=/path/to/workdir auth logout
  ````





