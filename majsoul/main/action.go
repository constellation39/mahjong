package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"majsoul/message"
	"time"
	"uakochan"
	"utils/logger"
)

func (m *Majsoul) ActionPrototype(notify *message.ActionPrototype) {
	logger.Debug("ActionPrototype", zap.Reflect("notify", notify))
	ret := message.GetActionType(notify.Name)
	err := proto.Unmarshal(notify.Data, ret)
	if err != nil {
		logger.Error("ActionPrototype", zap.Error(err))
		return
	}
	logger.Debug("ActionPrototype", zap.Reflect("ret", ret))
	switch notify.Name {
	case "ActionMJStart":
		m.ActionMJStart(ret.(*message.ActionMJStart))
	case "ActionNewCard":
		m.ActionNewCard(ret.(*message.ActionNewCard))
	case "ActionNewRound":
		m.ActionNewRound(ret.(*message.ActionNewRound))
	case "ActionSelectGap":
		m.ActionSelectGap(ret.(*message.ActionSelectGap))
	case "ActionChangeTile":
		m.ActionChangeTile(ret.(*message.ActionChangeTile))
	case "ActionRevealTile":
		m.ActionRevealTile(ret.(*message.ActionRevealTile))
	case "ActionUnveilTile":
		m.ActionUnveilTile(ret.(*message.ActionUnveilTile))
	case "ActionLockTile":
		m.ActionLockTile(ret.(*message.ActionLockTile))
	case "ActionDiscardTile":
		m.ActionDiscardTile(ret.(*message.ActionDiscardTile))
	case "ActionDealTile":
		m.ActionDealTile(ret.(*message.ActionDealTile))
	case "ActionChiPengGang":
		m.ActionChiPengGang(ret.(*message.ActionChiPengGang))
	case "ActionGangResult":
		m.ActionGangResult(ret.(*message.ActionGangResult))
	case "ActionGangResultEnd":
		m.ActionGangResultEnd(ret.(*message.ActionGangResultEnd))
	case "ActionAnGangAddGang":
		m.ActionAnGangAddGang(ret.(*message.ActionAnGangAddGang))
	case "ActionBaBei":
		m.ActionBaBei(ret.(*message.ActionBaBei))
	case "ActionHule":
		m.ActionHule(ret.(*message.ActionHule))
	case "ActionHuleXueZhanMid":
		m.ActionHuleXueZhanMid(ret.(*message.ActionHuleXueZhanMid))
	case "ActionHuleXueZhanEnd":
		m.ActionHuleXueZhanEnd(ret.(*message.ActionHuleXueZhanEnd))
	case "ActionLiuJu":
		m.ActionLiuJu(ret.(*message.ActionLiuJu))
	case "ActionNoTile":
		m.ActionNoTile(ret.(*message.ActionNoTile))
	}
}

func (m *Majsoul) ActionMJStart(in *message.ActionMJStart) {
	id := 0
	names := make([]string, 4)
	playerIndex := 0
	for i, player := range m.GameInfo.Players {
		playerIndex = i
		if player.AccountId == m.Account.AccountId {
			id = i
		}
		names[i] = player.Nickname
	}

	for i := playerIndex + 1; i < 4; i++ {
		names[i] = fmt.Sprintf("Bot%d", i+1)
	}

	m.UAkochan.StartGame(id, names)
	logger.Debug("ActionMJStart", zap.Reflect("in", in))
}
func (m *Majsoul) ActionNewCard(in *message.ActionNewCard) {
	logger.Debug("ActionNewCard", zap.Reflect("in", in))
}
func (m *Majsoul) ActionNewRound(in *message.ActionNewRound) {
	bakaze := uakochan.GetBakaze(in.Chang)
	kyoku := in.Ju + 1
	honba := in.Ben
	kyotaku := in.Liqibang
	oya := in.Ju
	dora_marker := in.Doras[0]
	tehais := make([][]string, 4)
	playerIndex := 0
	isOya := true
	actor := 0
	for i, player := range m.GameInfo.Players {
		playerIndex = i
		if player.AccountId == m.Account.AccountId {
			actor = i
			isOya = uint32(i) == oya
			tehais[i] = uakochan.GetTiles(in.Tiles[:13])
			continue
		}
		tehais[i] = []string{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"}
	}
	for i := playerIndex + 1; i < 4; i++ {
		tehais[i] = []string{"?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?", "?"}
	}
	m.UAkochan.StartKyoku(bakaze, kyoku, honba, kyotaku, oya, dora_marker, tehais)
	if !isOya {
		return
	}
	time.Sleep(time.Second * 3)
	m.Tsumo(actor, in.Tiles[13], in.Operation)
	logger.Debug("ActionNewRound", zap.Reflect("in", in))
}
func (m *Majsoul) ActionSelectGap(in *message.ActionSelectGap) {
	logger.Debug("ActionSelectGap", zap.Reflect("in", in))
}
func (m *Majsoul) ActionChangeTile(in *message.ActionChangeTile) {
	logger.Debug("ActionChangeTile", zap.Reflect("in", in))
}
func (m *Majsoul) ActionRevealTile(in *message.ActionRevealTile) {
	logger.Debug("ActionRevealTile", zap.Reflect("in", in))
}
func (m *Majsoul) ActionUnveilTile(in *message.ActionUnveilTile) {
	logger.Debug("ActionUnveilTile", zap.Reflect("in", in))
}
func (m *Majsoul) ActionLockTile(in *message.ActionLockTile) {
	logger.Debug("ActionLockTile", zap.Reflect("in", in))
}
func (m *Majsoul) ActionDiscardTile(in *message.ActionDiscardTile) {
	logger.Debug("ActionDiscardTile", zap.Reflect("in", in))
	m.Dahai(int(in.Seat), in.Tile, in.Moqie, in.Operation)
}
func (m *Majsoul) ActionDealTile(in *message.ActionDealTile) {
	if in.Tile == "" {
		in.Tile = "?"
	}
	m.Tsumo(int(in.Seat), in.Tile, in.Operation)
}
func (m *Majsoul) ActionChiPengGang(in *message.ActionChiPengGang) {
	logger.Debug("ActionChiPengGang", zap.Reflect("in", in))
	//switch in.Type {
	//case majsoul.CHI:
	//	m.UAkochan.Chi(int(in.Seat), in, uakochan.GetTiles(in.Tiles))
	//case majsoul.PON:
	//	m.UAkochan.Pon()
	//case majsoul.ANKAN:
	//	m.UAkochan.Ankan()
	//case majsoul.MINKAN:
	//	m.UAkochan.Daiminkan()
	//case majsoul.KAKAN:
	//	m.UAkochan.Kakan()
	//}
}
func (m *Majsoul) ActionGangResult(in *message.ActionGangResult)       {}
func (m *Majsoul) ActionGangResultEnd(in *message.ActionGangResultEnd) {}
func (m *Majsoul) ActionAnGangAddGang(in *message.ActionAnGangAddGang) {}
func (m *Majsoul) ActionBaBei(in *message.ActionBaBei)                 {}
func (m *Majsoul) ActionHule(in *message.ActionHule) {
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
func (m *Majsoul) ActionNoTile(in *message.ActionNoTile) {}
