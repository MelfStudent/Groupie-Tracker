package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

var tmpl *template.Template
var artist []Artist

func StartServer() {
	var err error

	tmpl, err = template.New("index").ParseFiles("web/HTML/index.html")
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.Dir(wd + "\\web"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			LoadArtist(w, r)
			tmpl.Execute(w, artist)
		} else {
			fileServer.ServeHTTP(w, r)
		}
	})

	fmt.Println("Pour accéder à la page web -> http://localhost:8080/")
	err1 := http.ListenAndServe(":8080", nil)
	if err1 != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func LoadArtist(w http.ResponseWriter, r *http.Request) {

	apiURL := "https://groupietrackers.herokuapp.com/api/artists"

	// Effectuez une requête HTTP GET
	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Erreur lors de la requête HTTP:", err)
		return
	}
	defer response.Body.Close()

	// Lisez le corps de la réponse
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la réponse HTTP:", err)
		return
	}

	// Vérifiez si la requête a réussi (statut 200 OK)
	if response.StatusCode != http.StatusOK {
		fmt.Println("La requête a échoué avec le statut:", response.Status)
		return
	}

	// Décodez le JSON dans une structure Artist
	err = json.Unmarshal(body, &artist)
	if err != nil {
		fmt.Println("Erreur lors du décodage JSON:", err)
		return
	}

	//fmt.Println(artist)
}
