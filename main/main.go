package main

import (
	"fmt"
	"log"
	_ "mahjong/ibukisaar"
	"runtime"
	"time"
)

var tilesList = [][]int{
	{11, 11, 11, 21, 21, 21, 22, 22, 22, 23, 23, 23, 24, 24},
	{13, 13, 14, 15, 16, 24, 25, 26, 26, 27, 28, 29, 29, 13},
	{15, 15, 15, 16, 16, 16, 17, 17, 17, 25, 25, 25, 26, 27},
	{11, 11, 11, 12, 12, 12, 13, 13, 13, 14, 14, 14, 16, 16},
	{11, 11, 12, 12, 13, 13, 21, 21, 21, 31, 31, 31, 33, 32},
	{15, 15, 15, 16, 16, 16, 17, 17, 17, 25, 25, 25, 26, 27},
	{11, 11, 12, 12, 13, 13, 21, 21, 21, 23, 31, 31, 31, 22},
	{11, 12, 13, 13, 14, 15, 15, 16, 17, 22, 21, 21, 21, 22},
	{11, 11, 11, 11, 12, 13, 14, 15, 16, 17, 18, 19, 19, 19},
	{11, 11, 12, 12, 13, 13, 21, 21, 21, 31, 31, 31, 33, 32},
	{11, 12, 13, 23, 24, 25, 37, 38, 47, 47, 39},
	{11, 12, 13, 14, 15, 16, 39, 39},
	{11, 12, 13, 23, 24, 25, 37, 38, 47, 47, 39},
	{17, 18, 19, 27, 27, 27, 29, 37, 37, 38, 38, 39, 39, 29},
	{18, 18, 18, 28, 28, 28, 38, 38, 38, 41, 41, 41, 42, 42},
	{11, 12, 13, 21, 22, 23, 24, 25, 26, 29, 37, 38, 39, 29},
	{14, 14, 15, 15, 16, 16, 24, 24, 24, 26, 26, 35, 35, 35},
	{15, 15, 15, 16, 16, 16, 17, 17, 17, 25, 25, 25, 26, 27},
	{21, 22, 22, 23, 23, 23, 24, 24, 24, 25, 25, 26, 29, 29},
	{15, 17, 17, 17, 17, 18, 18, 18, 18, 19, 19, 19, 19, 15},
	{11, 11, 11, 12, 12, 12, 13, 13, 13, 14, 14, 14, 16, 16},
	{11, 11, 12, 12, 13, 13, 15, 17, 17, 18, 18, 19, 19, 15},
	{11, 11, 11, 21, 21, 21, 29, 29, 29, 39, 39, 39, 19, 19},
	{11, 14, 17, 22, 25, 28, 33, 36, 39, 41, 42, 43, 44, 45},
	{31, 31, 41, 41, 41, 42, 42, 42, 43, 43, 43, 44, 44, 31},
	{12, 12, 12, 13, 13, 45, 45, 45, 46, 46, 46, 47, 47, 13},
	{12, 13, 14, 22, 22, 22, 28, 28, 29, 29, 29, 41, 41, 41},
	{13, 14, 15, 16, 12, 17, 35, 35, 35, 36, 36, 36, 42, 42},
	{12, 13, 14, 15, 16, 17, 24, 24, 24, 25, 25, 26, 27, 28},
	{12, 13, 13, 13, 14, 15, 16, 17, 25, 26, 26, 27, 27, 28},
}

func main() {
	//for _, ints := range tilesList {
	//	sort.Ints(ints)
	//}

	now := time.Now()
	//for i := 0; i < 100000; i++ {
	//	for _, tiles := range tilesList {
	//		keys := ibukisaar.Parse(tiles)
	//		key := ibukisaar.BuildKey(keys)
	//		info, _ := ibukisaar.ShantenMap.Load(key)
	//		if info == nil {
	//			log.Printf("%v", tiles)
	//			os.Exit(0)
	//		}
	//		list := ibukisaar.Analysis(info.(*analysis.Info), keys)
	//		_ = list
	//	}
	//}
	log.Printf("time use %dms", time.Now().Sub(now).Milliseconds())
	runtime.GC()
	fmt.Scanln()
}
