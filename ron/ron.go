package ron

import . "mahjong/utils"

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
