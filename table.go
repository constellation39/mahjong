package mahjong

import (
	"fmt"
	"mahjong/analysis"
	"mahjong/ron"
	. "mahjong/utils"
)

// 原理：https://zhuanlan.zhihu.com/p/31000381
//

func Shanten(value uint64) (int, []*analysis.Result) {
	var shift, tiles, remCount uint64
	for (value>>shift)&0b1111 != 0b1111 {
		continuous := value >> (shift + 2) & 0b11
		cnt := (value >> shift & 0b11) + 1
		remCount += cnt
		tiles |= (continuous*5 + cnt) << shift
		shift += 4
	}

	if ron.Ron(tiles, remCount) {
		ret := Analysis(tiles)
		return -1, ret
	}

	if value == 16252928 {
		fmt.Println()
	}

	args := &analysis.SyantenArgs{
		N:               remCount / 3,
		Result:          13,
		MaxUseTileCount: 0,
	}
	Syanten(tiles, remCount, args)
	return int(args.Result), nil
}

func Analysis(value uint64) (ret []*analysis.Result) {
	ret = make([]*analysis.Result, 0)
	AnalysisCutPair(value, &ret)
	return
}

func AnalysisCutPair(value uint64, ret *[]*analysis.Result) {
	for shift := uint64(0); (value >> shift) != 0; shift += 4 {
		continuous, singleCount := Get(value, shift)
		if singleCount >= 2 {
			AnalysisCut3(Set(value, shift, continuous, singleCount-2), 0, (shift>>2)+1, 0, 0, 0, 0, ret)
		}
	}
}

func AnalysisCut3(value, shift, pair, junkoCount, junkos, pungCount, pungs uint64, ret *[]*analysis.Result) {
	for (value>>shift) != 0 && ((value>>shift)&0xF)%5 == 0 {
		shift += 4
	}
	if value>>shift == 0 {
		newResult := analysis.New(pair, junkoCount, junkos, pungCount, pungs)
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

func Syanten(value, count uint64, args *analysis.SyantenArgs) {
	for shift := uint64(0); (value >> shift) != 0; shift += 4 {
		continuous, singleCount := Get(value, shift)
		if singleCount >= 2 {
			SyantenCut3(Set(value, shift, continuous, singleCount-2), 0, count-2, 0, 0, 1, args)
		}
	}
	SyantenCut3(value, 0, count, 0, 0, 0, args)
}

func SyantenCut3(value, shift, remCount, c3, c2, p uint64, args *analysis.SyantenArgs) {
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

func SyantenCut2(value, shift, remCount, c3, c2, p uint64, args *analysis.SyantenArgs) {
	if args.Result == 0 {
		return
	}
	if c3+c2 > args.N {
		return
	}
	useTileCount := c3 + (c3+c2+p)*2
	if remCount < args.MaxUseTileCount-useTileCount {
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
