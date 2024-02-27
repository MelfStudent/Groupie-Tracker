package Internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

var Artists []Artist

func LoadArtist() {

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
	err = json.Unmarshal(body, &Artists)
	if err != nil {
		fmt.Println("Erreur lors du décodage JSON:", err)
		return
	}

	//fmt.Println(Artists)
}
