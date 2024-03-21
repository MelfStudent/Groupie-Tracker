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

var tmpl *template.Template
var tmpl_map *template.Template
var tmpl_groupe *template.Template

func StartServer() {
	var err error

	tmpl, err = template.New("index").ParseFiles("Web/HTML/index.html")
	if err != nil {
		panic(err)
	}

	tmpl_map, err = template.New("map").ParseFiles("Web/HTML/map.html")
	if err != nil {
		panic(err)
	}

	tmpl_groupe, err = template.New("groupe").ParseFiles("Web/HTML/groupe.html")
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
			Internal.LoadArtist(w, r)
			//Internal.MapArtist()
			tmpl.Execute(w, Internal.Artists)
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})

	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/map" {
			tmpl_map.Execute(w, Internal.Artists)
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})

	http.HandleFunc("/group/", func(w http.ResponseWriter, r *http.Request) {
		// Extraire groupID de l'URL
		groupID := strings.TrimPrefix(r.URL.Path, "/group/")
		// Trouver le groupe correspondant dans Groups
		for _, group := range Internal.Artists {
			if strconv.Itoa(group.ID) == groupID {
				// Afficher le template avec les détails du groupe
				tmpl_groupe.ExecuteTemplate(w, "groupe", group)
				return
			}
		}
		// Gérer le cas où le groupe n'est pas trouvé
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
