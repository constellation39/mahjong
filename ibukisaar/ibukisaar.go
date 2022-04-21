// Package ibukisaar
package ibukisaar

import (
	"mahjong/ibukisaar/analysis"
	"mahjong/ibukisaar/table"
	"sync"
)

var (
	//ShantenMap = make(map[uint64]int)
	ShantenMap = sync.Map{}
	//ResultsMap = make(map[uint64][]*analysis.Result)
	ResultsMap = sync.Map{}
)

func init() {
	wg := sync.WaitGroup{}
	table.EnumTiles(2, 5, 8, 11).Range(func(key, _ interface{}) bool {
		tiles := key.(uint64)
		println(tiles)
		wg.Add(1)
		go func() {
			shanten, results := analysis.Shanten(tiles)
			if shanten == -1 && len(results) == 0 {
				panic(shanten)
			}
			ShantenMap.Store(tiles, shanten)
			ResultsMap.Store(tiles, results)
			wg.Done()
		}()
		return true
	})
	wg.Wait()
}

func BuildKey(tiles []int) (key uint64) {
	set := make([]int, 0)
	cnt := make(map[int]int, 0)
	for _, tile := range tiles {
		if _, ok := cnt[tile]; ok {
			cnt[tile]++
			continue
		}
		cnt[tile]++
		set = append(set, tile)
	}

	for index, tile := range set {
		key |= uint64(cnt[tile]-1) << (index * 4)
		if index+1 >= len(set) {
			key |= 0b10 << (index*4 + 2)
			continue
		}
		dis := uint64(set[index+1] - set[index])
		if dis > 3 {
			dis = 3
		}
		key |= dis << (index*4 + 2)
	}

	key |= 0b1111 << (len(set) * 4)
	return key
}
