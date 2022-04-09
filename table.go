package mahjong

import (
	"sort"
)

// QuantityEnum 对数量穷举
func QuantityEnum(sum int) [][]int {
	ret := make([][]int, 0)
	quantityEnum(make([]int, 0, sum), sum, &ret)
	return ret
}

func quantityEnum(stack []int, sum int, ret *[][]int) {
	if sum == 0 {
		temp := make([]int, len(stack))
		*ret = append(*ret, temp)
		copy(temp, stack)
		return
	}

	for i := 1; i <= 4; i++ {
		if sum >= i {
			stack = append(stack, i)
			quantityEnum(stack, sum-i, ret)
			stack = stack[:len(stack)-1]
		}
	}
}

// DistanceEnum 对距离穷举
func DistanceEnum(quantity []int) []uint64 {
	ret := make([]uint64, 0)
	distanceEnum(Build(quantity)|0b1000<<((len(quantity)-1)*4), 0, len(quantity)-1, 0, &ret)
	return ret
}

func distanceEnum(value uint64, deep, deepCnt int, index uint64, ret *[]uint64) {
	if deep >= deepCnt {
		*ret = append(*ret, value)
		return
	}
	for i := uint64(1); i <= 3; i++ {
		if i == 3 {
			distanceEnum(value|(i-1)<<(deep*4+2), deep+1, deepCnt, 0, ret)
			continue
		}
		if index+i > 8 {
			continue
		}
		distanceEnum(value|(i-1)<<(deep*4+2), deep+1, deepCnt, index+i, ret)
	}
}

// Build 对数组以指定方式储存
func Build(quantity []int) uint64 {
	var ret, shift uint64
	for i := 0; i < len(quantity); i, shift = i+1, shift+4 {
		ret |= (uint64(quantity[i]) - 1) << shift
	}
	ret |= 0b1111 << shift
	return ret
}

