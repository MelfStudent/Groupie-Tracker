package Internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time" // Importez le package time ici
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
	Coordinate   []DataCoordinate
	ConcertDates string `json:"concertDates"`
	Dates        [][]string
	Relations    string `json:"relations"`
}

type LocationsResponse struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	DatesURL  string   `json:"dates"`
}

type DataCoordinate struct {
	Latitude  float64
	Longitude float64
	Name      string
	Date      []string
}

type DatesResponse struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type OpenCageResponse struct {
	Results []struct {
		Geometry struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"geometry"`
	} `json:"results"`
}

type DateData struct {
	Date int64 `json:"date"`
}

var Artists []Artist

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
	err = json.Unmarshal(body, &Artists)
	if err != nil {
		fmt.Println("Erreur lors du décodage JSON:", err)
		return
	}

	for i := 0; i < len(Artists); i++ {
		response, err := http.Get(Artists[i].ConcertDates)
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
		var resultDate DatesResponse
		err = json.Unmarshal(body, &resultDate)
		if err != nil {
			fmt.Println("Erreur lors du décodage JSON:", err)
			return
		}

		Artists[i].Dates = [][]string{}
		var dates []string
		dates = append(dates, resultDate.Dates[0])

		for j := 1; j < len(resultDate.Dates); j++ {
			if string(resultDate.Dates[j][0]) == "*" {
				Artists[i].Dates = append(Artists[i].Dates, dates)
				dates = []string{}
				dates = append(dates, resultDate.Dates[j])
			} else {
				dates = append(dates, resultDate.Dates[j])
			}
		}
		Artists[i].Dates = append(Artists[i].Dates, dates)
	}

	times, err := os.Open("date.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer times.Close()

	// Décode le fichier JSON dans la structure DateData
	var dateInfo DateData
	decoder := json.NewDecoder(times)
	err = decoder.Decode(&dateInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	}
	if time.Now().Unix()-dateInfo.Date >= 86400 {
		for i := 0; i < len(Artists); i++ {
			MapArtist(i)
		}
		// Extrait uniquement les coordonnées des artistes
		var coordinates [][][2]float64
		for _, artist := range Artists {
			var coo [][2]float64
			for i := 0; i < len(artist.Coordinate); i++ {
				coo = append(coo, [2]float64{artist.Coordinate[i].Latitude, artist.Coordinate[i].Longitude})
			}
			coordinates = append(coordinates, coo)
		}

		// Ouvre le fichier en mode écriture et le tronque
		file, err := os.OpenFile("coordinate.json", os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Convertit les coordonnées en JSON
		jsonData, err := json.MarshalIndent(coordinates, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Écrit les coordonnées JSON dans le fichier
		_, err = file.Write(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Les coordonnées ont été remplacées dans le fichier coordinates.json")
		//fmt.Println(Artists)

		fileTime, err := os.OpenFile("date.json", os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer fileTime.Close()

		// Créer une nouvelle structure DateData avec la date actuelle
		newDateInfo := DateData{Date: time.Now().Unix()}

		// Convertir la nouvelle structure en JSON
		jsonDate, err := json.MarshalIndent(newDateInfo, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Écrire la nouvelle date JSON dans le fichier "time.json"
		_, err = fileTime.Write(jsonDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		file, err := os.Open("coordinate.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Décode le contenu JSON du fichier dans une structure de données
		var coordinates [][][2]float64
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&coordinates)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for i := 0; i < len(coordinates); i++ {
			Artists[i].Coordinate = []DataCoordinate{}
			for j := 0; j < len(coordinates[i]); j++ {
				Artists[i].Coordinate = append(Artists[i].Coordinate, DataCoordinate{})
				Artists[i].Coordinate[j].Latitude = coordinates[i][j][0]
				Artists[i].Coordinate[j].Longitude = coordinates[i][j][1]
			}
		}

		fmt.Println("Les coordonnées ont été mises à jour dans la struct Artist")
	}
	for i := 0; i < len(Artists); i++ {
		for j := 0; j < len(Artists[i].Coordinate); j++ {
			Artists[i].Coordinate[j].Name = Artists[i].Name
			Artists[i].Coordinate[j].Date = Artists[i].Dates[j]
		}
	}
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
		Artists[index].Coordinate = append(Artists[index].Coordinate, DataCoordinate{})
		Artists[index].Coordinate[i].Latitude = result.Results[0].Geometry.Lat
		Artists[index].Coordinate[i].Longitude = result.Results[0].Geometry.Lng

	}
	fmt.Printf("Latitude: %f, Longitude: %f\n", Artists[index].Coordinate)
}
