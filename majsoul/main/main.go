package main

import (
	"encoding/json"
	"go.uber.org/zap"
	"majsoul"
	"majsoul/message"
	"uakochan"
	"utils/logger"
)

func main() {
	m := NewMajsoul()

	loginRes, err := m.Login(m.Ctx, &message.ReqLogin{
		Account:   "1601198895@qq.com",
		Password:  majsoul.Hash("miku39.."),
		Reconnect: false,
		Device: &message.ClientDeviceInfo{
			Platform:       "pc",
			Hardware:       "pc",
			Os:             "windows",
			OsVersion:      "win10",
			IsBrowser:      true,
			Software:       "Chrome",
			SalePlatform:   "web",
			HardwareVendor: "",
			ModelNumber:    "",
			ScreenWidth:    914,
			ScreenHeight:   1316,
		},
		RandomKey: "cfc35be-f519-4cbc-9765-c4c124cdc6a16",
		ClientVersion: &message.ClientVersionInfo{
			Resource: m.Version.Version,
			Package:  "",
		},
		GenAccessToken:    true,
		CurrencyPlatforms: []uint32{2, 6, 8, 10, 11},
		// 电话1 邮箱0
		Type:                0,
		Version:             0,
		ClientVersionString: m.Version.ClientVersionString,
	})
	if err != nil {
		return
	}
	logger.Debug("Login success", zap.Reflect("Nickname", loginRes.Account.Nickname))

	select {
	case <-m.Ctx.Done():
	}
}

type Majsoul struct {
	*majsoul.Majsoul
	*uakochan.UAkochan
}

func NewMajsoul() *Majsoul {
	cfg := majsoul.LoadConfig()
	m := &Majsoul{Majsoul: majsoul.New(cfg), UAkochan: uakochan.New()}
	m.IFReceive = m
	return m
}

