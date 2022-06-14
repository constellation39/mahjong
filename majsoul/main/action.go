package main

import (
	"fmt"
	"majsoul"
	"majsoul/message"
	"strconv"
	"time"
	"utils/logger"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

func (m *Majsoul) ActionPrototype(notify *message.ActionPrototype) {
	data := message.GetActionType(notify.Name)
	err := proto.Unmarshal(notify.Data, data)
	if err != nil {
		logger.Error("ActionPrototype", zap.Error(err))
		return
	}
	logger.Debug("Majsoul.ActionPrototype", zap.Reflect("name", notify.Name), zap.Reflect("data", data))
	switch notify.Name {
	case "ActionMJStart":
		m.ActionMJStart(data.(*message.ActionMJStart))
	case "ActionNewCard":
		m.ActionNewCard(data.(*message.ActionNewCard))
	case "ActionNewRound":
		m.ActionNewRound(data.(*message.ActionNewRound))
	case "ActionSelectGap":
		m.ActionSelectGap(data.(*message.ActionSelectGap))
	case "ActionChangeTile":
		m.ActionChangeTile(data.(*message.ActionChangeTile))
	case "ActionRevealTile":
		m.ActionRevealTile(data.(*message.ActionRevealTile))
	case "ActionUnveilTile":
		m.ActionUnveilTile(data.(*message.ActionUnveilTile))
	case "ActionLockTile":
		m.ActionLockTile(data.(*message.ActionLockTile))
	case "ActionDiscardTile":
		m.ActionDiscardTile(data.(*message.ActionDiscardTile))
	case "ActionDealTile":
		m.ActionDealTile(data.(*message.ActionDealTile))
	case "ActionChiPengGang":
		m.ActionChiPengGang(data.(*message.ActionChiPengGang))
	case "ActionGangResult":
		m.ActionGangResult(data.(*message.ActionGangResult))
	case "ActionGangResultEnd":
		m.ActionGangResultEnd(data.(*message.ActionGangResultEnd))
	case "ActionAnGangAddGang":
		m.ActionAnGangAddGang(data.(*message.ActionAnGangAddGang))
	case "ActionBaBei":
		m.ActionBaBei(data.(*message.ActionBaBei))
	case "ActionHule":
		m.ActionHule(data.(*message.ActionHule))
	case "ActionHuleXueZhanMid":
		m.ActionHuleXueZhanMid(data.(*message.ActionHuleXueZhanMid))
	case "ActionHuleXueZhanEnd":
		m.ActionHuleXueZhanEnd(data.(*message.ActionHuleXueZhanEnd))
	case "ActionLiuJu":
		m.ActionLiuJu(data.(*message.ActionLiuJu))
	case "ActionNoTile":
		m.ActionNoTile(data.(*message.ActionNoTile))
	}
}

func (m *Majsoul) ActionMJStart(in *message.ActionMJStart) {
	names := make([]string, 4)
	for i, uid := range m.GameInfo.SeatList {
		if uid == m.Account.AccountId {
			m.Seat = i
		}
		if uid == 0 {
			names[i] = fmt.Sprintf("Bot%d", i+1)
			continue
		}
		names[i] = strconv.Itoa(int(uid))
	}

	m.UAkochan.StartGame(m.Seat, names)
}
func (m *Majsoul) ActionNewCard(in *message.ActionNewCard) {
	logger.Debug("ActionNewCard", zap.Reflect("in", in))
}
func (m *Majsoul) ActionNewRound(in *message.ActionNewRound) {
	bakaze := GetUAkochanBakaze(in.Chang)
	kyoku := in.Ju + 1
	honba := in.Ben
	kyotaku := in.Liqibang
	oya := in.Ju
	dora_marker := GetUAkochanTile(in.Doras[0])
	tehais := make([][]string, 4)
	for i := 0; i < 4; i++ {
		tehais[i] = []string{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"}
	}
	tehais[m.Seat] = GetUAkochanTiles(in.Tiles[:13])
	m.UAkochan.StartKyoku(bakaze, kyoku, honba, kyotaku, oya, dora_marker, tehais)
	if in.Ju != uint32(m.Seat) {
		return
	}
	time.Sleep(time.Second * 3)
	m.Tsumo(m.Seat, in.Tiles[13], in.Operation)
}
func (m *Majsoul) ActionSelectGap(in *message.ActionSelectGap) {
}
func (m *Majsoul) ActionChangeTile(in *message.ActionChangeTile) {
}
func (m *Majsoul) ActionRevealTile(in *message.ActionRevealTile) {
}
func (m *Majsoul) ActionUnveilTile(in *message.ActionUnveilTile) {
}
func (m *Majsoul) ActionLockTile(in *message.ActionLockTile) {
}
func (m *Majsoul) ActionDiscardTile(in *message.ActionDiscardTile) {
	m.Dahai(int(in.Seat), in.Tile, in.Moqie, in.Operation)
}
func (m *Majsoul) ActionDealTile(in *message.ActionDealTile) {
	if in.Tile == "" {
		in.Tile = "?"
	}
	m.Tsumo(int(in.Seat), in.Tile, in.Operation)
	if in.Liqi != nil {
		m.ReachAccepted(int(in.Liqi.Seat))
	}
}
func (m *Majsoul) ActionChiPengGang(in *message.ActionChiPengGang) {
	logger.Debug("ActionChiPengGang", zap.Reflect("in", in))
	f := 0
	for _, from := range in.Froms {
		if from != in.Seat {
			f = int(from)
		}
	}
	switch in.Type + 1 {
	case majsoul.CHI:
		m.UAkochan.Chi(int(in.Seat), f, GetUAkochanTiles(in.Tiles)[0], GetUAkochanTiles(in.Tiles))
	case majsoul.PON:
		m.UAkochan.Pon(int(in.Seat), f, GetUAkochanTiles(in.Tiles)[0], GetUAkochanTiles(in.Tiles))
	case majsoul.MINKAN:
		m.UAkochan.Daiminkan(int(in.Seat), f, GetUAkochanTiles(in.Tiles)[0], GetUAkochanTiles(in.Tiles))
	default:
		logger.Error("ActionChiPengGang in.Type not found", zap.Reflect("in", in))
	}
	if in.Liqi != nil {
		m.ReachAccepted(int(in.Liqi.Seat))
	}
}
func (m *Majsoul) ActionGangResult(in *message.ActionGangResult)       {}
func (m *Majsoul) ActionGangResultEnd(in *message.ActionGangResultEnd) {}
func (m *Majsoul) ActionAnGangAddGang(in *message.ActionAnGangAddGang) {
	consumed := make([]string, 4)
	for i := 0; i < len(consumed); i++ {
		consumed[i] = GetUAkochanTile(in.Tiles)
	}
	switch in.Type {
	case majsoul.ANKAN:
		m.UAkochan.Ankan(int(in.Seat), consumed)
	case majsoul.KAKAN:
		m.UAkochan.Kakan(int(in.Seat), GetSoulTile(in.Tiles), consumed)
	default:
		logger.Debug("ActionAnGangAddGang in.Type not found", zap.Reflect("in", in))
	}
}
func (m *Majsoul) ActionBaBei(in *message.ActionBaBei) {}
func (m *Majsoul) ActionHule(in *message.ActionHule) {
	f := in.Hules[0].Seat
	m.UAkochan.Hora_(int(f), int(in.Hules[0].Seat))
	m.UAkochan.EndKyoku()
	_, err := m.ConfirmNewRound(m.Ctx, &message.ReqCommon{})
	if err != nil {
		logger.Error("ActionHule", zap.Error(err))
		return
	}
}
func (m *Majsoul) ActionHuleXueZhanMid(in *message.ActionHuleXueZhanMid) {}
func (m *Majsoul) ActionHuleXueZhanEnd(in *message.ActionHuleXueZhanEnd) {}
func (m *Majsoul) ActionLiuJu(in *message.ActionLiuJu) {
	_, err := m.ConfirmNewRound(m.Ctx, &message.ReqCommon{})
	if err != nil {
		logger.Error("ActionLiuJu", zap.Error(err))
		return
	}
}
func (m *Majsoul) ActionNoTile(in *message.ActionNoTile) {
	m.UAkochan.Ryukyoku_()
	m.UAkochan.EndKyoku()
}
