package Internal

import (
	"math/rand"
	"time"
)

// SelectRandomGroups returns a list of 5 groups randomly from the Artists structure
func SelectRandomGroups(groups []Artist) []map[string]interface{} {
	rand.Seed(time.Now().UnixNano())

	selectedGroups := make([]map[string]interface{}, 5)

	selectedIndices := rand.Perm(len(groups))[:5]

	for i, index := range selectedIndices {
		selectedGroups[i] = map[string]interface{}{
			"id":    groups[index].ID,
			"image": groups[index].Image,
		}
	}

	return selectedGroups
}
