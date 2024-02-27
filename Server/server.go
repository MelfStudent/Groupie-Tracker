package Server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

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

	http.HandleFunc("/updateFilters", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Form data error", http.StatusBadRequest)
				return
			}

			minDateStr := r.Form.Get("dateSelectMin")
			maxDateStr := r.Form.Get("dateSelectMax")
			minDate, _ := strconv.Atoi(minDateStr)
			maxDate, _ := strconv.Atoi(maxDateStr)
			filteredArtists := Internal.ResultFilters(minDate, maxDate, Internal.Artists)
			err = tmpl.Execute(w, filteredArtists)
			if err != nil {
				return
			}
			return
		}
	})

	fmt.Println("Pour accéder à la page Web -> http://localhost:8080/")
	err1 := http.ListenAndServe(":8080", nil)
	if err1 != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
