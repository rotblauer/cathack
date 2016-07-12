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
