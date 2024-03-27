package Internal

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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

func LoadArtist(w http.ResponseWriter) {

	apiURL := "https://groupietrackers.herokuapp.com/api/artists"

	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Error during HTTP request:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading HTTP response:", err)
		return
	}

	if response.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status:", response.Status)
		return
	}

	err = json.Unmarshal(body, &Artists)
	if err != nil {
		fmt.Println("Error while decoding JSON:", err)
		return
	}

	for i := 0; i < len(Artists); i++ {
		response, err := http.Get(Artists[i].ConcertDates)
		if err != nil {
			fmt.Println("Error during HTTP request:", err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(response.Body)

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading HTTP response:", err)
			return
		}

		if response.StatusCode != http.StatusOK {
			fmt.Println("Request failed with status:", response.Status)
			return
		}

		var resultDate DatesResponse
		err = json.Unmarshal(body, &resultDate)
		if err != nil {
			fmt.Println("Error while decoding JSON:", err)
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
	defer func(times *os.File) {
		err := times.Close()
		if err != nil {

		}
	}(times)

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
			fmt.Println("Error during HTTP request:", err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(response.Body)

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading HTTP response:", err)
			return
		}

		if response.StatusCode != http.StatusOK {
			fmt.Println("Request failed with status:", response.Status)
			return
		}

		var result LocationsResponse
		err = json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println("Error while decoding JSON:", err)
			return
		}

		Artists[i].Locations = result.Locations
	}
	if time.Now().Unix()-dateInfo.Date >= 86400 {
		for i := 0; i < len(Artists); i++ {
			MapArtist(i)
		}

		var coordinates [][][2]float64
		for _, artist := range Artists {
			var coo [][2]float64
			for i := 0; i < len(artist.Coordinate); i++ {
				coo = append(coo, [2]float64{artist.Coordinate[i].Latitude, artist.Coordinate[i].Longitude})
			}
			coordinates = append(coordinates, coo)
		}

		file, err := os.OpenFile("coordinate.json", os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)

		jsonData, err := json.MarshalIndent(coordinates, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = file.Write(jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Coordinates have been replaced in the file coordinates.json")
		//fmt.Println(Artists)

		fileTime, err := os.OpenFile("date.json", os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(fileTime *os.File) {
			err := fileTime.Close()
			if err != nil {

			}
		}(fileTime)

		newDateInfo := DateData{Date: time.Now().Unix()}

		jsonDate, err := json.MarshalIndent(newDateInfo, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)

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

		fmt.Println("The coordinates have been updated in the Artist struct")
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
			fmt.Println("\nError during HTTP request :", err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(response.Body)

		if response.StatusCode != http.StatusOK {
			fmt.Println("Request failed with status code :", response.StatusCode)
			return
		}

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("\nError reading response :", err)
			return
		}

		var result OpenCageResponse
		err = json.Unmarshal(body, &result)
		if err != nil {
			fmt.Println("Error parsing JSON :", err)
			return
		}
		Artists[index].Coordinate = append(Artists[index].Coordinate, DataCoordinate{})
		Artists[index].Coordinate[i].Latitude = result.Results[0].Geometry.Lat
		Artists[index].Coordinate[i].Longitude = result.Results[0].Geometry.Lng

	}
	fmt.Printf("Latitude: %f, Longitude: %f\n", Artists[index].Coordinate)
}
