package main

import (
	"encoding/json"
	"go.uber.org/zap"
	"majsoul/message"
	"utils/logger"
)

func (m *Majsoul) NotifyClientMessage(notify *message.NotifyClientMessage) {
	type DetailRule struct {
		TimeFixed    int  `json:"time_fixed"`
		TimeAdd      int  `json:"time_add"`
		DoraCount    int  `json:"dora_count"`
		Shiduan      int  `json:"shiduan"`
		InitPoint    int  `json:"init_point"`
		Fandian      int  `json:"fandian"`
		Bianjietishi bool `json:"bianjietishi"`
		AiLevel      int  `json:"ai_level"`
		Fanfu        int  `json:"fanfu"`
		GuyiMode     int  `json:"guyi_mode"`
		OpenHand     int  `json:"open_hand"`
	}
	type Mode struct {
		Mode       int        `json:"mode"`
		Ai         bool       `json:"ai"`
		DetailRule DetailRule `json:"detail_rule"`
	}
	type InvitationRoom struct {
		RoomID    uint32 `json:"room_id"`
		Mode      Mode   `json:"mode"`
		Nickname  string `json:"nickname"`
		Verified  int    `json:"verified"`
		AccountID int    `json:"account_id"`
	}
	if notify.Type != 1 {
		logger.Info("Majsoul.NotifyClientMessage", zap.Reflect("notify", notify))
		return
	}
	invitationRoom := new(InvitationRoom)
	err := json.Unmarshal([]byte(notify.Content), invitationRoom)
	if err != nil {
		logger.Error("Majsoul.NotifyClientMessage", zap.Error(err))
		return
	}
	joinRoom, err := m.JoinRoom(m.Ctx, &message.ReqJoinRoom{
		RoomId:              invitationRoom.RoomID,
		ClientVersionString: m.Version.String(),
	})
	if err != nil {
		logger.Error("Majsoul.NotifyClientMessage", zap.Error(err))
		return
	}
	logger.Debug("Majsoul.NotifyClientMessage", zap.Reflect("JoinRoom", joinRoom))
	readyPlay, err := m.ReadyPlay(m.Ctx, &message.ReqRoomReady{Ready: true})
	if err != nil {
		return
	}
	logger.Debug("Majsoul.NotifyClientMessage", zap.Reflect("ReadyPlay", readyPlay))
}

func (m *Majsoul) NotifyEndGameVote(notify *message.NotifyEndGameVote) {
	end, err := m.VoteGameEnd(m.Ctx, &message.ReqVoteGameEnd{Yes: true})
	if err != nil {
		return
	}
	logger.Debug("VoteGameEnd", zap.Reflect("end", end))
}
