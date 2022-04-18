package mahjong

import (
	"log"
	"mahjong/table"
	"testing"
)

func TestShanten(t *testing.T) {
	shantenMap := make(map[int]uint64)
	for tiles := range table.EnumTiles(2, 5, 8, 11, 14) {
		shanten, results := Shanten(tiles)
		if shanten == -1 && len(results) == 0 {
			panic(shanten)
		}
		shantenMap[shanten]++
		//for _, result := range results {
		//	log.Printf("%d [%b] Shanten %d Analysis %s", tiles, tiles, shanten, result.String())
		//}
		//log.Printf("%d [%b] Shanten %d", tiles, tiles, shanten)
	}

	log.Printf("%+v", shantenMap)
}
