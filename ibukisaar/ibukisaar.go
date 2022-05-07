package ibukisaar

import "mahjong/ibukisaar/analysis"

func Query(tiles []int) (*analysis.Info, bool) {
	ks := parse(tiles)
	k := buildKey(ks)
	info, ok := ShantenMap.Load(k)
	if !ok {
		return nil, false
	}
	return info.(*analysis.Info), true
}

func Analysis(tiles []int) [][][]int {
	ks := parse(tiles)
	k := buildKey(ks)
	info, ok := ShantenMap.Load(k)
	if !ok {
		return nil
	}
	return dfs(info.(*analysis.Info), ks)
}
