package table

import (
	"log"
	"sort"
	"time"
)

// enumQuantity 对数量穷举
func enumQuantity(count int) (ret [][]int) {
	ret = make([][]int, 0)
	quantitys(make([]int, 0, count), count, &ret)
	return
}
func quantitys(stack []int, count int, ret *[][]int) {
	if count == 0 {
		temp := make([]int, len(stack))
		*ret = append(*ret, temp)
		copy(temp, stack)
		return
	}

	for i := 1; i <= 4; i++ {
		if count >= i {
			stack = append(stack, i)
			quantitys(stack, count-i, ret)
			stack = stack[:len(stack)-1]
		}
	}
}

func buildUint64(stack []int) (ret uint64) {
	var shift uint64
	for i := 0; i < len(stack); i, shift = i+1, shift+4 {
		ret |= uint64(stack[i]-1) << shift
	}
	ret |= 0b1111 << shift
	return
}

// enumDistance 对距离穷举
func enumDistance(stack []int) []uint64 {
	ret := make([]uint64, 0)
	distances(buildUint64(stack)|0b1000<<((len(stack)-1)*4), 0, len(stack)-1, 0, &ret)
	return ret
}

func distances(value uint64, deep, deepCnt int, index uint64, ret *[]uint64) {
	if deep >= deepCnt {
		*ret = append(*ret, value)
		return
	}
	for i := uint64(1); i <= 3; i++ {
		if i == 3 {
			distances(value|(i-1)<<(deep*4+2), deep+1, deepCnt, 0, ret)
			continue
		}
		if index+i > 8 {
			continue
		}
		distances(value|(i-1)<<(deep*4+2), deep+1, deepCnt, index+i, ret)
	}
}

type rebuildItems []*rebuildItem
type rebuildItem struct {
	Value  uint64
	Length uint64
}

func newRebuildItem(value, length uint64) *rebuildItem {
	return &rebuildItem{
		Value:  min(value, reverse(value, length)),
		Length: length,
	}
}

func min(a, b uint64) uint64 {
	if a > b {
		return b
	}
	return a
}

func reverse(value, length uint64) uint64 {
	temp := uint64(0b10 << (length - 2))
	maxShift := length - 4
	maxShift2 := length - 8
	for shift := uint64(0); shift < length; shift += 4 {
		temp |= (value >> shift) & 0b0011 << (maxShift - shift)
		if shift < length-4 {
			temp |= (value >> shift) & 0b1100 << (maxShift2 - shift)
		}
	}
	return temp
}

func (item rebuildItems) Len() int           { return len(item) }
func (item rebuildItems) Less(i, j int) bool { return item[i].Value < item[j].Value }
func (item rebuildItems) Swap(i, j int)      { item[i], item[j] = item[j], item[i] }

func sortUInt64(value uint64) uint64 {
	rebuildItems := make(rebuildItems, 0)

	var shift uint64
	for value>>shift != 0b1111 {
		if value>>shift&0b1000 == 0 {
			shift += 4
			continue
		}
		shift += 4
		rebuildItems = append(rebuildItems, newRebuildItem(value&((1<<shift)-1), shift))
		value >>= shift
		shift = 0
	}
	sort.Sort(rebuildItems)
	if value != 0b1111 {
		panic(value)
	}
	for _, rebuildItem := range rebuildItems {
		value = (value << rebuildItem.Length) | rebuildItem.Value
	}
	return value
}
func valid(value uint64) bool {
	var level, continuous, tempContinuous uint64
	tempContinuous = 1
	for shift := uint64(0); (value >> shift) != 0xF; shift += 4 {
		singleContinuous := (value >> (shift + 2)) & 3
		if singleContinuous < 2 {
			if level >= 3 {
				return false
			}
			tempContinuous += singleContinuous + 1
		} else if level < 3 {
			if continuous+2+tempContinuous <= 9 {
				continuous += 2 + tempContinuous
			} else {
				continuous = 0
				tempContinuous = 1
				level++
			}
		}
	}
	return true
}

func EnumTiles(values ...int) map[uint64]struct{} {
	set := make(map[uint64]struct{})
	for _, value := range values {
		now := time.Now()
		EnumValue(value, set)
		log.Printf("%d cnt %d use time %s", value, len(set), time.Now().Sub(now))
	}
	return set
}

func EnumValue(num int, set map[uint64]struct{}) {
	for _, quantity := range enumQuantity(num) {
		for _, distance := range enumDistance(quantity) {
			//注释则不过滤无效牌型
			if !valid(distance) {
				continue
			}
			distance = sortUInt64(distance)
			set[distance] = struct{}{}
		}
	}
}
