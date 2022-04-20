package mahjong

import (
	"log"
	"mahjong/table"
	"testing"
)

func TestShanten(t *testing.T) {
	shantenMap := make(map[uint64]int)
	for tiles := range table.EnumTiles(2, 5, 8, 11) {
		shanten, results := Shanten(tiles)
		if shanten == -1 && len(results) == 0 {
			panic(shanten)
		}
		shantenMap[tiles] = shanten
		//for _, result := range results {
		//	log.Printf("%d [%b] Shanten %d Analysis %s", tiles, tiles, shanten, result.String())
		//}
		//log.Printf("%d [%b] Shanten %d", tiles, tiles, shanten)
	}

<<<<<<< HEAD
	check := []int{11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 18}
	shanten, ok := Suggest(shantenMap, check)

	log.Printf("shanten %d , ok %t", shanten, ok)
}

func Suggest(shantenMap map[uint64]int, tiles []int) (int, bool) {
	key := BuildKey(tiles)
	shanten, ok := shantenMap[key]
	if !ok {
		return 99, false
	}
	return shanten, true
=======
	log.Printf("%+v", shantenMap)
>>>>>>> 1353c8b588778e6f64c722dcc80b9de5329367de
}
