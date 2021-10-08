package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string        `json:"Name"`
	Email    string        `json:"Email"`
	Password string        `json:"Password"`
}
type Posts struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Caption   string
	imgurl    string
	Timestamp time.Time
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("users")
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("user.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()

		err = c.Insert(&User{Name: r.FormValue("name"), Email: r.FormValue("Email"), Password: r.FormValue("Password")})

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("username:", r.Form["name"])
		fmt.Println("password:", r.Form["Password"])
		fmt.Println("email:", r.Form["Email"])
		defer session.Close()
		t, _ := template.ParseFiles("usercreated.html")
		t.Execute(w, nil)
	}
}
func CreatePost(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("posts")
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		m, _ := template.ParseFiles("post.html")
		m.Execute(w, nil)

	} else {
		r.ParseForm()

		err = c.Insert(&Posts{Caption: r.FormValue("caption"), imgurl: r.FormValue("img"), Timestamp: time.Now()})

		if err != nil {
			log.Fatal(err)
		}

		defer session.Close()

	}
}
func main() {

	http.HandleFunc("/users", CreateUser)
	http.HandleFunc("/posts", CreatePost)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
