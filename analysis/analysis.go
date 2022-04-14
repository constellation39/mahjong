package analysis

import (
	"fmt"
	"mahjong/utils"
)

type SyantenArgs struct {
	N               uint64
	Result          uint64
	MaxUseTileCount uint64
}

type Result struct {
	JunkoCount uint64
	Pair       uint64
	Groups     uint64
}

func New(pair, junkoCount, junkos, pungCount, pungs uint64) *Result {
	junkos = Sort(junkos, junkoCount)
	pungs = Sort(pungs, pungCount)
	return &Result{
		Pair:       pair,
		JunkoCount: junkoCount,
		Groups:     pungs<<(junkoCount*8) | junkos,
	}
}

func BuildAnalysisResult(arithmetic uint64) *Result {
	junkoCount := arithmetic % 5
	arithmetic = arithmetic / 5
	pair := arithmetic % 15
	arithmetic = arithmetic / 15
	groups := uint64(0)
	x := utils.ToBytes(groups)
	i := 0
	for arithmetic != 0 {
		group := arithmetic % 15
		arithmetic = arithmetic / 15
		if group > 255 {
			panic(group)
		}
		x[i] = byte(groups)
		i++
	}
	return &Result{
		JunkoCount: junkoCount,
		Pair:       pair,
		Groups:     utils.ToUInt64(x),
	}
}

func GetUInt64Bytes4(value uint64) uint64 {
	x := utils.ToBytes(value)
	for i := uint64(3); i >= 0; i-- {
		if x[i] != 0 {
			return i + 1
		}
	}
	return 0
}

func GetUInt64Bytes8(value uint64) uint64 {
	x := utils.ToBytes(value)
	for i := uint64(7); i >= 0; i-- {
		if x[i] != 0 {
			return i + 1
		}
	}
	return 0
}

func Sort(value, length uint64) uint64 {
	x := utils.ToBytes(value)
	for i := uint64(0); i < length; i++ {
		for j := i + 1; j < length; j++ {
			if x[i] > x[j] {
				x[i], x[j] = x[j], x[i]
			}
		}
	}
	return utils.ToUInt64(x)
}

func (analysis *Result) GetArithmetic(end uint64) uint64 {
	result := uint64(0)
	groups := analysis.Groups
	x := utils.ToBytes(groups)
	for i := 3; i >= 0; i-- {

		if x[i] == 0 {
			continue
		}
		result = result*15 + uint64(x[i])
	}
	result = result*15 + analysis.Pair
	result = result*5 + analysis.JunkoCount
	result <<= 3
	bytes := GetUInt64Bytes4(result)
	return result | ((bytes - 1) << 1) | end
}

func (analysis *Result) String() string {
	groups := analysis.Groups
	junkoCnt := make([]uint64, 0)
	for junkoIndex := uint64(0); junkoIndex < analysis.JunkoCount; junkoIndex, groups = junkoIndex+1, groups>>8 {
		junkoCnt = append(junkoCnt, groups&0xFF)
	}
	pungCnt := make([]uint64, 0)
	if groups != 0 {
		for ; groups != 0; groups >>= 8 {
			pungCnt = append(pungCnt, groups&0xFF)
		}
	}
	return fmt.Sprintf("Pair %d, Junko %+v, Pung %+v", analysis.Pair, junkoCnt, pungCnt)
}
