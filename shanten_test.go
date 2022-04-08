package mahjong

import (
	"log"
	"sort"
	"testing"
	"time"
)

type ShantenNode struct {
	Level  int                  // 向听数
	Blocks [][]int              // 拆解后的手牌
	Nodes  map[int]*ShantenNode // 打出某张牌后的向听树
}

func TestCheckVis(t *testing.T) {

}

func TestJudge3MakePack(t *testing.T) {

}

func TestJudgeTaTsu(t *testing.T) {

}

func TestPackHashCode(t *testing.T) {

}

func TestShanten(t *testing.T) {
	tilesList := [][]int{
		{12, 23, 14, 22, 22, 22, 28, 28, 29, 29, 29, 41, 11, 41},
		{13, 24, 15, 16, 12, 17, 35, 35, 35, 36, 36, 36, 12, 42},
		{12, 23, 14, 15, 16, 17, 24, 24, 24, 25, 25, 26, 17, 28},
		{12, 23, 13, 13, 14, 15, 16, 17, 25, 26, 26, 27, 17, 28},
		{12, 23, 13, 14, 15, 16, 17, 25, 26, 26, 27, 27, 18, 13},
		{11, 21, 11, 12, 13, 14, 15, 16, 17, 18, 19, 19, 29, 19},
		{11, 21, 12, 12, 13, 13, 14, 14, 15, 15, 16, 16, 27, 17},
		{22, 32, 22, 22, 24, 24, 24, 24, 26, 28, 28, 28, 18, 26},
		{12, 22, 12, 13, 13, 45, 45, 45, 46, 46, 46, 47, 27, 13},
		{31, 21, 41, 41, 41, 42, 42, 42, 43, 43, 43, 44, 14, 31},
		{11, 24, 17, 22, 25, 28, 33, 36, 39, 41, 42, 43, 14, 45},
		{11, 21, 11, 21, 21, 21, 29, 29, 29, 39, 39, 39, 39, 19},
		{11, 21, 12, 12, 13, 13, 15, 17, 17, 18, 18, 19, 39, 15},
		{11, 21, 11, 12, 12, 12, 13, 13, 13, 14, 14, 14, 36, 16},
		{15, 17, 17, 17, 17, 18, 18, 18, 18, 19, 19, 19, 29, 15},
		{21, 12, 22, 23, 23, 23, 24, 24, 24, 25, 25, 26, 19, 29},
		{15, 25, 15, 16, 16, 16, 17, 17, 17, 25, 25, 25, 16, 27},
		{14, 24, 15, 15, 16, 16, 24, 24, 24, 26, 26, 35, 25, 35},
		{11, 22, 13, 21, 22, 23, 24, 25, 26, 29, 37, 38, 29, 29},
		{18, 28, 18, 28, 28, 28, 38, 38, 38, 41, 41, 41, 32, 42},
		{17, 28, 19, 27, 27, 27, 29, 37, 37, 38, 38, 39, 19, 29},
		{11, 21, 12, 12, 13, 13, 21, 21, 21, 31, 31, 31, 13, 32},
		{11, 22, 13, 13, 14, 15, 15, 16, 17, 22, 21, 21, 31, 22},
		{11, 21, 12, 12, 13, 13, 21, 21, 21, 23, 31, 31, 21, 22},
		{15, 25, 15, 16, 16, 16, 17, 17, 17, 25, 25, 25, 46, 27},
		{11, 21, 12, 12, 13, 13, 21, 21, 21, 31, 31, 31, 13, 32},
		{11, 21, 11, 12, 12, 12, 13, 13, 13, 14, 14, 14, 26, 16},
		{15, 25, 15, 16, 16, 16, 17, 17, 17, 25, 25, 25, 36, 27},
		{13, 23, 14, 15, 16, 24, 25, 26, 26, 27, 28, 29, 19, 13},
	}

	for _, tiles := range tilesList {
		sort.Sort(sort.IntSlice(tiles))
	}

	for index, tiles := range tilesList {
		now := time.Now()
		shantenNode := ShantenLevel(tiles, -1)
		log.Printf("%d %+v use time %s shantenNode %+v", index, tiles, time.Now().Sub(now), shantenNode)
	}
}

func ShantenLevel(tiles []int, target int) (shantenNode *ShantenNode) {
	shanten, blocks := Shanten(tiles)
	shantenNode = &ShantenNode{
		Level:  shanten,
		Blocks: blocks,
		Nodes:  make(map[int]*ShantenNode),
	}

	if shanten > target {
		shantenLevel(tiles, target, shanten, shantenNode)
	}
	return
}

func shantenLevel(tiles []int, target int, shanten int, shantenNode *ShantenNode) {
	temp := make([]int, len(tiles))
	copy(temp, tiles)

	for index := range tiles {
		for suit := 1; suit <= 4; suit++ {
			for i := 1; i <= 9; i++ {
				copy(tiles, temp)
				add := suit*10 + i
				tiles[index] = add
				sort.Ints(tiles)
				s, b := Shanten(tiles)
				if shanten <= s {
					continue
				}
				shantenNode.Nodes[temp[index]] = &ShantenNode{
					Level:  s,
					Blocks: b,
					Nodes:  make(map[int]*ShantenNode),
				}
				//if s > target {
				//	shantenLevel(tiles, target, s, shantenNode.Nodes[temp[index]])
				//}
			}
		}
	}
}
