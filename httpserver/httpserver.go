package httpserver

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var temp *template.Template

func init() {
	temp = template.Must(template.ParseGlob("public/html/*.html"))
}

func RunServer() {
	fileserver := http.FileServer(http.Dir("./public"))
	http.Handle("/", fileserver)
	http.HandleFunc("/home", home)
	http.HandleFunc("/contact", contact)
	fmt.Printf("Http Server is running...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal((err))
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" {
		http.Error(w, "Not Found", 404)
		return
	}
	temp.ExecuteTemplate(w, "index.html", nil)
}
func contact(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/contact" {
		http.Error(w, "Not Found", 404)
		return
	}
	temp.ExecuteTemplate(w, "contact.html", nil)
}