func Valid(value uint64) bool {
	var continuous, level, tempContinuous uint64
	tempContinuous = 1
	for shift := 0; (value >> shift) != 0b1111; shift += 4 {
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

type RebuildItem struct {
	Value  uint64
	Length uint64
}

type RebuildItems []*RebuildItem

func (items RebuildItems) Len() int {
	return len(items)
}

func (items RebuildItems) Less(i, j int) bool {
	return items[i].Value > items[j].Value
}

func (items RebuildItems) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func NewRebuildItem(value uint64, length uint64) *RebuildItem {
	return &RebuildItem{
		Value:  Min(value, Flip(value, length)),
		Length: length,
	}
}

func Min(a, b uint64) uint64 {
	if a > b {
		return b
	}
	return a
}

func Flip(value uint64, length uint64) uint64 {
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

func Rebuild(value uint64) uint64 {
	var shift uint64
	var cache RebuildItems
	for value != 0b1111 {
		if ((value >> shift) & 0b1000) == 0 {
			shift += 4
			continue
		}
		shift += 4
		cache = append(cache, NewRebuildItem(value&((1<<shift)-1), shift))
		value >>= shift
		shift = 0
	}
	sort.Sort(cache)
	for i := 0; i < len(cache); i++ {
		value = (value << cache[i].Length) | cache[i].Value
	}
	return value
}

func TilesEnum(sum int) []uint64 {
	ret := make([]uint64, 0)
	var cnt int
	quantitys := QuantityEnum(sum)
	for _, quantity := range quantitys {
		distances := DistanceEnum(quantity)
		cnt += len(distances)
		for _, tiles := range distances {
			if !Valid(tiles) {
				continue
			}
			ret = append(ret, Rebuild(tiles))
		}
	}
	return ret
}

func GetInfo(values uint64) {
	var shift, tiles uint64
	var remCount int
	for (values>>shift)&0b1111 != 0b1111 {
		continuous := values >> (shift + 2) & 0b11
		cnt := (values >> shift & 0b11) + 1
		remCount += int(cnt)
		tiles |= (continuous*5 + cnt) << shift
		shift += 4
	}
	if RonCutPair(tiles, remCount) {
		//log.Printf("%d True", values)
		Analysis(tiles)
	}
}

func Get(value, shift uint64, continuous, singleCount *uint64) {
	meta := value >> shift & 0b1111
	*continuous = meta / 5
	*singleCount = meta % 5
}

func Set(value, shift uint64, distance, cnt *uint64) uint64 {
	return (value & ^(0b1111 << shift)) | ((*distance*5 + *cnt) << shift)
}

func RonCutPair(value uint64, cnt int) bool {
	for shift := uint64(0); (value >> shift) != 0; shift += 4 {
		var continuous, singleCount uint64
		Get(value, shift, &continuous, &singleCount)
		if singleCount < 2 {
			continue
		}
		temp := singleCount - 2
		if RonCutPung(Set(value, shift, &continuous, &temp), 0, cnt-2) {
			return true
		}
	}
	return false
}

func RonCutPung(value, shift uint64, cnt int) bool {
	if cnt == 0 {
		return true
	}
	for (value>>shift) != 0 && ((value>>shift)&0xF)%5 == 0 {
		shift += 4
	}
	if (value >> shift) == 0 {
		return RonCutJunko(value, 0, cnt)
	}
	var continuous, singleCount uint64
	Get(value, shift, &continuous, &singleCount)
	if singleCount >= 3 {
		temp := singleCount - 3
		if RonCutPung(Set(value, shift, &continuous, &temp), shift, cnt-3) {
			return true
		}
	}
	return RonCutPung(value, shift+4, cnt)
}

func RonCutJunko(value, shift uint64, cnt int) bool {
	if cnt == 0 {
		return true
	}

	for (value>>shift) != 0 && ((value>>shift)&0xF)%5 == 0 {
		shift += 4
	}
	var continuous1, singleCount1 uint64
	Get(value, shift, &continuous1, &singleCount1)
	if continuous1 == 0 {
		var continuous2, singleCount2 uint64
		Get(value, shift+4, &continuous2, &singleCount2)
		if continuous2 == 0 && singleCount2 > 0 {
			var continuous3, singleCount3 uint64
			Get(value, shift+8, &continuous3, &singleCount3)
			if singleCount3 > 0 {
				temp1 := singleCount1 - 1
				var valueT = Set(value, shift, &continuous1, &temp1)
				temp2 := singleCount2 - 1
				valueT = Set(valueT, shift+4, &continuous2, &temp2)
				temp3 := singleCount3 - 1
				valueT = Set(valueT, shift+8, &continuous3, &temp3)
				return RonCutJunko(valueT, shift, cnt-3)
			}
		}
	}
	return false
}

func Analysis(value uint64) {
	ret := make([][]int, 0)
	AnalysisCut2(value, &ret)
}

func AnalysisCut2(value uint64, ret *[][]int) {
	for shift := uint64(0); (value >> shift) != 0; shift += 4 {
		var continuous, singleCount uint64
		Get(value, shift, &continuous, &singleCount)
		if singleCount >= 2 {
			temp := singleCount - 2
			AnalysisCut3(Set(value, shift, &continuous, &temp), 0, (shift>>2)+1, 0, 0, 0, 0, ret)
		}
	}
}

func AnalysisCut3(value, shift, pair, junkoCount, junkos, pungCount, pungs uint64, ret *[][]int) {
	for (value>>shift) != 0 && ((value>>shift)&0xF)%5 == 0 {
		shift += 4
	}
	if value>>shift == 0 {
		return
	}
	var continuous1, singleCount1 uint64
	Get(value, shift, &continuous1, &singleCount1)
	if singleCount1 >= 3 {
		temp1 := singleCount1 - 3
		AnalysisCut3(Set(value, shift, &continuous1, &temp1), shift, pair, junkoCount, junkos, pungCount+1, (pungs<<8)|((shift>>2)+1), ret)
	}
}
