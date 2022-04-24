// Package ibukisaar
package ibukisaar

import (
	"bytes"
	"github.com/RoaringBitmap/roaring/roaring64"
	"io/ioutil"
	"log"
	"mahjong/ibukisaar/analysis"
	"mahjong/ibukisaar/table"
	"sort"
	"sync"
	"time"
)

var (
	ShantenMap = sync.Map{}
	//ResultsMap    = make(map[uint64][]*analysis.Result)
)

func init() {
	now := time.Now()
	bitMap, err := Load("table.data")
	if err != nil {
		bitMap = table.EnumTiles(2, 5, 8, 11, 14)
		Store("table.data", bitMap)
	}
	iterator := bitMap.Iterator()

	for iterator.HasNext() {
		tiles := iterator.Next()
		go func() {
			info := analysis.Shanten(tiles)
			ShantenMap.Store(tiles, info)
		}()
	}
	log.Printf("time use %+v", time.Now().Sub(now))
}

func Store(path string, bitMap *roaring64.Bitmap) error {
	data, err := bitMap.ToBytes()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func Load(path string) (*roaring64.Bitmap, error) {
	bitMap := &roaring64.Bitmap{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	read := bytes.NewReader(data)
	_, err = bitMap.ReadFrom(read)
	if err != nil {
		return nil, err
	}
	return bitMap, nil
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

func Parse(tiles []int) (ret [][][]int) {
	sort.Ints(tiles)
	key := BuildKey(tiles)
	key = table.SortUInt64(key)
	data, ok := ShantenMap.Load(key)
	if !ok {
		return nil
	}
	info := data.(*analysis.Info)
	if info.Shanten != -1 {
		return nil
	}
	sort.Sort(sort.Reverse(sort.IntSlice(tiles)))
	tMap := make(map[int]int)
	sortKey := make([]int, 0, 14)
	for _, tile := range tiles {
		if _, ok := tMap[tile]; ok {
			tMap[tile]++
			continue
		}
		tMap[tile]++
		sortKey = append(sortKey, tile)
	}
	for _, result := range info.Results {
		r := [][]int{}
		mentsuToitsu := result.MentsuToitsu()
		toitsu := mentsuToitsu.Toitsu
		toitsu--
		tMap[sortKey[toitsu]] -= 2
		r = append(r, []int{sortKey[toitsu], sortKey[toitsu]})

		for _, koTsu := range mentsuToitsu.KoTsu {
			koTsu--
			tMap[sortKey[koTsu]] -= 3
			r = append(r, []int{sortKey[koTsu], sortKey[koTsu], sortKey[koTsu]})
		}

		for _, shuntsu := range mentsuToitsu.Shuntsu {
			shuntsu--
			tMap[sortKey[shuntsu]] -= 1
			tMap[sortKey[shuntsu+1]] -= 1
			tMap[sortKey[shuntsu+2]] -= 1
			r = append(r, []int{sortKey[shuntsu], sortKey[shuntsu+1], sortKey[shuntsu+2]})
		}
		ret = append(ret, r)
	}
	return
}
