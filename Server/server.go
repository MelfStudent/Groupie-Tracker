package Server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"Groupie-Tracker/Internal"
)

var tmplFilters *template.Template
var tmplHomepage *template.Template
var tmplMap *template.Template
var tmplGroup *template.Template

func StartServer() {
	var err error

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	tmplFilters, err = template.New("index").ParseFiles("template/index.html")
	if err != nil {
		panic(err)
	}

	tmplHomepage, err = template.New("homepage").ParseFiles("template/homepage.html")
	if err != nil {
		panic(err)
	}

	tmplMap, err = template.New("map").ParseFiles("template/map.html")
	if err != nil {
		panic(err)
	}

	tmplGroup, err = template.New("groupe").ParseFiles("template/group.html")
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.Dir(wd + "/template"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			Internal.LoadArtist(w)
			err := tmplHomepage.Execute(w, Internal.SelectRandomGroups(Internal.Artists))
			if err != nil {
				return
			}
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})

	http.HandleFunc("/filters", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/filters" {
			//Internal.LoadArtist(w)
			err := tmplFilters.Execute(w, Internal.Artists)
			if err != nil {
				return
			}
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})

	http.HandleFunc("/homepage", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/homepage" {
			//Internal.LoadArtist(w)
			err := tmplHomepage.Execute(w, Internal.SelectRandomGroups(Internal.Artists))
			if err != nil {
				return
			}
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})

	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/map" {
			err := tmplMap.Execute(w, Internal.Artists)
			if err != nil {
				return
			}
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})

	http.HandleFunc("/group/", func(w http.ResponseWriter, r *http.Request) {
		groupID := strings.TrimPrefix(r.URL.Path, "/group/")
		for _, group := range Internal.Artists {
			if strconv.Itoa(group.ID) == groupID {
				err := tmplGroup.ExecuteTemplate(w, "groupe", group)
				if err != nil {
					return
				}
				return
			}
		}
		http.NotFound(w, r)
	})

	http.HandleFunc("/updateFilters", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Form data error", http.StatusBadRequest)
				return
			}

			formValues := make(map[string]string)
			for key := range r.Form {
				formValues[key] = r.Form.Get(key)
			}
			filteredArtists := Internal.ResultFilters(formValues, Internal.Artists)
			err = tmplFilters.Execute(w, filteredArtists)
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
