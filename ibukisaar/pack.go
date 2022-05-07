package ibukisaar

import (
	"fmt"
	"mahjong/ibukisaar/analysis"
)

const (
	Pair int = iota
	Junko
	Pung
)

type Pack struct {
	Type   int
	Header int
	Tiles  []int
}

func (pack *Pack) String() string {
	return fmt.Sprintf("%v", pack.Tiles)
}

func dfs(info *analysis.Info, keys keys) [][]*Pack {
	keys.Reverse()
	temp := make([][]int, 0)
	for _, key := range keys {
		temp = append(temp, key.Tiles...)
	}
	groupsList := make([][]*Pack, 0)
	for i := 0; i < len(info.Results); i++ {
		indexes := make([]int, len(temp))
		var result = info.Results[i]
		groups := make([]*Pack, 0, 5)

		pairIndex := result.Pair - 1
		cur := indexes[pairIndex]
		indexes[pairIndex]++
		pairs := make([]int, 0, 2)
		pairs = append(pairs, temp[pairIndex][cur])
		cur = indexes[pairIndex]
		indexes[pairIndex]++
		pairs = append(pairs, temp[pairIndex][cur])
		groups = append(groups, &Pack{Type: Pair, Header: pairs[0], Tiles: pairs})
		groupIds := result.Groups

		for junkoIndex := uint(0); junkoIndex < uint(result.JunkoCount); junkoIndex, groupIds = junkoIndex+1, groupIds>>8 {
			index := int(groupIds & 0xFF)
			junko := make([]int, 0, 3)
			cur2 := indexes[index-2]
			indexes[index-2]++
			cur1 := indexes[index-1]
			indexes[index-1]++
			cur0 := indexes[index-0]
			indexes[index-0]++
			junko = append(junko, temp[index-2][cur2])
			junko = append(junko, temp[index-1][cur1])
			junko = append(junko, temp[index-0][cur0])
			groups = append(groups, &Pack{Type: Junko, Header: junko[0], Tiles: junko})
		}

		for ; groupIds != 0; groupIds >>= 8 {
			index := (int)(groupIds&0xFF) - 1
			cur2 := indexes[index]
			indexes[index]++
			cur1 := indexes[index]
			indexes[index]++
			cur0 := indexes[index]
			indexes[index]++
			pung := make([]int, 0, 3)
			pung = append(pung, temp[index][cur2])
			pung = append(pung, temp[index][cur1])
			pung = append(pung, temp[index][cur0])
			groups = append(groups, &Pack{Type: Pung, Header: pung[0], Tiles: pung})
		}
		groupsList = append(groupsList, groups)
	}
	return groupsList
}
