package mahjong

// EnumQuantity 对数量穷举
func EnumQuantity(sum int) []uint64 {
	ret := make([]uint64, 0)
	enumQuantity(0, 0, sum, &ret)
	return ret
}

func enumQuantity(value uint64, deep, sum int, ret *[]uint64) {
	if sum == 0 {
		*ret = append(*ret, value)
		return
	}

	for i := 1; i <= 4; i++ {
		if sum >= i {
			enumQuantity(value|(i-1)<<(deep*4), deep+1, sum-i, ret)
		}
	}
}

// EnumDistance 对距离穷举
func EnumDistance(value []int) []uint64 {
	ret := make([]uint64, 0)
	enumDistance(value, 0, len(quantity)-1, 0, &ret)
	return ret
}

func enumDistance(value uint64, deep, deepCnt int, index uint64, ret *[]uint64) {
	if deep >= deepCnt {
		value |= 0b11 << (deep*4 + 4 - 2)
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

func LenTiles(value uint64) int {
	shift := uint64(0)
	for value >> shift {

	}
}

func TilesEnum(sum int) []uint64 {
	ret := make([]uint64, 0)
	for _, quantity := range EnumQuantity(sum) {
		for _, distance := range EnumDistance(quantity) {
			ret = append(ret, distance)
		}
	}
	return ret
}
