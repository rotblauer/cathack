## Syncthing

#### Installing on local Mac.
> [digital ocean tutorial on installing server-server](https://www.digitalocean.com/community/tutorials/how-to-install-and-configure-syncthing-to-synchronize-directories-on-ubuntu-14-04)

> [syncthing getting started docs](https://docs.syncthing.net/intro/getting-started.html)

> [latest release builds](https://github.com/syncthing/syncthing/releases/tag/v0.14.0)

##### I did this:
__local__
```
$ uname -a # determine cpu architecture
$ `clicked download` #auto unpacks. cli wasn't having it
$ sudo cp syncthing /usr/local/bin
$ syncthing
$ visit the ip:port address listed.

$ OR... # better forget everything above
$ brew install syncthing
$ brew link --overwrite syncthing # oh but actually the brew is old (0.13) which is incompatible with freya's newer version (released like 6 hours ago)
$ brew services start syncthing # so it start automaticlike
```

Add Freya's ID as a Remote Device as per the DO tutorial. 

Select the folder/s you want to sync and attend to the invites to sync and rock and rolla. 

> __launchctling syncthing__ was a bitch. i wound up using [this](https://github.com/syncthing/syncthing/tree/master/etc/macosx-launchd) as a working setup, with the `syncthing` executable in $HOME/bin. I also had to use the -w option in `launchctl load -w ~/Library/LaunchAgents/syncthing.plist ` to override a launchctl disabled service problem per [this SO question](http://stackoverflow.com/questions/26334414/cant-make-postgresql-load-at-startup-in-mac-os)

__freya__
```
$ more /proc/cpuinfo # get cpu info for which .tar to wget
$ wget https://github.com/syncthing/syncthing/releases/download/v0.14.0/syncthing-linux-amd64-v0.14.0.tar.gz
$ tar xzvf syncthing*.tar.gz # unpack tar
freyabison@ribbon:~ $ cd syncthing-linux-amd64-v0.14.0
freyabison@ribbon:~/syncthing-linux-amd64-v0.14.0 $ ll
total 16352
-rw-r--r-- 1 freyabison pg8746732    16725 May 21 18:48 LICENSE.txt
-rw-r--r-- 1 freyabison pg8746732     5169 Jun 26 03:52 AUTHORS.txt
-rw-r--r-- 1 freyabison pg8746732     3262 Jul 10 00:36 README.txt
-rwxr-xr-x 1 freyabison pg8746732 16707449 Jul 19 01:56 syncthing*
-rw-rw-r-- 1 freyabison pg8746732      241 Jul 19 02:23 syncthing.sig
drwxrwxr-x 2 freyabison pg8746732       58 Jul 19 02:23 extra/
drwxrwxr-x 7 freyabison pg8746732      147 Jul 19 02:23 etc/
drwxrwxr-x 2 freyabison pg8746732       32 Jul 19 02:23 .metadata/
freyabison@ribbon:~/syncthing-linux-amd64-v0.14.0 $ sudo cp syncthing /usr/local/bin
[sudo] password for freyabison: 
freyabison is not in the sudoers file.  This incident will be reported.
```

So add `PATH=$PATH:$HOME/syncthing-linux-amd64-v0.14.0/syncthing` to `.bash_profile` before exporting PATH and re source it. 

```
$ syncthing
```

> Note that Freya _does_ have the `Upstart` utility baked right in. Fanatastico. 

```
freyabison@ribbon:~/.config/syncthing $ sudo nano /etc/init/syncthing.conf
[sudo] password for freyabison: 
freyabison is not in the sudoers file.  This incident will be reported.
# Damn.
```

... edited `syncthing.conf` per:
`freyabison@ribbon:~/syncthing-linux-amd64-v0.14.0/etc/linux-upstart/user $ vim syncthing.conf`

```
freyabison@ribbon:~/syncthing-linux-amd64-v0.14.0/etc/linux-upstart/user $ sudo initctl start syncthing
[sudo] password for freyabison: 
freyabison is not in the sudoers file.  This incident will be reported.
```

OK, so skip the whole Upstarter thing.

```
# Figure out Freya's public ip. 
$ wget http://ipinfo.io/ip -qO -
$ cd ~/syncthing-linux-amd64-v0.14.0/ && ./syncthing 
$ open https://173.236.168.108:8384
```

- (see `.bash_profile` > `alias gosync="forever start -c bash ~/foreversyncthing.sh"`)


----
Found it [here](https://medium.com/@olahol/writing-real-time-web-apps-in-go-chat-4aa058644f73#.8whi6swyn), too slick: 


Startup if not working:

Run ~/chat.areteh.co/chat
[http://www.chat.areteh.co:5000](http://www.chat.areteh.co:5000)

```
$ ./freyafy.sh # send to freya
$ ./kickstart.sh # rebuild and restart server, run from freya
```

----

### Architectural concerns
How Massad does it. 

> Controllers and models are referenced via ArticleController{}
> and ArticleModel{}.
> They export anonymous functions which accept their structs 
> and which return named functions. 
> The main method initializes the db and a new struct for each 
> controller. (The given controller then initializes a new 
> model struct). 

_db_
- sets a global var and exports a getter
- db.Init() called from main.go

_models_
- do all of the db manipulation
- defines models structs

_controllers_
- call model methods and send responses

_main.go_
- calls db.Init()
