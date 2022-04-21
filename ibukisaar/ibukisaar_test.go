package ibukisaar

import (
	"fmt"
	"mahjong/ibukisaar/table"
	"sort"
	"testing"
)

func TestBuildKey(t *testing.T) {
	tiles := []int{11, 12}
	sort.Ints(tiles)
	key := BuildKey(tiles)
	key = table.SortUInt64(key)
	fmt.Printf("key %d \n", key)
	shanten, ok := ShantenMap.Load(key)
	fmt.Printf("shanten %d ok %t", shanten, ok)
}
