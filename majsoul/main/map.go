package main

import (
	"go.uber.org/zap"
	"strings"
	"utils/logger"
)

func GetSoulTile(t string) (ret string) {
	ret = t
	switch t {
	case "E":
		ret = "1z"
	case "S":
		ret = "2z"
	case "W":
		ret = "3z"
	case "N":
		ret = "4z"
	case "P":
		ret = "5z"
	case "F":
		ret = "6z"
	case "C":
		ret = "7z"
	case "5mr":
		ret = "0m"
	case "5pr":
		ret = "0p"
	case "5sr":
		ret = "0s"
	}
	return
}

func GetSoulTiles(tiles []string) (ret []string) {
	ret = make([]string, len(tiles))
	for i, tile := range tiles {
		ret[i] = GetSoulTile(tile)
	}
	return
}

func GetUAkochanBakaze(id uint32) (ret string) {
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

func GetUAkochanTile(t string) (ret string) {
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

func GetUAkochanTiles(tiles []string) (ret []string) {
	ret = make([]string, len(tiles))
	for i, tile := range tiles {
		ret[i] = GetUAkochanTile(tile)
	}
	return
}

func GetSoulComb(comb []string) (ret string) {
	return strings.Join(GetSoulTiles(comb), "|")
}
