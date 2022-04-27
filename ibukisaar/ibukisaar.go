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
	//ShantenMap = make(map[uint64]*analysis.Info)
	sortTiles = map[int]int{
		11: 0, 21: 9, 31: 18, 41: 27,
		12: 1, 22: 10, 32: 19, 42: 28,
		13: 2, 23: 11, 33: 20, 43: 29,
		14: 3, 24: 12, 34: 21, 44: 30,
		15: 4, 25: 13, 35: 22, 45: 31,
		16: 5, 26: 14, 36: 23, 46: 32,
		17: 6, 27: 15, 37: 24, 47: 33,
		18: 7, 28: 16, 38: 25,
		19: 8, 29: 17, 39: 26,
	}
)

func getIndex(tile int) int {
	return sortTiles[tile]
}

func init() {
	now := time.Now()
	bitMap, err := Load("table.data")
	if err != nil {
		bitMap = table.EnumTiles(2, 5, 8, 11, 14)
		Store("table.data", bitMap)
	}
	iterator := bitMap.Iterator()
	wg := sync.WaitGroup{}
	for iterator.HasNext() {
		tiles := iterator.Next()
		wg.Add(1)
		go func() {
			info := analysis.Shanten(tiles)
			//ShantenMap[tiles] = info
			ShantenMap.Store(tiles, info)
			wg.Done()
		}()
	}
	wg.Wait()
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

type Keys []*Key

type Key struct {
	Value uint64
	Bits  uint64
	Tiles [][]int
}

func newKey(value, bits uint64, tiles [][]int) *Key {
	flip := Flip(value, bits)
	for flip < value {
		value = flip
		reverse(tiles)
	}
	return &Key{
		Value: value,
		Bits:  bits,
		Tiles: tiles,
	}
}

func reverse(arr [][]int) {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}
func (items Keys) Reverse() {
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
}
func (items Keys) Len() int           { return len(items) }
func (items Keys) Less(i, j int) bool { return items[i].Value < items[j].Value }
func (items Keys) Swap(i, j int)      { items[i], items[j] = items[j], items[i] }

func Flip(value, length uint64) uint64 {
	temp := uint64(0b10 << (length - 2))
	maxShift := length - 4
	maxShift2 := length - 8
	for shift := uint64(0); shift < length; shift += 4 {
		temp |= ((value >> shift) & 0b0011) << (maxShift - shift)
		if shift < length-4 {
			temp |= ((value >> shift) & 0b1100) << (maxShift2 - shift)
		}
	}
	return temp
}

func Parse(tiles []int) (result Keys) {
	tilesCnt := [34][]int{}
	pointCnt := 0
	for _, tile := range tiles {
		index := getIndex(tile)
		if len(tilesCnt[index]) == 0 {
			pointCnt++
		}
		tilesCnt[index] = append(tilesCnt[index], tile)
	}
	temp := make([][]int, 0, pointCnt)
	for i := 0; i < 34; i++ {
		if len(tilesCnt[i]) == 0 {
			continue
		}
		temp = append(temp, tilesCnt[i])
	}

	var start, length, value, shift uint64
	diffIndex := 0
	prevIndex := -3
	for i := 0; i <= 34; i++ {
		if i%9 == 0 || i >= 27 {
			prevIndex = -3
		}
		if i < 34 && len(tilesCnt[i]) == 0 {
			continue
		}
		diffIndex = i - prevIndex - 1
		if diffIndex >= 2 {
			value |= (2 << shift) >> 2
			if value != 0 {
				result = append(result, newKey(value, shift, temp[start:start+length]))
				//result.Add(new Key(value, shift, new ArraySegment<List<Tile>>(template, start, length)));
				start += length
				length = 0
				value = 0
				shift = 0
			}
		} else {
			value |= uint64(diffIndex) << (shift - 2)
		}
		if i == 34 {
			continue
		}
		value |= uint64(len(tilesCnt[i])-1) << shift
		length++
		shift += 4
		prevIndex = i
	}
	sort.Sort(result)
	return result
}

func BuildKey(keys Keys) uint64 {
	value := uint64(0xF)
	for i := 0; i < keys.Len(); i++ {
		value = (value << keys[i].Bits) | keys[i].Value
	}
	return value
}

func Analysis(info *analysis.Info, keys Keys) [][][]int {
	keys.Reverse()
	temp := make([][]int, 0)
	for _, key := range keys {
		temp = append(temp, key.Tiles...)
	}
	groupsList := make([][][]int, 0)
	for i := 0; i < len(info.Results); i++ {
		indexes := make([]int, len(temp))
		var result = info.Results[i]
		groups := make([][]int, 0, 5)

		pairIndex := result.Pair - 1
		cur := indexes[pairIndex]
		indexes[pairIndex]++
		pairs := make([]int, 0, 2)
		pairs = append(pairs, temp[pairIndex][cur])
		cur = indexes[pairIndex]
		indexes[pairIndex]++
		pairs = append(pairs, temp[pairIndex][cur])
		groups = append(groups, pairs)
		groupIds := result.Groups

		for junkoIndex := uint(0); junkoIndex < uint(result.JunkoCount); junkoIndex, groupIds = junkoIndex+1, groupIds>>8 {
			index := int(groupIds & 0xFF)
			junko := make([]int, 0, 3)
			cur2 := indexes[index-2]
			indexes[index-2]++
			cur1 := indexes[index-1]
			indexes[index-1]++
			cur0 := indexes[index-0]
			indexes[index-0]++
			junko = append(junko, temp[index-2][cur2])
			junko = append(junko, temp[index-1][cur1])
			junko = append(junko, temp[index-0][cur0])
			groups = append(groups, junko)
		}

		for ; groupIds != 0; groupIds >>= 8 {
			index := (int)(groupIds&0xFF) - 1
			cur2 := indexes[index]
			indexes[index]++
			cur1 := indexes[index]
			indexes[index]++
			cur0 := indexes[index]
			indexes[index]++
			pung := make([]int, 0, 3)
			pung = append(pung, temp[index][cur2])
			pung = append(pung, temp[index][cur1])
			pung = append(pung, temp[index][cur0])
			groups = append(groups, pung)
		}
		groupsList = append(groupsList, groups)
	}
	return groupsList

}
