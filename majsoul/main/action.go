package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"majsoul/message"
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
		// 血战换牌
	case "ActionNewRound":
		m.ActionNewRound(ret.(*message.ActionNewRound))
	case "ActionSelectGap":
	case "ActionChangeTile":
	case "ActionRevealTile":
	case "ActionUnveilTile":
	case "ActionLockTile":
	case "ActionDiscardTile":
	case "ActionDealTile":
	case "ActionChiPengGang":
	case "ActionGangResult":
	case "ActionGangResultEnd":
	case "ActionAnGangAddGang":
	case "ActionBaBei":
	case "ActionHule":
	case "ActionHuleXueZhanMid":
	case "ActionHuleXueZhanEnd":
	case "ActionLiuJu":

	case "ActionNoTile":
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
		names[i] = fmt.Sprintf("Bot%d\n", i+1)
	}

	m.UAkochan.StartGame(id, names)
	logger.Debug("ActionMJStart")
}
func (m *Majsoul) ActionNewCard(in *message.ActionNewCard) {}
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
			tehais[i] = uakochan.GetAkoChanTiles(in.Tiles[:13])
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
	m.Tsumo(actor, in.Tiles[13])
}
func (m *Majsoul) ActionSelectGap(in *message.ActionSelectGap)         {}
func (m *Majsoul) ActionChangeTile(in *message.ActionChangeTile)       {}
func (m *Majsoul) ActionRevealTile(in *message.ActionRevealTile)       {}
func (m *Majsoul) ActionUnveilTile(in *message.ActionUnveilTile)       {}
func (m *Majsoul) ActionLockTile(in *message.ActionLockTile)           {}
func (m *Majsoul) ActionDiscardTile(in *message.ActionDiscardTile)     {}
func (m *Majsoul) ActionDealTile(in *message.ActionDealTile)           {}
func (m *Majsoul) ActionChiPengGang(in *message.ActionChiPengGang)     {}
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