func (m *Majsoul) NotifyCaptcha(notify *message.NotifyCaptcha) {
	logger.Debug("NotifyCaptcha", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyRoomGameStart(notify *message.NotifyRoomGameStart) {
	logger.Debug("NotifyRoomGameStart", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyMatchGameStart(notify *message.NotifyMatchGameStart) {
	logger.Debug("NotifyMatchGameStart", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyRoomPlayerReady(notify *message.NotifyRoomPlayerReady) {
	logger.Debug("NotifyRoomPlayerReady", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyRoomPlayerDressing(notify *message.NotifyRoomPlayerDressing) {
	logger.Debug("NotifyRoomPlayerDressing", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyRoomPlayerUpdate(notify *message.NotifyRoomPlayerUpdate) {
	logger.Debug("NotifyRoomPlayerUpdate", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyRoomKickOut(notify *message.NotifyRoomKickOut) {
	logger.Debug("NotifyRoomKickOut", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyFriendStateChange(notify *message.NotifyFriendStateChange) {
	logger.Debug("NotifyFriendStateChange", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyFriendViewChange(notify *message.NotifyFriendViewChange) {
	logger.Debug("NotifyFriendViewChange", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyFriendChange(notify *message.NotifyFriendChange) {
	logger.Debug("NotifyFriendChange", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyNewFriendApply(notify *message.NotifyNewFriendApply) {
	logger.Debug("NotifyNewFriendApply", zap.Reflect("notify", notify))
}

type InvitationRoom struct {
	RoomID    uint32 `json:"room_id"`
	Mode      Mode   `json:"mode"`
	Nickname  string `json:"nickname"`
	Verified  int    `json:"verified"`
	AccountID int    `json:"account_id"`
}
type Mode struct {
	Mode       int        `json:"mode"`
	Ai         bool       `json:"ai"`
	DetailRule DetailRule `json:"detail_rule"`
}
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

func (m *Majsoul) NotifyClientMessage(notify *message.NotifyClientMessage) {
	if notify.Type != 1 {
		logger.Info("NotifyClientMessage", zap.Reflect("notify", notify))
		return
	}
	invitationRoom := new(InvitationRoom)
	err := json.Unmarshal([]byte(notify.Content), invitationRoom)
	if err != nil {
		logger.Error("NotifyClientMessage", zap.Error(err))
		return
	}
	joinRoom, err := m.JoinRoom(m.Ctx, &message.ReqJoinRoom{
		RoomId:              invitationRoom.RoomID,
		ClientVersionString: m.ClientVersionString,
	})
	if err != nil {
		logger.Error("NotifyClientMessage", zap.Error(err))
		return
	}
	logger.Debug("Join Room", zap.Reflect("res", joinRoom))
}

func (m *Majsoul) NotifyAccountUpdate(notify *message.NotifyAccountUpdate) {
	logger.Debug("NotifyAccountUpdate", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyAnotherLogin(notify *message.NotifyAnotherLogin) {
	logger.Debug("NotifyAnotherLogin", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyAccountLogout(notify *message.NotifyAccountLogout) {
	logger.Debug("NotifyAccountLogout", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyAnnouncementUpdate(notify *message.NotifyAnnouncementUpdate) {
	logger.Debug("NotifyAnnouncementUpdate", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyNewMail(notify *message.NotifyNewMail) {
	logger.Debug("NotifyNewMail", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyDeleteMail(notify *message.NotifyDeleteMail) {
	logger.Debug("NotifyDeleteMail", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyReviveCoinUpdate(notify *message.NotifyReviveCoinUpdate) {
	logger.Debug("NotifyReviveCoinUpdate", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyDailyTaskUpdate(notify *message.NotifyDailyTaskUpdate) {
	logger.Debug("NotifyDailyTaskUpdate", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivityTaskUpdate(notify *message.NotifyActivityTaskUpdate) {
	logger.Debug("NotifyActivityTaskUpdate", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivityPeriodTaskUpdate(notify *message.NotifyActivityPeriodTaskUpdate) {
	logger.Debug("NotifyActivityPeriodTaskUpdate", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyAccountRandomTaskUpdate(notify *message.NotifyAccountRandomTaskUpdate) {
	logger.Debug("NotifyAccountRandomTaskUpdate", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivitySegmentTaskUpdate(notify *message.NotifyActivitySegmentTaskUpdate) {
	logger.Debug("NotifyActivitySegmentTaskUpdate", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivityUpdate(notify *message.NotifyActivityUpdate) {
	logger.Debug("NotifyActivityUpdate", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyAccountChallengeTaskUpdate(notify *message.NotifyAccountChallengeTaskUpdate) {
	logger.Debug("NotifyAccountChallengeTaskUpdate", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyNewComment(notify *message.NotifyNewComment) {
	logger.Debug("NotifyNewComment", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyRollingNotice(notify *message.NotifyRollingNotice) {
	logger.Debug("NotifyRollingNotice", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyGiftSendRefresh(notify *message.NotifyGiftSendRefresh) {
	logger.Debug("NotifyGiftSendRefresh", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyShopUpdate(notify *message.NotifyShopUpdate) {
	logger.Debug("NotifyShopUpdate", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyVipLevelChange(notify *message.NotifyVipLevelChange) {
	logger.Debug("NotifyVipLevelChange", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyServerSetting(notify *message.NotifyServerSetting) {
	logger.Debug("NotifyServerSetting", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyPayResult(notify *message.NotifyPayResult) {
	logger.Debug("NotifyPayResult", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyCustomContestAccountMsg(notify *message.NotifyCustomContestAccountMsg) {
	logger.Debug("NotifyCustomContestAccountMsg", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyCustomContestSystemMsg(notify *message.NotifyCustomContestSystemMsg) {
	logger.Debug("NotifyCustomContestSystemMsg", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyMatchTimeout(notify *message.NotifyMatchTimeout) {
	logger.Debug("NotifyMatchTimeout", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyCustomContestState(notify *message.NotifyCustomContestState) {
	logger.Debug("NotifyCustomContestState", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivityChange(notify *message.NotifyActivityChange) {
	logger.Debug("NotifyActivityChange", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyAFKResult(notify *message.NotifyAFKResult) {
	logger.Debug("NotifyAFKResult", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyGameFinishRewardV2(notify *message.NotifyGameFinishRewardV2) {
	logger.Debug("NotifyGameFinishRewardV2", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivityRewardV2(notify *message.NotifyActivityRewardV2) {
	logger.Debug("NotifyActivityRewardV2", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivityPointV2(notify *message.NotifyActivityPointV2) {
	logger.Debug("NotifyActivityPointV2", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyLeaderboardPointV2(notify *message.NotifyLeaderboardPointV2) {
	logger.Debug("NotifyLeaderboardPointV2", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyNewGame(notify *message.NotifyNewGame) {
	logger.Debug("NotifyNewGame", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyPlayerLoadGameReady(notify *message.NotifyPlayerLoadGameReady) {
	logger.Debug("NotifyPlayerLoadGameReady", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyGameBroadcast(notify *message.NotifyGameBroadcast) {
	logger.Debug("NotifyGameBroadcast", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyGameEndResult(notify *message.NotifyGameEndResult) {
	logger.Debug("NotifyGameEndResult", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyGameTerminate(notify *message.NotifyGameTerminate) {
	logger.Debug("NotifyGameTerminate", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyPlayerConnectionState(notify *message.NotifyPlayerConnectionState) {
	logger.Debug("NotifyPlayerConnectionState", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyAccountLevelChange(notify *message.NotifyAccountLevelChange) {
	logger.Debug("NotifyAccountLevelChange", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyGameFinishReward(notify *message.NotifyGameFinishReward) {
	logger.Debug("NotifyGameFinishReward", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivityReward(notify *message.NotifyActivityReward) {
	logger.Debug("NotifyActivityReward", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivityPoint(notify *message.NotifyActivityPoint) {
	logger.Debug("NotifyActivityPoint", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyLeaderboardPoint(notify *message.NotifyLeaderboardPoint) {
	logger.Debug("NotifyLeaderboardPoint", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyGamePause(notify *message.NotifyGamePause) {
	logger.Debug("NotifyGamePause", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyEndGameVote(notify *message.NotifyEndGameVote) {
	logger.Debug("NotifyEndGameVote", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyObserveData(notify *message.NotifyObserveData) {
	logger.Debug("NotifyObserveData", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyRoomPlayerReady_AccountReadyState(notify *message.NotifyRoomPlayerReady_AccountReadyState) {
	logger.Debug("NotifyRoomPlayerReady_AccountReadyState", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyRoomPlayerDressing_AccountDressingState(notify *message.NotifyRoomPlayerDressing_AccountDressingState) {
	logger.Debug("NotifyRoomPlayerDressing_AccountDressingState", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyAnnouncementUpdate_AnnouncementUpdate(notify *message.NotifyAnnouncementUpdate_AnnouncementUpdate) {
	logger.Debug("NotifyAnnouncementUpdate_AnnouncementUpdate", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivityUpdate_FeedActivityData(notify *message.NotifyActivityUpdate_FeedActivityData) {
	logger.Debug("NotifyActivityUpdate_FeedActivityData", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivityUpdate_FeedActivityData_CountWithTimeData(notify *message.NotifyActivityUpdate_FeedActivityData_CountWithTimeData) {
	logger.Debug("NotifyActivityUpdate_FeedActivityData_CountWithTimeData", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivityUpdate_FeedActivityData_GiftBoxData(notify *message.NotifyActivityUpdate_FeedActivityData_GiftBoxData) {
	logger.Debug("NotifyActivityUpdate_FeedActivityData_GiftBoxData", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyPayResult_ResourceModify(notify *message.NotifyPayResult_ResourceModify) {
	logger.Debug("NotifyPayResult_ResourceModify", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyGameFinishRewardV2_LevelChange(notify *message.NotifyGameFinishRewardV2_LevelChange) {
	logger.Debug("NotifyGameFinishRewardV2_LevelChange", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyGameFinishRewardV2_MatchChest(notify *message.NotifyGameFinishRewardV2_MatchChest) {
	logger.Debug("NotifyGameFinishRewardV2_MatchChest", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyGameFinishRewardV2_MainCharacter(notify *message.NotifyGameFinishRewardV2_MainCharacter) {
	logger.Debug("NotifyGameFinishRewardV2_MainCharacter", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyGameFinishRewardV2_CharacterGift(notify *message.NotifyGameFinishRewardV2_CharacterGift) {
	logger.Debug("NotifyGameFinishRewardV2_CharacterGift", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivityRewardV2_ActivityReward(notify *message.NotifyActivityRewardV2_ActivityReward) {
	logger.Debug("NotifyActivityRewardV2_ActivityReward", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivityPointV2_ActivityPoint(notify *message.NotifyActivityPointV2_ActivityPoint) {
	logger.Debug("NotifyActivityPointV2_ActivityPoint", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyLeaderboardPointV2_LeaderboardPoint(notify *message.NotifyLeaderboardPointV2_LeaderboardPoint) {
	logger.Debug("NotifyLeaderboardPointV2_LeaderboardPoint", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyGameFinishReward_LevelChange(notify *message.NotifyGameFinishReward_LevelChange) {
	logger.Debug("NotifyGameFinishReward_LevelChange", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyGameFinishReward_MatchChest(notify *message.NotifyGameFinishReward_MatchChest) {
	logger.Debug("NotifyGameFinishReward_MatchChest", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyGameFinishReward_MainCharacter(notify *message.NotifyGameFinishReward_MainCharacter) {
	logger.Debug("NotifyGameFinishReward_MainCharacter", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyGameFinishReward_CharacterGift(notify *message.NotifyGameFinishReward_CharacterGift) {
	logger.Debug("NotifyGameFinishReward_CharacterGift", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivityReward_ActivityReward(notify *message.NotifyActivityReward_ActivityReward) {
	logger.Debug("NotifyActivityReward_ActivityReward", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyActivityPoint_ActivityPoint(notify *message.NotifyActivityPoint_ActivityPoint) {
	logger.Debug("NotifyActivityPoint_ActivityPoint", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyLeaderboardPoint_LeaderboardPoint(notify *message.NotifyLeaderboardPoint_LeaderboardPoint) {
	logger.Debug("NotifyLeaderboardPoint_LeaderboardPoint", zap.Reflect("notify", notify))
}

func (m *Majsoul) NotifyEndGameVote_VoteResult(notify *message.NotifyEndGameVote_VoteResult) {
	logger.Debug("NotifyEndGameVote_VoteResult", zap.Reflect("notify", notify))
}

func (m *Majsoul) ActionPrototype(notify *message.ActionPrototype) {
	logger.Debug("ActionPrototype", zap.Reflect("notify", notify))
}
