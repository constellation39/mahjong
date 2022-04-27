package analysis

import (
	"fmt"
	. "mahjong/ibukisaar/utils"
	"reflect"
)

type SyantenArgs struct {
	N               uint64
	Result          uint64
	MaxUseTileCount uint64
}

type Result struct {
	Pair       byte
	JunkoCount byte
	Groups     uint
}

type Info struct {
	Shanten int8
	Results []*Result
}

func New(pair, junkoCount, junkos, pungCount, pungs uint64) *Result {
	junkos = Sort(junkos, junkoCount)
	pungs = Sort(pungs, pungCount)
	return &Result{
		Pair:       byte(pair),
		JunkoCount: byte(junkoCount),
		Groups:     uint(pungs<<(junkoCount*8) | junkos),
	}
}

func BuildAnalysisResult(arithmetic uint) *Result {
	junkoCount := arithmetic % 5
	arithmetic = arithmetic / 5

	pair := arithmetic % 15
	arithmetic = arithmetic / 15
	groups := uint64(0)
	x := ToBytes(groups)
	i := 0
	for arithmetic != 0 {
		group := arithmetic % 15
		arithmetic = arithmetic / 15
		x[i] = byte(group)
		i++
	}
	return &Result{
		JunkoCount: byte(junkoCount),
		Pair:       byte(pair),
		Groups:     uint(ToUInt64(x)),
	}
}

func (analysis *Result) GetArithmetic(end uint64) uint {
	result := uint64(0)
	groups := analysis.Groups
	x := ToBytes(uint64(groups))
	for i := 3; i >= 0; i-- {
		if x[i] == 0 {
			continue
		}
		result = result*15 + uint64(x[i])
	}
	result = result*15 + uint64(analysis.Pair)
	result = result*5 + uint64(analysis.JunkoCount)
	result <<= 3
	bytes := GetUInt64Bytes8(result)
	return uint(result | ((bytes - 1) << 1) | end)
}

func (analysis *Result) String() string {
	groups := analysis.Groups
	junkoCnt := make([]uint, 0)
	for junkoIndex := byte(0); junkoIndex < analysis.JunkoCount; junkoIndex, groups = junkoIndex+1, groups>>8 {
		junkoCnt = append(junkoCnt, groups&0xFF)
	}
	pungCnt := make([]uint, 0)
	if groups != 0 {
		for ; groups != 0; groups >>= 8 {
			pungCnt = append(pungCnt, groups&0xFF)
		}
	}
	return fmt.Sprintf("Pair %d, Junko %+v, Pung %+v", analysis.Pair, junkoCnt, pungCnt)
}

func GetUInt64Bytes4(value uint64) uint64 {
	x := ToBytes(value)
	for i := uint64(3); i <= 0; i-- {
		if x[i] != 0 {
			return i + 1
		}
	}
	return 0
}

func GetUInt64Bytes8(value uint64) uint64 {
	x := ToBytes(value)
	for i := uint64(7); i >= 0; i-- {
		if x[i] != 0 {
			return i + 1
		}
	}
	return 0
}

func Sort(value, length uint64) uint64 {
	x := ToBytes(value)
	for i := uint64(0); i < length; i++ {
		for j := i + 1; j < length; j++ {
			if x[i] > x[j] {
				x[i], x[j] = x[j], x[i]
			}
		}
	}
	return ToUInt64(x)
}

func Ron(value, cnt uint64) bool {
	for shift := uint64(0); (value >> shift) != 0; shift += 4 {
		continuous, singleCount := Get(value, shift)
		if singleCount < 2 {
			continue
		}
		if CutPung(Set(value, shift, continuous, singleCount-2), 0, cnt-2) {
			return true
		}
	}
	return false
}

func CutPung(value, shift, cnt uint64) bool {
	if cnt == 0 {
		return true
	}
	for (value>>shift) != 0 && ((value>>shift)&0xF)%5 == 0 {
		shift += 4
	}
	if (value >> shift) == 0 {
		return CutJunko(value, 0, cnt)
	}
	continuous, singleCount := Get(value, shift)
	if singleCount >= 3 {
		if CutPung(Set(value, shift, continuous, singleCount-3), shift, cnt-3) {
			return true
		}
	}
	return CutPung(value, shift+4, cnt)
}

func CutJunko(value, shift, cnt uint64) bool {
	if cnt == 0 {
		return true
	}

	for (value>>shift) != 0 && ((value>>shift)&0xF)%5 == 0 {
		shift += 4
	}

	continuous1, singleCount1 := Get(value, shift)
	if continuous1 == 0 {
		continuous2, singleCount2 := Get(value, shift+4)
		if continuous2 == 0 && singleCount2 > 0 {
			continuous3, singleCount3 := Get(value, shift+8)
			if singleCount3 > 0 {
				var valueT = Set(value, shift, continuous1, singleCount1-1)
				valueT = Set(valueT, shift+4, continuous2, singleCount2-1)
				valueT = Set(valueT, shift+8, continuous3, singleCount3-1)
				return CutJunko(valueT, shift, cnt-3)
			}
		}
	}
	return false
}

func Shanten(value uint64) *Info {
	var shift, tiles, remCount uint64
	for (value>>shift)&0b1111 != 0b1111 {
		continuous := value >> (shift + 2) & 0b11
		cnt := (value >> shift & 0b11) + 1
		remCount += cnt
		tiles |= (continuous*5 + cnt) << shift
		shift += 4
	}

	if Ron(tiles, remCount) {
		return &Info{
			Shanten: -1,
			Results: Analysis(tiles),
		}
	}

	args := &SyantenArgs{
		N:               remCount / 3,
		Result:          13,
		MaxUseTileCount: 0,
	}
	Syanten(tiles, remCount, args)
	return &Info{
		Shanten: int8(args.Result),
		Results: nil,
	}
}

