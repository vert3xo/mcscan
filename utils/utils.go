package utils

import "github.com/vert3xo/mcscan/types"

func RemoveDuplicatesFromPlayersList(playersList []types.Player) []types.Player {
	occurred := map[types.Player]bool{}
	result := []types.Player{}

	for _, player := range playersList {
		if !occurred[player] {
			occurred[player] = true

			result = append(result, player)
		}
	}
	return result
}