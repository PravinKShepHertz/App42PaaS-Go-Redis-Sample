package main

import (
	"fmt"
	"github.com/astaxie/goredis"
	"net/http"
	"log"
	"text/template"
)

type User struct {
	UserName string
}

var (
	client goredis.Client
)

func setupDB(){
	client.Addr = "127.0.0.1:6379"
	// client, err := DialURL("tcp://auth:password@127.0.0.1:6379/0?timeout=10s&maxidle=1")
}

func indexHanlder(w http.ResponseWriter, r *http.Request){
	fmt.Println("In index handler!")
    username,_ := client.Lrange("username", 0, 10000)
	users := []User{}
	for _, v := range username {
		user := User{}
		user.UserName = string(v)
		users = append(users, user)
	}
	t := template.New("index.html")
	t.ParseFiles("templates/index.html")
	t.Execute(w, users)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In new handler!")
	t := template.New("new.html")
	t.ParseFiles("templates/new.html")
	t.Execute(w, t)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In save handler!")
	username := r.FormValue("username")
 	client.Rpush("username", []byte(username))
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	setupDB()
	fmt.Println("Redis db connected: ", client)

	http.HandleFunc("/", indexHanlder)
	http.HandleFunc("/new/", newHandler)
	http.HandleFunc("/save/", saveHandler)
	fmt.Println("Listening on port 3000....")
	http.Handle("/public/css/", http.StripPrefix("/public/css/", http.FileServer(http.Dir("public/css"))))
	http.Handle("/public/images/", http.StripPrefix("/public/images/", http.FileServer(http.Dir("public/images"))))
	if err := http.ListenAndServe("0.0.0.0:3000", nil); err != nil {
		log.Fatalf("Error in listening:", err)	
	}
}
