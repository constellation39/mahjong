package uakochan

import (
	"go.uber.org/zap"
	"strings"
	"utils/logger"
)

func GetBakaze(id uint32) (ret string) {
	switch id {
	case 0:
		ret = "E"
	case 1:
		ret = "S"
	case 2:
		ret = "W"
	case 3:
		ret = "N"
	default:
		logger.Panic("GetBakaze not found", zap.Uint32("id", id))
	}
	return
}

func GetTile(t string) (ret string) {
	ret = t
	switch t {
	case "1z":
		ret = "E"
	case "2z":
		ret = "S"
	case "3z":
		ret = "W"
	case "4z":
		ret = "N"
	case "5z":
		ret = "P"
	case "6z":
		ret = "F"
	case "7z":
		ret = "C"
	case "0m":
		ret = "5mr"
	case "0p":
		ret = "5pr"
	case "0s":
		ret = "5sr"
	}
	return
}

func GetTiles(tiles []string) (ret []string) {
	ret = make([]string, len(tiles))
	for i, tile := range tiles {
		ret[i] = GetTile(tile)
	}
	return
}

func GetComb(comb []string) (ret string) {
	return strings.Join(comb, "|")
}
