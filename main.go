package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/kr/pretty"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	
)

const (
    host     = "db"  // This is the service name in docker-compose
    port     = 5432
    user     = "postgres"
    password = "your_password"
    dbname   = "postgres"
)

type Post struct {
	gorm.Model
	Title string
	Content string
}

type PostHandler struct {
	ps []Post
}
func (ph PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("hello")
	log.Println(ph.ps)
	posts(ph.ps).Render(r.Context(), w)
}

var psc []Post

func postGetHandler(w http.ResponseWriter, r *http.Request) {
	posts(psc).Render(r.Context(), w)
}

func postPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("%# v", pretty.Formatter(r.Form))
	psc = append(psc, Post{Title: r.Form.Get("title"), Content: r.Form.Get("content")})
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	if r.Method == http.MethodPost {
		postPostHandler(w, r)
	}	
	postGetHandler(w, r)
} 

func main() {
	dsn := "host=localhost user=postgres password=your_password dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Post{})

	psc = append(psc, Post{Title: "Test1", Content: "Lorem ipsum"})
	// Use a template that doesn't take parameters.
	http.Handle("/", templ.Handler(home()))
	http.Handle("/add-post", templ.Handler(addPost()))

	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request){
		log.Printf("Handling Posts")
		postsHandler(w, r)
	})
	// Start the server.
	fmt.Println("listening on http://localhost:8000")
	if err := http.ListenAndServe("localhost:8000", nil); err != nil {
		log.Printf("error listening: %v", err)
	}
}
