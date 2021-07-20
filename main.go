package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type user struct {
	Id       int
	UserName string
	Email    string
	First    string
	Last     string
	Password []byte
}
type post struct {
	Id       int
	UserId   int
	CategId  int
	Title    string
	Text     string
	Like     int
	DisLike  int
	UserName string
	Categ    string
}
type categ struct {
	Id    int
	Categ string
}
type comments struct {
	Id      int
	User    string
	UserId  int
	PostId  int
	Comment string
}
type emotion struct {
	Id       int
	likes    int
	dislikes int
	epost    post
	ecomment comments
}

var tmp *template.Template
var dbUsers = map[string]user{}      // user ID, user
var dbSessions = map[string]string{} // session ID, user ID
var idCateg int
var idPost int
var idUser int

func init() {
	tmp = template.Must(template.ParseGlob("assets/*.html"))

}
func main() {
	createTables()
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/theme", theme)
	mux.HandleFunc("/postes", postes)
	mux.HandleFunc("/comment", comment)
	mux.HandleFunc("/title", title)
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
