package mahjong

import (
	"sort"
)

const (
	Error   = -1
	TaTsu   = 0    // 搭子
	Toitsu  = iota // 对子 (对子是搭子)
	Kanta          // 嵌塔
	Penta          // 辺塔
	Wanzu          // 万子
	Pinzu          // 筒子
	Shuntsu        // 顺子
	KoTsu          // 刻子
	Kantsu         // 杠子
	Mentsu         // 面子
)

var TextMap = map[int]string{
	Error:   "错误",
	TaTsu:   "搭子",
	Toitsu:  "对子",
	Kanta:   "嵌塔",
	Penta:   "辺塔",
	Wanzu:   "万子",
	Pinzu:   "筒子",
	Shuntsu: "顺子",
	KoTsu:   "刻子",
	Kantsu:  "杠子",
	Mentsu:  "面子",
}

type Pack struct {
	Type  int // 牌组类型
	Tile  int // 用一张牌表示牌组，该牌为牌组的第一张
	Offer int
	//供牌信息，0为暗手，1、2、3为副露
	//若为吃，则分别表示吃的牌是_tile的左、中、右
	//若为碰或杠，则分别表示碰的牌来自于上家、对家、下家
	//若为加杠，则用5、6、7来表示（便于mod 4） 3-补杠 4-点杠 5-暗杠
	//若包含最后所和的牌，使用-1表示自摸牌，-2表示铳和牌。实际上是Tile的_drawflag取负
}

// Shanten 求向听数的函数
func Shanten(tiles []int) (int, [][]int) {
	sort.Ints(tiles)
	return Dfs(tiles)
}

func Dfs(tiles []int) (int, [][]int) {
	shanten := 99
	blocks := new([][]int)
	dfs(tiles, make([]byte, len(tiles)), []*Pack{}, map[int]struct{}{}, [][]int{}, &shanten, blocks)
	return shanten, *blocks
}

// dfs 深度遍历
func dfs(tiles []int, vis []byte, packs []*Pack, st map[int]struct{}, stack [][]int, shanten *int, blocks *[][]int) (ret int) {
	tilesLen := len(tiles)
	start := -1
	if CheckVis(vis) {
		var x, m, d, c, q int
		for _, pack := range packs {
			switch pack.Type {
			case Shuntsu, KoTsu:
				m++
			case Toitsu:
				q = 1
				fallthrough
			case Kanta, Penta:
				d++
			}
		}

		if m+d > 5 {
			c = m + d - 5
		}

		if m+d <= 4 {
			q = 1
		}

		x = 9 - 2*m - d + c - q

		if *shanten > x {
			*shanten = x
			*blocks = make([][]int, len(stack))
			for index, ints := range stack {
				(*blocks)[index] = ints
			}
		}
		return 1
	}

	for index, v := range vis {
		if v == 0 {
			start = index
			break
		}
	}

	if start == -1 {
		return
	}

	for mid := start + 1; mid < tilesLen; mid++ {
		if vis[mid] > 0 {
			continue
		}
		taTsu := JudgeTaTsu(tiles[start], tiles[mid])

		if taTsu != Error {
			packs = append(packs, &Pack{
				Type:  taTsu,
				Tile:  tiles[start],
				Offer: 0,
			})
			vis[start] = 1
			vis[mid] = 1
			hash := PackHashCode(packs)
			stack = append(stack, []int{tiles[start], tiles[mid]})
			if _, ok := st[hash]; !ok {
				st[hash] = struct{}{}
				ret |= dfs(tiles, vis, packs, st, stack, shanten, blocks)
			}
			vis[start] = 0
			vis[mid] = 0
			stack = stack[:len(stack)-1]
			packs = packs[:len(packs)-1]
		}

		if taTsu == Toitsu || taTsu == Penta {
			for end := mid + 1; end < tilesLen; end++ {
				if vis[end] > 0 {
					continue
				}
				menTsu := Judge3MakePack(tiles[start], tiles[mid], tiles[end])
				if menTsu == Error {
					continue
				}
				packs = append(packs, &Pack{
					Type:  menTsu,
					Tile:  tiles[start],
					Offer: 0,
				})
				stack = append(stack, []int{tiles[start], tiles[mid], tiles[end]})
				vis[start] = 1
				vis[mid] = 1
				vis[end] = 1
				hash := PackHashCode(packs)
				if _, ok := st[hash]; !ok {
					st[hash] = struct{}{}
					ret |= dfs(tiles, vis, packs, st, stack, shanten, blocks)
				}
				vis[start] = 0
				vis[mid] = 0
				vis[end] = 0
				stack = stack[:len(stack)-1]
				packs = packs[:len(packs)-1]
				if ret > 0 {
					return
				}
			}
		}
	}

	if !CheckVis(vis) {
		vis[start] = 1
		stack = append(stack, []int{tiles[start]})
		dfs(tiles, vis, packs, st, stack, shanten, blocks)
		stack = stack[:len(stack)-1]
		vis[start] = 0
	}

	return
}

// PackHashCode hashCode
func PackHashCode(packs []*Pack) (ret int) {
	for i := 0; i < len(packs); i++ {
		pack := packs[i]
		ret = ret<<8 | pack.Tile<<4 | pack.Type<<1
	}
	return
}

// CheckVis 检查vis是否全部使用
func CheckVis(vis []byte) bool {
	for _, v := range vis {
		if v == 0 {
			return false
		}
	}
	return true
}

// Judge3MakePack 判断a < b < c时是否为顺子或者刻子
func Judge3MakePack(a, b, c int) int {
	if a == b && a == c {
		return KoTsu
	} else if a == b-1 && a == c-2 {
		return Shuntsu
	}
	return Error
}

// JudgeTaTsu 判断a < b时是否为搭子
func JudgeTaTsu(a, b int) int {
	// a == b || a == b-1 || a == b-2
	if a == b {
		return Toitsu
	} else if a == b-1 {
		return Penta
	} else if a == b-2 {
		return Kanta
	}
	return Error
}
