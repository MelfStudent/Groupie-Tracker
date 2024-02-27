package Internal

import (
	"fmt"
	"sort"
)

type FilteredArtist struct {
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

// ResultFilters Retrieves input filters and returns a structure with the new data
func ResultFilters(minDateStr int, maxDateStr int, _Artists []Artist) []FilteredArtist {
	var Result []FilteredArtist

	for i := 0; i < len(_Artists); i++ {
		if _Artists[i].CreationDate >= minDateStr && _Artists[i].CreationDate <= maxDateStr {
			fmt.Println(_Artists[i].ID, _Artists[i].CreationDate)
			filtered := FilteredArtist{
				ID:           _Artists[i].ID,
				Image:        _Artists[i].Image,
				Name:         _Artists[i].Name,
				Members:      _Artists[i].Members,
				CreationDate: _Artists[i].CreationDate,
				FirstAlbum:   _Artists[i].FirstAlbum,
				Locations:    _Artists[i].Locations,
				ConcertDates: _Artists[i].ConcertDates,
				Relations:    _Artists[i].Relations,
			}
			Result = append(Result, filtered)
		}
	}

	sort.Slice(Result, func(i, j int) bool {
		return Result[i].CreationDate < Result[j].CreationDate
	})

	return Result
}
