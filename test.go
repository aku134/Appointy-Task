package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Name     string             `json:"Name"`
	Email    string             `json:"Email"`
	Password string             `json:"Password"`
}
type Posts struct {
	ID        primitive.ObjectID `bson:"_id, omitempty" json:"id"`
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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("user.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()

		val := User{Name: r.FormValue("name"), Email: r.FormValue("Email"), Password: r.FormValue("Password")}
		insertOneResult, err := insertOne(client, ctx, "test1", "users", val)

		// handle the error
		if err != nil {
			panic(err)
		}
		fmt.Println(insertOneResult)
	}
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("post.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()

		val := Posts{Caption: r.FormValue("caption"), imgurl: r.FormValue("img"), Timestamp: time.Now()}
		insertOneResult, err := insertOne(client, ctx, "test1", "posts", val)

		// handle the error
		if err != nil {
			panic(err)
		}
		fmt.Println(insertOneResult)
	}
}

func insertOne(client *mongo.Client, ctx context.Context, dataBase, col string, val interface{}) (*mongo.InsertOneResult, error) {

	collection := client.Database(dataBase).Collection(col)
	result, err := collection.InsertOne(ctx, val)
	return result, err
}

func main() {

	http.HandleFunc("/users", CreateUser)
	http.HandleFunc("/posts", CreatePost)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
