package ibukisaar

import (
	"fmt"
	"mahjong/ibukisaar/table"
	"sort"
	"testing"
)

func TestBuildKey(t *testing.T) {
	tiles := []int{11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 19, 19, 19, 19}
	sort.Ints(tiles)
	key := BuildKey(tiles)
	key = table.SortUInt64(key)
	fmt.Printf("key %d %b \n", key, key)

	results, ok := ResultsMap.Load(key)
	fmt.Printf("Contains %+v(%t) results %+v ok %t \n", tiles, ShantenBitMap.Contains(key), results, ok)
	//log.Println(ShantenBitMap.String())
}
