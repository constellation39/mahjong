package main

func GetTile(t string) (ret string) {
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
	return t
}

func GetTiles(tiles []string) (ret []string) {
	ret = make([]string, len(tiles))
	for i, tile := range tiles {
		ret[i] = GetTile(tile)
	}
	return
}
