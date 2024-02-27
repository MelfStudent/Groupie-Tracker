package Server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"Groupie-Tracker/Internal"
)

var tmpl *template.Template

func StartServer() {
	var err error

	tmpl, err = template.New("index").ParseFiles("Web/HTML/index.html")
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.Dir(wd + "/Web"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			Internal.LoadArtist()
			tmpl.Execute(w, Internal.Artists)
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})

	fmt.Println("Pour accéder à la page Web -> http://localhost:8080/")
	err1 := http.ListenAndServe(":8080", nil)
	if err1 != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
