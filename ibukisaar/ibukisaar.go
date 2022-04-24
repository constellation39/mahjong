// Package ibukisaar
package ibukisaar

import (
	"bytes"
	"github.com/RoaringBitmap/roaring/roaring64"
	"io/ioutil"
	"log"
	"mahjong/ibukisaar/analysis"
	"mahjong/ibukisaar/table"
	"sync"
	"time"
)

var (
	ShantenBitMap = &roaring64.Bitmap{}
	ShantenMap    = sync.Map{}
	ResultsMap    = sync.Map{}
	//ResultsMap    = make(map[uint64][]*analysis.Result)
)

func init() {
	now := time.Now()
	if err := Load("table.data"); err != nil {
		table.EnumTiles(2, 5, 8, 11, 14).Range(func(key, _ interface{}) bool {
			tiles := key.(uint64)
			//shanten, results := analysis.Shanten(tiles)
			//_ = results
			//if shanten == -1 && len(results) == 0 {
			//	panic(shanten)
			//}
			ShantenBitMap.Add(tiles)
			//ShantenMap[tiles] = shanten
			//ResultsMap[tiles] = results
			return true
		})
		Store("table.data")
	}

	iterator := ShantenBitMap.Iterator()

	count := 0
	wg := sync.WaitGroup{}
	for iterator.HasNext() {
		count++
		wg.Add(1)
		tiles := iterator.Next()
		go func() {

			shanten, results := analysis.Shanten(tiles)
			_ = results

			if shanten == -1 && len(results) == 0 {
				panic(shanten)
			}
			ShantenMap.Store(tiles, shanten)
			ResultsMap.Store(tiles, results)
			wg.Done()
		}()
	}
	wg.Wait()

	log.Printf("%d use time %s", count, time.Now().Sub(now))
}

func Store(path string) error {
	data, err := ShantenBitMap.ToBytes()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	read := bytes.NewReader(data)
	_, err = ShantenBitMap.ReadFrom(read)
	if err != nil {
		return err
	}
	return nil
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
		key |= (dis - 1) << (index*4 + 2)
	}

	key |= 0b1111 << (len(set) * 4)
	return key
}
