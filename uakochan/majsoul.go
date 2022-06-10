package uakochan

import (
	"encoding/json"
	"go.uber.org/zap"
	"majsoul"
	"majsoul/message"
	"utils/logger"
)

type Majsoul struct {
	*majsoul.Majsoul
}

func (m *Majsoul) NotifyCaptcha(notify *message.NotifyCaptcha) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyRoomGameStart(notify *message.NotifyRoomGameStart) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyMatchGameStart(notify *message.NotifyMatchGameStart) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyRoomPlayerReady(notify *message.NotifyRoomPlayerReady) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyRoomPlayerDressing(notify *message.NotifyRoomPlayerDressing) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyRoomPlayerUpdate(notify *message.NotifyRoomPlayerUpdate) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyRoomKickOut(notify *message.NotifyRoomKickOut) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyFriendStateChange(notify *message.NotifyFriendStateChange) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyFriendViewChange(notify *message.NotifyFriendViewChange) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyFriendChange(notify *message.NotifyFriendChange) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyNewFriendApply(notify *message.NotifyNewFriendApply) {
	//TODO implement me
	panic("implement me")
}

type InvitationRoom struct {
	RoomID    int    `json:"room_id"`
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

	//invitationRoom.RoomID

	//notify.Content
}

func (m *Majsoul) NotifyAccountUpdate(notify *message.NotifyAccountUpdate) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyAnotherLogin(notify *message.NotifyAnotherLogin) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyAccountLogout(notify *message.NotifyAccountLogout) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyAnnouncementUpdate(notify *message.NotifyAnnouncementUpdate) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyNewMail(notify *message.NotifyNewMail) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyDeleteMail(notify *message.NotifyDeleteMail) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyReviveCoinUpdate(notify *message.NotifyReviveCoinUpdate) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyDailyTaskUpdate(notify *message.NotifyDailyTaskUpdate) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivityTaskUpdate(notify *message.NotifyActivityTaskUpdate) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivityPeriodTaskUpdate(notify *message.NotifyActivityPeriodTaskUpdate) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyAccountRandomTaskUpdate(notify *message.NotifyAccountRandomTaskUpdate) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivitySegmentTaskUpdate(notify *message.NotifyActivitySegmentTaskUpdate) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivityUpdate(notify *message.NotifyActivityUpdate) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyAccountChallengeTaskUpdate(notify *message.NotifyAccountChallengeTaskUpdate) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyNewComment(notify *message.NotifyNewComment) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyRollingNotice(notify *message.NotifyRollingNotice) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyGiftSendRefresh(notify *message.NotifyGiftSendRefresh) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyShopUpdate(notify *message.NotifyShopUpdate) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyVipLevelChange(notify *message.NotifyVipLevelChange) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyServerSetting(notify *message.NotifyServerSetting) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyPayResult(notify *message.NotifyPayResult) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyCustomContestAccountMsg(notify *message.NotifyCustomContestAccountMsg) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyCustomContestSystemMsg(notify *message.NotifyCustomContestSystemMsg) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyMatchTimeout(notify *message.NotifyMatchTimeout) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyCustomContestState(notify *message.NotifyCustomContestState) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivityChange(notify *message.NotifyActivityChange) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyAFKResult(notify *message.NotifyAFKResult) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyGameFinishRewardV2(notify *message.NotifyGameFinishRewardV2) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivityRewardV2(notify *message.NotifyActivityRewardV2) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivityPointV2(notify *message.NotifyActivityPointV2) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyLeaderboardPointV2(notify *message.NotifyLeaderboardPointV2) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyNewGame(notify *message.NotifyNewGame) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyPlayerLoadGameReady(notify *message.NotifyPlayerLoadGameReady) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyGameBroadcast(notify *message.NotifyGameBroadcast) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyGameEndResult(notify *message.NotifyGameEndResult) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyGameTerminate(notify *message.NotifyGameTerminate) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyPlayerConnectionState(notify *message.NotifyPlayerConnectionState) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyAccountLevelChange(notify *message.NotifyAccountLevelChange) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyGameFinishReward(notify *message.NotifyGameFinishReward) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivityReward(notify *message.NotifyActivityReward) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivityPoint(notify *message.NotifyActivityPoint) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyLeaderboardPoint(notify *message.NotifyLeaderboardPoint) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyGamePause(notify *message.NotifyGamePause) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyEndGameVote(notify *message.NotifyEndGameVote) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyObserveData(notify *message.NotifyObserveData) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyRoomPlayerReady_AccountReadyState(notify *message.NotifyRoomPlayerReady_AccountReadyState) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyRoomPlayerDressing_AccountDressingState(notify *message.NotifyRoomPlayerDressing_AccountDressingState) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyAnnouncementUpdate_AnnouncementUpdate(notify *message.NotifyAnnouncementUpdate_AnnouncementUpdate) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivityUpdate_FeedActivityData(notify *message.NotifyActivityUpdate_FeedActivityData) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivityUpdate_FeedActivityData_CountWithTimeData(notify *message.NotifyActivityUpdate_FeedActivityData_CountWithTimeData) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivityUpdate_FeedActivityData_GiftBoxData(notify *message.NotifyActivityUpdate_FeedActivityData_GiftBoxData) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyPayResult_ResourceModify(notify *message.NotifyPayResult_ResourceModify) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyGameFinishRewardV2_LevelChange(notify *message.NotifyGameFinishRewardV2_LevelChange) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyGameFinishRewardV2_MatchChest(notify *message.NotifyGameFinishRewardV2_MatchChest) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyGameFinishRewardV2_MainCharacter(notify *message.NotifyGameFinishRewardV2_MainCharacter) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyGameFinishRewardV2_CharacterGift(notify *message.NotifyGameFinishRewardV2_CharacterGift) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivityRewardV2_ActivityReward(notify *message.NotifyActivityRewardV2_ActivityReward) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivityPointV2_ActivityPoint(notify *message.NotifyActivityPointV2_ActivityPoint) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyLeaderboardPointV2_LeaderboardPoint(notify *message.NotifyLeaderboardPointV2_LeaderboardPoint) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyGameFinishReward_LevelChange(notify *message.NotifyGameFinishReward_LevelChange) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyGameFinishReward_MatchChest(notify *message.NotifyGameFinishReward_MatchChest) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyGameFinishReward_MainCharacter(notify *message.NotifyGameFinishReward_MainCharacter) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyGameFinishReward_CharacterGift(notify *message.NotifyGameFinishReward_CharacterGift) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivityReward_ActivityReward(notify *message.NotifyActivityReward_ActivityReward) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyActivityPoint_ActivityPoint(notify *message.NotifyActivityPoint_ActivityPoint) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyLeaderboardPoint_LeaderboardPoint(notify *message.NotifyLeaderboardPoint_LeaderboardPoint) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) NotifyEndGameVote_VoteResult(notify *message.NotifyEndGameVote_VoteResult) {
	//TODO implement me
	panic("implement me")
}

func (m *Majsoul) ActionPrototype(notify *message.ActionPrototype) {
	//TODO implement me
	panic("implement me")
}
