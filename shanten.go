package mahjong

import "sort"

const (
	Error   = -1
	TaTsu   = iota // 搭子
	Toitsu         // 对子
	Wanzu          // 万子
	Pinzu          // 筒子
	Shuntsu        // 顺子
	KoTsu          // 刻子
	Kantsu         // 杠子
	Mentsu         // 面子
	Kanta          // 嵌塔
	Penta          // 辺塔
)

// Shanten 求向听数的函数
func Shanten(tiles []int) int {
	sort.Ints(tiles)
	return 0
}

// dfs 深度遍历
// pack 当表示搭子时 [第一位, 是否为顺子, 是否为跳搭]
// pack 当表示面子时 [第一位, 是否为顺子,]
func dfs(tiles []int, vis []byte, packs [][3]int) {
	tilesLen := len(tiles)
	vis = make([]byte, tilesLen)
	start := 0

	for index, v := range vis {
		if v == 0 {
			start = index
			break
		}
	}

	for mid := start + 1; mid < tilesLen; mid++ {
		if JudgeTaTsu(tiles[start], tiles[mid]) {
			vis[start] = 1
			mid = 1
			koTsu, kanta := 1, 1
			if tiles[start] == tiles[mid] {
				koTsu = 0
			}
			if tiles[start] == tiles[mid]-1 {
				kanta = 0
			}
			packs = append(packs, [3]int{tiles[start], koTsu, kanta})
		}

		if Judge2SameOrAdjacent(tiles[start], tiles[mid]) {
			for end := mid + 1; end < tilesLen; end++ {
				mentsu := Judge3MakePack(tiles[start], tiles[mid], tiles[end])
				if mentsu == Error {
					continue
				}
			}
		}
	}

	// if !CheckVis(vis) {
	// 	vis[start] = 1
	// }
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

func Judge2MakePack() {

}

// Judge3MakePack 判断a < b < c时是否为顺子或者刻子
func Judge3MakePack(a, b, c int) int {
	if a == b && a == c {
		return KoTsu
	} else if a == b-1 && a == b-2 {
		return Shuntsu
	}
	return Error
}

// JudgeTaTsu 判断a < b时是否为搭子
func JudgeTaTsu(a, b int) bool {
	return a == b || a == b-1 || a == b-2
}

// Judge2SameOrAdjacent 判断a < b时是否相邻或者相同
func Judge2SameOrAdjacent(a, b int) bool {
	return a == b || a == b-1
}
