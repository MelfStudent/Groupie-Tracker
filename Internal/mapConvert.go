package Internal

import (
	"errors"
	"strconv"
)

type FormValues struct {
	MinCreationDate int
	MaxCreationDate int
	MinFirstAlbum   int
	MaxFirstAlbum   int
	LocationConcert string
	NumberMembers   []int
}

func MapConvert(formValues map[string]string) (FormValues, error) {
	var values FormValues

	if formValues["creationDateSelectMin"] != "" {
		minCreationDateStr, ok := formValues["creationDateSelectMin"]
		if !ok {
			return values, errors.New("creationDateSelectMin is missing")
		}
		minCreationDate, err := strconv.Atoi(minCreationDateStr)
		if err != nil {
			return values, err
		}
		values.MinCreationDate = minCreationDate
		//fmt.Println(values.MinCreationDate)
	}

	if formValues["creationDateSelectMax"] != "" {
		maxCreationDateStr, ok := formValues["creationDateSelectMax"]
		if !ok {
			return values, errors.New("creationDateSelectMax is missing")
		}
		maxCreationDate, err := strconv.Atoi(maxCreationDateStr)
		if err != nil {
			return values, err
		}
		values.MaxCreationDate = maxCreationDate
		//fmt.Println(values.MaxCreationDate)
	}

	if formValues["firstAlbumSelectMin"] != "" {
		minFirstAlbumStr, ok := formValues["firstAlbumSelectMin"]
		if !ok {
			return values, errors.New("firstAlbumSelectMin is missing")
		}
		minFirstAlbum, err := strconv.Atoi(minFirstAlbumStr)
		if err != nil {
			return values, err
		}
		values.MinFirstAlbum = minFirstAlbum
		//fmt.Println(values.MinFirstAlbum)
	}

	if formValues["firstAlbumSelectMax"] != "" {
		maxFirstAlbumStr, ok := formValues["firstAlbumSelectMax"]
		if !ok {
			return values, errors.New("firstAlbumSelectMax is missing")
		}
		maxFirstAlbum, err := strconv.Atoi(maxFirstAlbumStr)
		if err != nil {
			return values, err
		}
		values.MaxFirstAlbum = maxFirstAlbum
		//fmt.Println(values.MaxFirstAlbum)
	}

	if formValues["locationConcert"] != "" {
		locationConcertStr, ok := formValues["locationConcert"]
		if !ok {
			return values, errors.New("locationConcert is missing")
		}
		values.LocationConcert = locationConcertStr
		//fmt.Println(values.NumberMembers)
	}

	if _, ok := formValues["number1"]; ok {
		values.NumberMembers = append(values.NumberMembers, 1)
	}
	if _, ok := formValues["number2"]; ok {
		values.NumberMembers = append(values.NumberMembers, 2)
	}
	if _, ok := formValues["number3"]; ok {
		values.NumberMembers = append(values.NumberMembers, 3)
	}
	if _, ok := formValues["number4"]; ok {
		values.NumberMembers = append(values.NumberMembers, 4)
	}
	if _, ok := formValues["number5"]; ok {
		values.NumberMembers = append(values.NumberMembers, 5)
	}
	if _, ok := formValues["number6"]; ok {
		values.NumberMembers = append(values.NumberMembers, 6)
	}
	if _, ok := formValues["number7"]; ok {
		values.NumberMembers = append(values.NumberMembers, 7)
	}

	return values, nil
}
