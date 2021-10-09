package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"encoding/json"
	"strconv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"time"
	"crypto/md5"
)

type Handler func(*Context)

type Route struct {
	Pattern *regexp.Regexp
	Handler Handler
}

type App struct {
	Routes       []Route
	DefaultRoute Handler
}

type User struct {
	Id int
	Name string
	Email string
	Password string
}

type Post struct {
	Id int
	Caption string
	Image_url string
	Posted_time string
}


func NewApp() *App {
	app := &App{
		DefaultRoute: func(ctx *Context) {
			fmt.Println("Path not found")
		},
	}

	return app
}

func (a *App) Handle(pattern string, handler Handler) {
	re := regexp.MustCompile(pattern)
	route := Route{Pattern: re, Handler: handler}

	a.Routes = append(a.Routes, route)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{Request: r, ResponseWriter: w}

	for _, rt := range a.Routes {
		if matches := rt.Pattern.FindStringSubmatch(ctx.URL.Path); len(matches) > 0 {
			if len(matches) > 1 {
				ctx.Params = matches[1:]
			}

			rt.Handler(ctx)
			return
		}
	}

	a.DefaultRoute(ctx)
}

type Context struct {
	http.ResponseWriter
	*http.Request
	Params []string
}

var Users []User
var Posts []Post
var ctx = context.TODO()
var collection *mongo.Collection

func main() {
	Users = []User{
		User{Id:5, Name : "Aditya", Email : "baditya@gmail.com", Password : "asdf" },
		User{Id:9, Name : "tya", Email : "bada@gmail.com", Password : "asdf" },
	}
	Posts = []Post{
		Post{Id:6, Caption : "First post", Image_url : "htsdfdds", Posted_time : "34:24:34"},
		Post{Id:18, Caption : "second post", Image_url : "htds", Posted_time : "31:45:53"},
	}
	
	
	var mongoURI string = "mongodb+srv://dbUser:dbPassw0rd@cluster0.xqhl1.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
	
	// setup connection with mongoDb cluster
	
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil { log.Fatal(err) }
	collection = client.Database("tasker").Collection("tasks")
	
	app := NewApp()
	
		
	app.Handle(`/users/([^/]+)$`, func(ctx *Context) {
		if ctx.Request.Method == "GET" {
			var user User
			var id, erro = strconv.Atoi(ctx.Params[0])
			if erro == nil {
				fmt.Println("")
			}
			fmt.Println(id)
			collection := client.Database("trial").Collection("users")
			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			err := collection.FindOne(ctx, User{Id: id}).Decode(&user)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(user)
			//json.NewEncoder(ctx.ResponseWriter).Encode(user)
		}
		
	})
	
	
	app.Handle(`/posts/([^/]+)$`, func(ctx *Context) {
		if ctx.Request.Method == "GET" {
			for i:=0; i < len(Posts); i++ {
				if strconv.Itoa(Posts[i].Id) == ctx.Params[0] {
				json.NewEncoder(ctx.ResponseWriter).Encode(Posts[i])	
				}
			}
		}
	})
	
	
	app.Handle(`/users`, func(ctx *Context) {
		if ctx.Request.Method == "POST" {
			var name = ctx.Request.FormValue("Name")
			var email = ctx.Request.FormValue("Email")
			var password = ctx.Request.FormValue("Password")
			var tmp = User{Id: 3, Name : name, Email : email, Password : password}
			collection := client.Database("trial").Collection("users")
			tmp.Password = fmt.Sprint(md5.Sum([]byte(tmp.Password)))
			ctx2, _ := context.WithTimeout(context.Background(), 10*time.Second)
			result, _ := collection.InsertOne(ctx2, tmp)
			json.NewEncoder(ctx.ResponseWriter).Encode(result)
		}
	})
	
	
	app.Handle(`/posts`, func(ctx *Context) {
		if ctx.Request.Method == "POST" {
			var caption = ctx.Request.FormValue("Caption")
			var image_url = ctx.Request.FormValue("Image URL")
			var posted_time = ctx.Request.FormValue("Posted time")
			var tmp = Post{Id: 3, Caption : caption, Image_url : image_url, Posted_time : posted_time}
			collection := client.Database("trial").Collection("posts")
			ctx2, _ := context.WithTimeout(context.Background(), 10*time.Second)
			result, _ := collection.InsertOne(ctx2, tmp)
			json.NewEncoder(ctx.ResponseWriter).Encode(result)
		}
	})
	


	err2 := http.ListenAndServe(":9000", app)

	if err2 != nil {
		log.Fatalf("Could not start server: %s\n", err2.Error())
	}

}
