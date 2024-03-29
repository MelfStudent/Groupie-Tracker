package Internal

import (
	"fmt"
	"strings"
)

// FilterPassed Filters API data with the desired filters
func FilterPassed(artist Artist, filters FormValues) bool {
	var result = false
	if filters.MinCreationDate != 0 && filters.MaxCreationDate != 0 {
		if artist.CreationDate > filters.MinCreationDate && artist.CreationDate < filters.MaxCreationDate {
			result = true
		} else {
			result = false
		}
	}

	if filters.MinFirstAlbum != 0 && filters.MaxFirstAlbum != 0 {
		if filters.MinFirstAlbum < artist.CreationDate && filters.MaxFirstAlbum > artist.CreationDate {
			result = true
			fmt.Println("a")
		} else {
			result = false
		}
	}

	if filters.NumberMembers != nil {
		var number = [7]bool{false, false, false, false, false, false, false}
		for i := range filters.NumberMembers {
			if filters.NumberMembers[i] != 0 {
				if len(artist.Members) == filters.NumberMembers[i] {
					number[i] = true
				} else {
					number[i] = false
				}
			}
		}
		if number[0] == true || number[1] == true || number[2] == true || number[3] == true || number[4] == true || number[5] == true || number[6] == true {
			result = true
		} else {
			result = false
		}
	}

	if filters.LocationConcert != "" {
		for i := range artist.Locations {
			research := WordProcessing(filters.LocationConcert)
			for j := 0; j < len(research); j++ {
				fmt.Println(strings.ToLower(artist.Locations[i]), strings.ToLower(research[j]))
				if strings.Contains(strings.ToLower(artist.Locations[i]), strings.ToLower(research[j])) {
					result = true
				}
			}
		}
	}

	fmt.Println(result)
	return result
}

// WordProcessing removed all non-important spaces
func WordProcessing(research string) []string {
	var wordsList []string
	words := strings.Split(research, ",")
	for i := range words {
		deletes := strings.ReplaceAll(words[i], " ", "")
		wordsList = append(wordsList, deletes)
	}
	return wordsList
}
