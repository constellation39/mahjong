package mahjong

import (
	"sort"
)

// 原理：https://zhuanlan.zhihu.com/p/31000381
//

// EnumQuantity 对数量穷举
func EnumQuantity(count int) (ret [][]int) {
	ret = make([][]int, 0)
	enumQuantity(make([]int, 0, count), count, &ret)
	return
}
func enumQuantity(stack []int, count int, ret *[][]int) {
	if count == 0 {
		temp := make([]int, len(stack))
		*ret = append(*ret, temp)
		copy(temp, stack)
		return
	}

	for i := 1; i <= 4; i++ {
		if count >= i {
			stack = append(stack, i)
			enumQuantity(stack, count-i, ret)
			stack = stack[:len(stack)-1]
		}
	}
}

func Build(stack []int) (ret uint64) {
	var shift uint64
	for i := 0; i < len(stack); i, shift = i+1, shift+4 {
		ret |= uint64(stack[i]-1) << shift
	}
	ret |= 0b1111 << shift
	return
}

// EnumDistance 对距离穷举
func EnumDistance(stack []int) []uint64 {
	ret := make([]uint64, 0)
	enumDistance(Build(stack)|0b1000<<((len(stack)-1)*4), 0, len(stack)-1, 0, &ret)
	return ret
}

func enumDistance(value uint64, deep, deepCnt int, index uint64, ret *[]uint64) {
	if deep >= deepCnt {
		*ret = append(*ret, value)
		return
	}
	for i := uint64(1); i <= 3; i++ {
		if i == 3 {
			enumDistance(value|(i-1)<<(deep*4+2), deep+1, deepCnt, 0, ret)
			continue
		}
		if index+i > 8 {
			continue
		}
		enumDistance(value|(i-1)<<(deep*4+2), deep+1, deepCnt, index+i, ret)
	}
}

type RebuildItems []*RebuildItem
type RebuildItem struct {
	Value  uint64
	Length uint64
}

func NewRebuildItem(value, length uint64) *RebuildItem {
	return &RebuildItem{
		Value:  Min(value, Reverse(value, length)),
		Length: length,
	}
}

func Min(a, b uint64) uint64 {
	if a > b {
		return b
	}
	return a
}

func Reverse(value, length uint64) uint64 {
	temp := uint64(0b10 << (length - 2))
	maxShift := length - 4
	maxShift2 := length - 8
	for shift := uint64(0); shift < length; shift++ {
		temp |= (value >> shift) & 0b0011 << (maxShift - shift)
		if shift < length-4 {
			temp |= (value >> shift) & 0b1100 << (maxShift2 - shift)
		}
	}
	return temp
}

func (item RebuildItems) Len() int           { return len(item) }
func (item RebuildItems) Less(i, j int) bool { return item[i].Value < item[j].Value }
func (item RebuildItems) Swap(i, j int)      { item[i], item[j] = item[j], item[i] }

func ReBuild(value uint64) uint64 {
	rebuildItems := make(RebuildItems, 0)

	var shift uint64
	for value>>shift != 0b1111 {
		if value>>shift&0b1000 == 0 {
			shift += 4
			continue
		}
		shift += 4
		rebuildItems = append(rebuildItems, NewRebuildItem(value&((1<<shift)-1), shift))
		value >>= shift
		shift = 0
	}
	sort.Sort(rebuildItems)
	if value != 0b1111 {
		panic(value)
	}
	for _, rebuildItem := range rebuildItems {
		value = value<<rebuildItem.Length | rebuildItem.Value
	}
	return value
}
func Valid(value uint64) bool {
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

func TilesEnum(sum int) []uint64 {
	ret := make([]uint64, 0)
	for _, quantity := range EnumQuantity(sum) {
		//log.Printf("len(%+v) == %d [%b] \n", quantity, len(quantity), Build(quantity))
		for _, distance := range EnumDistance(quantity) {
			if !Valid(distance) {
				continue
			}
			ret = append(ret, ReBuild(distance))
		}
	}
	return ret
}
