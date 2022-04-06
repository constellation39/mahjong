package mahjong

import (
	"log"
	"sort"
	"testing"
	"time"
)

func TestCheckVis(t *testing.T) {

}

func TestJudge3MakePack(t *testing.T) {

}

func TestJudgeTaTsu(t *testing.T) {

}

func TestPackHashCode(t *testing.T) {

}

func TestShanten(t *testing.T) {

}

func Test_dfs(t *testing.T) {
	tilesList := [][]int{
		{12, 13, 14, 22, 22, 22, 28, 28, 28, 29, 29, 29, 41, 41},
		{12, 13, 14, 15, 16, 17, 24, 24, 24, 25, 25, 26, 27},
		{12, 13, 14, 15, 16, 17, 35, 35, 35, 36, 36, 36, 42},
		{12, 13, 13, 13, 14, 15, 16, 17, 25, 26, 26, 27, 27},
		{12, 13, 13, 14, 15, 16, 17, 25, 26, 26, 27, 27, 28},
		{11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 19, 19, 19},
		{11, 11, 12, 12, 13, 13, 14, 14, 15, 15, 16, 16, 17},
		{22, 22, 22, 22, 24, 24, 24, 24, 26, 28, 28, 28, 23},
		{13, 13, 14, 15, 16, 24, 25, 26, 26, 27, 28, 29, 29},
		{15, 15, 15, 16, 16, 16, 17, 17, 17, 25, 25, 25, 26},
		{11, 11, 11, 12, 12, 12, 13, 13, 13, 14, 14, 14, 16},
		{11, 11, 12, 12, 13, 13, 21, 21, 21, 31, 31, 31, 33},
		{15, 15, 15, 16, 16, 16, 17, 17, 17, 25, 25, 25, 26},
		{11, 11, 12, 12, 13, 13, 21, 21, 21, 23, 31, 31, 31},
		{11, 12, 13, 13, 14, 15, 15, 16, 17, 22, 21, 21, 21},
		{12, 17, 26, 27, 32, 34, 35, 39, 42, 43, 43, 46, 47, 47},
	}
	//tiles := []int{12, 13, 14, 22, 22, 22, 28, 28, 29, 29, 29, 41, 41, 41}

	for _, tiles := range tilesList {
		sort.Sort(sort.IntSlice(tiles))
	}

	now := time.Now()

	for _, tiles := range tilesList {
		//sort.Sort(sort.IntSlice(tiles))
		shanten, blocks := Shanten(tiles)
		_ = shanten
		_ = blocks
		log.Printf("%+v shanten %d blocks %+v", tiles, shanten, blocks)
	}

	log.Printf("use time %s", time.Now().Sub(now))
}
