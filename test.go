package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
)

var user_count = 0
var post_count = 1001

type User struct {
	Userid   int    ` json:"id"`
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}
type Posts struct {
	Userid    int
	Postid    int
	Caption   string
	imgurl    string
	Timestamp time.Time
}

func close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {

	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

func createuser(w http.ResponseWriter, r *http.Request) {
	
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./assets/user.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
        user_count++
		val := User{Userid: user_count, Name: r.FormValue("name"), Email: r.FormValue("Email"), Password: r.FormValue("Password")}
		insertResult, err := insert(client, ctx, "test1", "users", val)

		// handle the error
		if err != nil {
			panic(err)
		}
		fmt.Println(insertResult)
		fmt.Println(user_count)
		
	}
}

func showuser(w http.ResponseWriter, r *http.Request) {
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Print(id)
	var option interface{}
	var dataBase = "test1"
	var col = "users"
	collection := client.Database(dataBase).Collection(col)
	option = bson.D{{"_id", 0}}
	filterCursor, err := collection.Find(ctx, bson.M{"userid": id}, options.Find().SetProjection(option))
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = filterCursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	fmt.Println(results)
	

}

func createpost(w http.ResponseWriter, r *http.Request) {
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./assets/post.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()

		val := Posts{Userid:user_count-1, Postid: post_count, Caption: r.FormValue("caption"), imgurl: r.FormValue("img"), Timestamp: time.Now()}
		insertResult, err := insert(client, ctx, "test1", "posts", val)

		// handle the error
		if err != nil {
			panic(err)
		}
		fmt.Println(insertResult)
		post_count++
	}
}

func showpost(w http.ResponseWriter, r *http.Request) {
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Print(id)
	var option interface{}
	var dataBase = "test1"
	var col = "posts"
	collection := client.Database(dataBase).Collection(col)
	option = bson.D{{"_id", 0}}
	filterCursor, err := collection.Find(ctx, bson.M{"postid": id}, options.Find().SetProjection(option))
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = filterCursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	b,err:=json.Marshal(results)
	fmt.Println(string(b))

}

func listposts(w http.ResponseWriter, r *http.Request) {
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Print(id)
	var option interface{}
	var dataBase = "test1"
	var col = "posts"
	collection := client.Database(dataBase).Collection(col)
	option = bson.D{{"_id", 0}}
	filterCursor, err := collection.Find(ctx, bson.M{"userid": id}, options.Find().SetProjection(option))
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = filterCursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	b,err:=json.Marshal(results)
	fmt.Println(string(b))

}

func insert(client *mongo.Client, ctx context.Context, dataBase, col string, val interface{}) (*mongo.InsertOneResult, error) {

	collection := client.Database(dataBase).Collection(col)
	result, err := collection.InsertOne(ctx, val)
	return result, err

}

func main() {

	http.HandleFunc("/users", createuser)
	http.HandleFunc("/posts", createpost)
	http.HandleFunc("/user", showuser)
	http.HandleFunc("/post", showpost)
	http.HandleFunc("/posts/users",listposts)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
