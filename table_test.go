package mahjong

import (
	"log"
	"testing"
	"time"
)

func TestDistanceEnum(t *testing.T) {

}

func TestQuantityEnum(t *testing.T) {
	for i := 2; i <= 14; i += 3 {
		ret := QuantityEnum(i)
		log.Printf("%+v", ret)
		break
	}
}

func Test_distanceEnum(t *testing.T) {

}

func Test_quantityEnum(t *testing.T) {

}

func PrintTiles(tiles uint64) {
	log.Printf("%b", tiles)
}

func TestTilesEnum(t *testing.T) {
	ret := make(map[uint64]struct{})
	cnt := 0
	now := time.Now()
	for i := 2; i <= 14; i += 3 {
		for _, tiles := range TilesEnum(i) {
			if _, ok := ret[tiles]; !ok {
				cnt++
				ret[tiles] = struct{}{}
			}
		}
		log.Printf("i %d cnt %d time %s", i, cnt, time.Now().Sub(now))
	}

	for tiles := range ret {
		GetInfo(tiles)
	}
}
