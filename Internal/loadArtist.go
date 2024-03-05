package Internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	LocationsUrl string   `json:"locations"`
	Locations    []string
	Latitude     []float64
	Longitude    []float64
	ConcertDates string `json:"concertDates"`
	Relations    string `json:"relations"`
}

type LocationsResponse struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	DatesURL  string   `json:"dates"`
}

type OpenCageResponse struct {
	Results []struct {
		Geometry struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"geometry"`
	} `json:"results"`
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

	for i := 0; i < len(Artists); i++ {
		response, err := http.Get(Artists[i].LocationsUrl)
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

		// Décodez le JSON dans une structure LocationsResponse
		var result LocationsResponse
		err = json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println("Erreur lors du décodage JSON:", err)
			return
		}

		Artists[i].Locations = result.Locations

		MapArtist(i)
	}
	//fmt.Println(Artists)
}

func MapArtist(index int) {
	apiKey := "4bb307b157f34868b7cc9acc4878e4f1"
	for i := 0; i < len(Artists[index].Locations); i++ {
		address := Artists[index].Locations[i]
		url := fmt.Sprintf("https://api.opencagedata.com/geocode/v1/json?q=%s&key=%s", address, apiKey)

		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		response, err := client.Get(url)
		if err != nil {
			fmt.Println("Erreur lors de la requête HTTP :", err)
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			fmt.Println("La requête a échoué avec le code de statut :", response.StatusCode)
			return
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Erreur lors de la lecture de la réponse :", err)
			return
		}

		var result OpenCageResponse
		err = json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println("Erreur lors de l'analyse JSON :", err)
			return
		}

		Artists[index].Latitude = append(Artists[index].Latitude, result.Results[0].Geometry.Lat)
		Artists[index].Longitude = append(Artists[index].Longitude, result.Results[0].Geometry.Lng)

	}
	fmt.Printf("Latitude: %f, Longitude: %f\n", Artists[index].Latitude, Artists[index].Longitude)
	/*
		if len(result.Results) > 0 {
			Artists[0].Latitude = result.Results[0].Geometry.Lat
			Artists[0].Longitude = result.Results[0].Geometry.Lng
			fmt.Printf("Latitude: %f, Longitude: %f\n", Artists[0].Latitude, Artists[0].Longitude)
		} else {
			fmt.Println("Aucun résultat trouvé.")
		}
	*/
}
