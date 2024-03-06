package Internal

import (
	"strings"
)

// FilterPassed Filters API data with the desired filters
func FilterPassed(artist Artist, filters FormValues) bool {
	if filters.MinCreationDate != 0 && filters.MaxCreationDate != 0 {
		if artist.CreationDate > filters.MinCreationDate && artist.CreationDate < filters.MaxCreationDate {
			return true
		}
	}

	if filters.MinFirstAlbum != 0 && filters.MaxFirstAlbum != 0 {
		if filters.MinFirstAlbum > filters.MinFirstAlbum && filters.MaxFirstAlbum < filters.MaxFirstAlbum {
			return true
		}
	}

	for i := range filters.NumberMembers {
		if filters.NumberMembers[i] != 0 {
			if len(artist.Members) == filters.NumberMembers[i] {
				return true
			}
		}
	}

	if filters.LocationConcert != "" {
		for i := range artist.Locations {
			research := WordProcessing(filters.LocationConcert)
			for j := 0; j < len(research); j++ {
				if strings.Contains(strings.ToLower(artist.Locations[i]), strings.ToLower(research[j])) {
					return true
				}
			}
		}
	}

	return false
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