func Analysis(value uint64) (ret []*Result) {
	ret = make([]*Result, 0)
	AnalysisCutPair(value, &ret)
	return
}

func AnalysisCutPair(value uint64, ret *[]*Result) {
	for shift := uint64(0); (value >> shift) != 0; shift += 4 {
		continuous, singleCount := Get(value, shift)
		if singleCount >= 2 {
			AnalysisCut3(Set(value, shift, continuous, singleCount-2), 0, (shift>>2)+1, 0, 0, 0, 0, ret)
		}
	}
}

func AnalysisCut3(value, shift, pair, junkoCount, junkos, pungCount, pungs uint64, ret *[]*Result) {
	for (value>>shift) != 0 && ((value>>shift)&0xF)%5 == 0 {
		shift += 4
	}
	if value>>shift == 0 {
		newResult := New(pair, junkoCount, junkos, pungCount, pungs)
		for _, result := range *ret {
			if reflect.DeepEqual(result, newResult) {
				return
			}
		}
		*ret = append(*ret, newResult)
		return
	}

	continuous, singleCount := Get(value, shift)
	if singleCount >= 3 {
		var temp = Set(value, shift, continuous, singleCount-3)
		AnalysisCut3(temp, shift, pair, junkoCount, junkos, pungCount+1, (pungs<<8)|((shift>>2)+1), ret)
	}
	if continuous == 0 {
		continuous2, singleCount2 := Get(value, shift+4)
		if continuous2 == 0 && singleCount2 > 0 {
			continuous3, singleCount3 := Get(value, shift+8)
			if singleCount3 > 0 {
				var temp = Set(value, shift, continuous, singleCount-1)
				temp = Set(temp, shift+4, continuous2, singleCount2-1)
				temp = Set(temp, shift+8, continuous3, singleCount3-1)
				AnalysisCut3(temp, shift, pair, junkoCount+1, (junkos<<8)|((shift>>2)+2), pungCount, pungs, ret)
			}
		}
	}
}

func Syanten(value, count uint64, args *SyantenArgs) {
	for shift := uint64(0); (value >> shift) != 0; shift += 4 {
		continuous, singleCount := Get(value, shift)
		if singleCount >= 2 {
			SyantenCut3(Set(value, shift, continuous, singleCount-2), 0, count-2, 0, 0, 1, args)
		}
	}
	SyantenCut3(value, 0, count, 0, 0, 0, args)
}

func SyantenCut3(value, shift, remCount, c3, c2, p uint64, args *SyantenArgs) {
	for value>>shift != 0 && ((value>>shift)&0xF)%5 == 0 {
		shift += 4
	}
	if value>>shift == 0 {
		SyantenCut2(value, 0, remCount, c3, c2, p, args)
		return
	}
	continuous, singleCount := Get(value, shift)
	if singleCount >= 3 {
		SyantenCut3(Set(value, shift, continuous, singleCount-3), shift+4, remCount-3, c3+1, c2, p, args)
	}
	if continuous == 0 {
		continuous2, singleCount2 := Get(value, shift+4)
		if continuous2 == 0 && singleCount2 > 0 {
			continuous3, singleCount3 := Get(value, shift+8)
			if singleCount3 > 0 {
				var temp = Set(value, shift, continuous, singleCount-1)
				temp = Set(temp, shift+4, continuous2, singleCount2-1)
				temp = Set(temp, shift+8, continuous3, singleCount3-1)
				SyantenCut3(temp, shift, remCount-3, c3+1, c2, p, args)
			}
		}
	}
	SyantenCut3(value, shift+4, remCount, c3, c2, p, args)
}

func SyantenCut2(value, shift, remCount, c3, c2, p uint64, args *SyantenArgs) {
	if args.Result == 0 {
		return
	}
	if c3+c2 > args.N {
		return
	}
	useTileCount := c3 + (c3+c2+p)*2
	t := int(args.MaxUseTileCount) - int(useTileCount)
	if int(remCount) < t {
		return
	}

	if remCount == 0 {
		num := (args.N-c3)*2 - c2 - p
		if num < args.Result {
			args.Result = num
		}
		if args.MaxUseTileCount < useTileCount {
			args.MaxUseTileCount = useTileCount
		}
		return
	}

	for ((value>>shift)&0xF)%5 == 0 {
		shift += 4
	}

	continuous, singleCount := Get(value, shift)
	if singleCount >= 2 {
		SyantenCut2(Set(value, shift, continuous, singleCount-2), shift, remCount-2, c3, c2+1, p, args)
	}
	if continuous == 0 {
		continuous2, singleCount2 := Get(value, shift+4)
		if singleCount2 > 0 {
			var temp = Set(value, shift, continuous, singleCount-1)
			temp = Set(temp, shift+4, continuous2, singleCount2-1)
			SyantenCut2(temp, shift, remCount-2, c3, c2+1, p, args)
		}
		if continuous2 == 0 {
			continuous3, singleCount3 := Get(value, shift+8)
			if singleCount3 > 0 {
				var temp = Set(value, shift, continuous, singleCount-1)
				temp = Set(temp, shift+8, continuous3, singleCount3-1)
				SyantenCut2(temp, shift, remCount-2, c3, c2+1, p, args)
			}
		}
	} else if continuous == 1 {
		continuous3, singleCount3 := Get(value, shift+4)
		if singleCount3 > 0 {
			var temp = Set(value, shift, continuous, singleCount-1)
			temp = Set(temp, shift+4, continuous3, singleCount3-1)
			SyantenCut2(temp, shift, remCount-2, c3, c2+1, p, args)
		}
	}

	SyantenCut2(value, shift+4, remCount-singleCount, c3, c2, p, args)
}
