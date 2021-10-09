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
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
)
//Global variables
var user_id = 0
var post_id = 1001
var option interface{}
var dataBase = "instadata"   //database name


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

//This function closes the MongoDB connection and cancels contexts and resources
func close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	defer cancel() //This cancels context

	defer func() {  //This closes the connection with MongoDB
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

//This function establishes connection with MongoDB 
func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {

	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}


//This function takes input from userform and parses the values and stores in the database.
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
        user_id++ // user_id is globally initialized to 0.
		          //Everytime this function is invoked user_id value is incremented.
		
		//This takes input from userform(HTML) and parses the values and stores in the database.
		
		hash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("Password")), bcrypt.DefaultCost)
		if err != nil {
			
			log.Fatal(err)
		}
		fmt.Println("Hash to store:", string(hash))
		    val := User{
			Userid: user_id, 
			Name: r.FormValue("name"), 
			Email: r.FormValue("Email"),
			Password: string(hash),
			}
		
		insertResult, err := insert(client, ctx, "instadata", "users", val)

		// handle the error
		if err != nil {
			panic(err)
		}
		fmt.Println(insertResult)
		fmt.Println(user_id)
		
	}
}


//showuser function
func showuser(w http.ResponseWriter, r *http.Request) {
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)

	id, err := strconv.Atoi(r.URL.Query().Get("id")) //This command taked id param from url and stores in variable id
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Print(id)
	var col = "users"           //collection name
	collection := client.Database(dataBase).Collection(col)
	option = bson.D{{"_id", 0}}  //  option remove objectid field from all documents
	filter, err := collection.Find(ctx, bson.M{"userid": id}, options.Find().SetProjection(option))
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = filter.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	fmt.Println(results)
	

}


//CreatePost function

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
//This takes input from postform(HTML) and parses the values and stores in the database.
		val := Posts{
			Userid:user_id, 
			Postid: post_id, 
			Caption: r.FormValue("caption"), 
			imgurl: r.FormValue("img"), 
			Timestamp: time.Now()}
		
		insertResult, err := insert(client, ctx, "instadata", "posts", val)//insert operation 

		// handle the error
		if err != nil {
			panic(err)
		}
		fmt.Println(insertResult)
		post_id++ // user_id is globally initialized to 0.
		          //Everytime this function is invoked user_id value is incremented.
	}
}


//Showpost function
func showpost(w http.ResponseWriter, r *http.Request) {
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))//This command taked id param from url and stores in variable id
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Print(id)
	var col = "posts"           //collection name
	collection := client.Database(dataBase).Collection(col)
	option = bson.D{{"_id", 0}} //  option remove objectid field from all documents
	filter, err := collection.Find(ctx, bson.M{"postid": id}, options.Find().SetProjection(option))
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = filter.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	b,err:=json.Marshal(results)
	fmt.Println(string(b))

}

//listposts function
func listposts(w http.ResponseWriter, r *http.Request) {
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)

	id, err := strconv.Atoi(r.URL.Query().Get("id"))//This command taked id param from url and stores in variable id
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Print(id)
	
	collection := client.Database(dataBase).Collection(col)
	option = bson.D{{"_id", 0}}//  option remove objectid field from all documents
	filter, err := collection.Find(ctx, bson.M{"userid": id}, options.Find().SetProjection(option))//Compares userid with id in the url 
	                                                                                               //and finds the posts of the particular user with that id
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = filter.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	b,err:=json.Marshal(results)
	fmt.Println(string(b))

}

//performs insert operation
func insert(client *mongo.Client, ctx context.Context, dataBase, col string, val interface{}) (*mongo.InsertOneResult, error) {

	collection := client.Database(dataBase).Collection(col)
	result, err := collection.InsertOne(ctx, val)
	return result, err

}


//routing
func main() {

	http.HandleFunc("/users", createuser)
	http.HandleFunc("/posts", createpost)
	http.HandleFunc("/user", showuser)
	http.HandleFunc("/post", showpost)
	http.HandleFunc("/posts/users",listposts)

	log.Fatal(http.ListenAndServe(":8000", nil))

}
