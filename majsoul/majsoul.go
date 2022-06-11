package majsoul

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"majsoul/message"
	"math/rand"
	"time"
	"utils/config"
	"utils/logger"
	"utils/net"
)

const (
	MsgTypeNotify   uint8 = 1
	MsgTypeRequest  uint8 = 2
	MsgTypeResponse uint8 = 3
)

const (
	DISCARD = 1
	CHI     = 2
	PON     = 3
	ANKAN   = 4
	MINKAN  = 5
	KAKAN   = 6
	RIICHI  = 7
	TSUMO   = 8
	RON     = 9
	KUKU    = 10
	KITA    = 11
	PASS    = 12
)

type Config struct {
	ServerAddress  string `json:"serverAddress"`
	GatewayAddress string `json:"gatewayAddress"`
	GameAddress    string `json:"gameAddress"`
}

type IFReceive interface {
	NotifyCaptcha(notify *message.NotifyCaptcha)
	NotifyRoomGameStart(notify *message.NotifyRoomGameStart)
	NotifyMatchGameStart(notify *message.NotifyMatchGameStart)
	NotifyRoomPlayerReady(notify *message.NotifyRoomPlayerReady)
	NotifyRoomPlayerDressing(notify *message.NotifyRoomPlayerDressing)
	NotifyRoomPlayerUpdate(notify *message.NotifyRoomPlayerUpdate)
	NotifyRoomKickOut(notify *message.NotifyRoomKickOut)
	NotifyFriendStateChange(notify *message.NotifyFriendStateChange)
	NotifyFriendViewChange(notify *message.NotifyFriendViewChange)
	NotifyFriendChange(notify *message.NotifyFriendChange)
	NotifyNewFriendApply(notify *message.NotifyNewFriendApply)
	NotifyClientMessage(notify *message.NotifyClientMessage)
	NotifyAccountUpdate(notify *message.NotifyAccountUpdate)
	NotifyAnotherLogin(notify *message.NotifyAnotherLogin)
	NotifyAccountLogout(notify *message.NotifyAccountLogout)
	NotifyAnnouncementUpdate(notify *message.NotifyAnnouncementUpdate)
	NotifyNewMail(notify *message.NotifyNewMail)
	NotifyDeleteMail(notify *message.NotifyDeleteMail)
	NotifyReviveCoinUpdate(notify *message.NotifyReviveCoinUpdate)
	NotifyDailyTaskUpdate(notify *message.NotifyDailyTaskUpdate)
	NotifyActivityTaskUpdate(notify *message.NotifyActivityTaskUpdate)
	NotifyActivityPeriodTaskUpdate(notify *message.NotifyActivityPeriodTaskUpdate)
	NotifyAccountRandomTaskUpdate(notify *message.NotifyAccountRandomTaskUpdate)
	NotifyActivitySegmentTaskUpdate(notify *message.NotifyActivitySegmentTaskUpdate)
	NotifyActivityUpdate(notify *message.NotifyActivityUpdate)
	NotifyAccountChallengeTaskUpdate(notify *message.NotifyAccountChallengeTaskUpdate)
	NotifyNewComment(notify *message.NotifyNewComment)
	NotifyRollingNotice(notify *message.NotifyRollingNotice)
	NotifyGiftSendRefresh(notify *message.NotifyGiftSendRefresh)
	NotifyShopUpdate(notify *message.NotifyShopUpdate)
	NotifyVipLevelChange(notify *message.NotifyVipLevelChange)
	NotifyServerSetting(notify *message.NotifyServerSetting)
	NotifyPayResult(notify *message.NotifyPayResult)
	NotifyCustomContestAccountMsg(notify *message.NotifyCustomContestAccountMsg)
	NotifyCustomContestSystemMsg(notify *message.NotifyCustomContestSystemMsg)
	NotifyMatchTimeout(notify *message.NotifyMatchTimeout)
	NotifyCustomContestState(notify *message.NotifyCustomContestState)
	NotifyActivityChange(notify *message.NotifyActivityChange)
	NotifyAFKResult(notify *message.NotifyAFKResult)
	NotifyGameFinishRewardV2(notify *message.NotifyGameFinishRewardV2)
	NotifyActivityRewardV2(notify *message.NotifyActivityRewardV2)
	NotifyActivityPointV2(notify *message.NotifyActivityPointV2)
	NotifyLeaderboardPointV2(notify *message.NotifyLeaderboardPointV2)
	NotifyNewGame(notify *message.NotifyNewGame)
	NotifyPlayerLoadGameReady(notify *message.NotifyPlayerLoadGameReady)
	NotifyGameBroadcast(notify *message.NotifyGameBroadcast)
	NotifyGameEndResult(notify *message.NotifyGameEndResult)
	NotifyGameTerminate(notify *message.NotifyGameTerminate)
	NotifyPlayerConnectionState(notify *message.NotifyPlayerConnectionState)
	NotifyAccountLevelChange(notify *message.NotifyAccountLevelChange)
	NotifyGameFinishReward(notify *message.NotifyGameFinishReward)
	NotifyActivityReward(notify *message.NotifyActivityReward)
	NotifyActivityPoint(notify *message.NotifyActivityPoint)
	NotifyLeaderboardPoint(notify *message.NotifyLeaderboardPoint)
	NotifyGamePause(notify *message.NotifyGamePause)
	NotifyEndGameVote(notify *message.NotifyEndGameVote)
	NotifyObserveData(notify *message.NotifyObserveData)
	NotifyRoomPlayerReady_AccountReadyState(notify *message.NotifyRoomPlayerReady_AccountReadyState)
	NotifyRoomPlayerDressing_AccountDressingState(notify *message.NotifyRoomPlayerDressing_AccountDressingState)
	NotifyAnnouncementUpdate_AnnouncementUpdate(notify *message.NotifyAnnouncementUpdate_AnnouncementUpdate)
	NotifyActivityUpdate_FeedActivityData(notify *message.NotifyActivityUpdate_FeedActivityData)
	NotifyActivityUpdate_FeedActivityData_CountWithTimeData(notify *message.NotifyActivityUpdate_FeedActivityData_CountWithTimeData)
	NotifyActivityUpdate_FeedActivityData_GiftBoxData(notify *message.NotifyActivityUpdate_FeedActivityData_GiftBoxData)
	NotifyPayResult_ResourceModify(notify *message.NotifyPayResult_ResourceModify)
	NotifyGameFinishRewardV2_LevelChange(notify *message.NotifyGameFinishRewardV2_LevelChange)
	NotifyGameFinishRewardV2_MatchChest(notify *message.NotifyGameFinishRewardV2_MatchChest)
	NotifyGameFinishRewardV2_MainCharacter(notify *message.NotifyGameFinishRewardV2_MainCharacter)
	NotifyGameFinishRewardV2_CharacterGift(notify *message.NotifyGameFinishRewardV2_CharacterGift)
	NotifyActivityRewardV2_ActivityReward(notify *message.NotifyActivityRewardV2_ActivityReward)
	NotifyActivityPointV2_ActivityPoint(notify *message.NotifyActivityPointV2_ActivityPoint)
	NotifyLeaderboardPointV2_LeaderboardPoint(notify *message.NotifyLeaderboardPointV2_LeaderboardPoint)
	NotifyGameFinishReward_LevelChange(notify *message.NotifyGameFinishReward_LevelChange)
	NotifyGameFinishReward_MatchChest(notify *message.NotifyGameFinishReward_MatchChest)
	NotifyGameFinishReward_MainCharacter(notify *message.NotifyGameFinishReward_MainCharacter)
	NotifyGameFinishReward_CharacterGift(notify *message.NotifyGameFinishReward_CharacterGift)
	NotifyActivityReward_ActivityReward(notify *message.NotifyActivityReward_ActivityReward)
	NotifyActivityPoint_ActivityPoint(notify *message.NotifyActivityPoint_ActivityPoint)
	NotifyLeaderboardPoint_LeaderboardPoint(notify *message.NotifyLeaderboardPoint_LeaderboardPoint)
	NotifyEndGameVote_VoteResult(notify *message.NotifyEndGameVote_VoteResult)
	ActionPrototype(notify *message.ActionPrototype)
}

type Majsoul struct {
	Ctx    context.Context
	Cancel context.CancelFunc

	message.LobbyClient
	message.FastTestClient

	IFReceive

	request *net.Request
	conn    *ClientConn
	game    *ClientConn

	Config   *Config
	Version  *Version
	Account  *message.Account
	GameInfo *message.ResAuthGame
}

func New(c *Config) *Majsoul {
	ctx, cancel := context.WithCancel(context.Background())
	cConn := NewClientConn(ctx, c.GatewayAddress)
	m := &Majsoul{
		Ctx:         ctx,
		Cancel:      cancel,
		request:     net.NewRequest(c.ServerAddress),
		LobbyClient: message.NewLobbyClient(cConn),
		Config:      c,
		conn:        cConn,
	}
	m.IFReceive = m
	m.init()
	go m.heatbeat()
	go m.receiveConn()
	return m
}

func Hash(data string) string {
	hash := hmac.New(sha256.New, []byte("lailai"))
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func LoadConfig() *Config {
	cfg := new(Config)
	err := config.Read("majsoul.json", cfg)
	if err != nil {
		logger.Panic("init client fail", zap.Error(err))
	}
	return cfg
}

func (majsoul *Majsoul) Start() {
}

type Version struct {
	Version      string `json:"version"`
	ForceVersion string `json:"force_version"`
	Code         string `json:"code"`
}

func (v *Version) String() string {
	return fmt.Sprintf("web-%s", v.Version[:len(v.Version)-2])
}

func (majsoul *Majsoul) Login(username, password string) {
	loginRes, err := majsoul.LobbyClient.Login(majsoul.Ctx, &message.ReqLogin{
		Account:   username,
		Password:  Hash(password),
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
			Resource: majsoul.Version.Version,
			Package:  "",
		},
		GenAccessToken:    true,
		CurrencyPlatforms: []uint32{2, 6, 8, 10, 11},
		// 电话1 邮箱0
		Type:                0,
		Version:             0,
		ClientVersionString: majsoul.Version.String(),
	})
	if err != nil {
		return
	}
	majsoul.Account = loginRes.Account
	logger.Debug("Login success", zap.Reflect("Nickname", loginRes.Account.Nickname))
}

func (majsoul *Majsoul) GetVersion() *Version {
	body, err := majsoul.request.Get(fmt.Sprintf("1/version.json?randv=%d", rand.Intn(1000000000)))
	if err != nil {
		logger.Panic("GetVersion", zap.Error(err))
	}
	majsoul.Version = new(Version)
	err = json.Unmarshal(body, majsoul.Version)
	if err != nil {
		logger.Panic("GetVersion", zap.Error(err))
	}
	return majsoul.Version
}

func (majsoul *Majsoul) init() {
	version := majsoul.GetVersion()
	if version.Version != "0.10.105.w" {
		logger.Info("liqi.json的版本为0.10.105.w,雀魂当前版本为", zap.String("Version", version.Version))
	}
	logger.Debug("当前雀魂版本为: ", zap.String("Version", version.Version))
}

func (majsoul *Majsoul) heatbeat() {
	t := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-t.C:
			_, err := majsoul.Heatbeat(majsoul.Ctx, &message.ReqHeatBeat{})
			if err != nil {
				logger.Info("Heatbeat", zap.Error(err))
				return
			}
		}
	}
}

func (majsoul *Majsoul) receiveConn() {
	for data := range majsoul.conn.Receive() {
		logger.Debug("majsoul.receiveConn", zap.Reflect("raw", data))
		majsoul.handleNotify(data)
	}
}

func (majsoul *Majsoul) receiveGame() {
	for data := range majsoul.game.Receive() {
		logger.Debug("majsoul.receiveGame", zap.Reflect("raw", data))
		majsoul.handleNotify(data)
	}
}

func (majsoul *Majsoul) handleNotify(data proto.Message) {
	switch notify := data.(type) {
	case *message.NotifyCaptcha:
		majsoul.IFReceive.NotifyCaptcha(notify)
	case *message.NotifyRoomGameStart:
		majsoul.IFReceive.NotifyRoomGameStart(notify)
	case *message.NotifyMatchGameStart:
		majsoul.IFReceive.NotifyMatchGameStart(notify)
	case *message.NotifyRoomPlayerReady:
		majsoul.IFReceive.NotifyRoomPlayerReady(notify)
	case *message.NotifyRoomPlayerDressing:
		majsoul.IFReceive.NotifyRoomPlayerDressing(notify)
	case *message.NotifyRoomPlayerUpdate:
		majsoul.IFReceive.NotifyRoomPlayerUpdate(notify)
	case *message.NotifyRoomKickOut:
		majsoul.IFReceive.NotifyRoomKickOut(notify)
	case *message.NotifyFriendStateChange:
		majsoul.IFReceive.NotifyFriendStateChange(notify)
	case *message.NotifyFriendViewChange:
		majsoul.IFReceive.NotifyFriendViewChange(notify)
	case *message.NotifyFriendChange:
		majsoul.IFReceive.NotifyFriendChange(notify)
	case *message.NotifyNewFriendApply:
		majsoul.IFReceive.NotifyNewFriendApply(notify)
	case *message.NotifyClientMessage:
		majsoul.IFReceive.NotifyClientMessage(notify)
	case *message.NotifyAccountUpdate:
		majsoul.IFReceive.NotifyAccountUpdate(notify)
	case *message.NotifyAnotherLogin:
		majsoul.IFReceive.NotifyAnotherLogin(notify)
	case *message.NotifyAccountLogout:
		majsoul.IFReceive.NotifyAccountLogout(notify)
	case *message.NotifyAnnouncementUpdate:
		majsoul.IFReceive.NotifyAnnouncementUpdate(notify)
	case *message.NotifyNewMail:
		majsoul.IFReceive.NotifyNewMail(notify)
	case *message.NotifyDeleteMail:
		majsoul.IFReceive.NotifyDeleteMail(notify)
	case *message.NotifyReviveCoinUpdate:
		majsoul.IFReceive.NotifyReviveCoinUpdate(notify)
	case *message.NotifyDailyTaskUpdate:
		majsoul.IFReceive.NotifyDailyTaskUpdate(notify)
	case *message.NotifyActivityTaskUpdate:
		majsoul.IFReceive.NotifyActivityTaskUpdate(notify)
	case *message.NotifyActivityPeriodTaskUpdate:
		majsoul.IFReceive.NotifyActivityPeriodTaskUpdate(notify)
	case *message.NotifyAccountRandomTaskUpdate:
		majsoul.IFReceive.NotifyAccountRandomTaskUpdate(notify)
	case *message.NotifyActivitySegmentTaskUpdate:
		majsoul.IFReceive.NotifyActivitySegmentTaskUpdate(notify)
	case *message.NotifyActivityUpdate:
		majsoul.IFReceive.NotifyActivityUpdate(notify)
	case *message.NotifyAccountChallengeTaskUpdate:
		majsoul.IFReceive.NotifyAccountChallengeTaskUpdate(notify)
	case *message.NotifyNewComment:
		majsoul.IFReceive.NotifyNewComment(notify)
	case *message.NotifyRollingNotice:
		majsoul.IFReceive.NotifyRollingNotice(notify)
	case *message.NotifyGiftSendRefresh:
		majsoul.IFReceive.NotifyGiftSendRefresh(notify)
	case *message.NotifyShopUpdate:
		majsoul.IFReceive.NotifyShopUpdate(notify)
	case *message.NotifyVipLevelChange:
		majsoul.IFReceive.NotifyVipLevelChange(notify)
	case *message.NotifyServerSetting:
		majsoul.IFReceive.NotifyServerSetting(notify)
	case *message.NotifyPayResult:
		majsoul.IFReceive.NotifyPayResult(notify)
	case *message.NotifyCustomContestAccountMsg:
		majsoul.IFReceive.NotifyCustomContestAccountMsg(notify)
	case *message.NotifyCustomContestSystemMsg:
		majsoul.IFReceive.NotifyCustomContestSystemMsg(notify)
	case *message.NotifyMatchTimeout:
		majsoul.IFReceive.NotifyMatchTimeout(notify)
	case *message.NotifyCustomContestState:
		majsoul.IFReceive.NotifyCustomContestState(notify)
	case *message.NotifyActivityChange:
		majsoul.IFReceive.NotifyActivityChange(notify)
	case *message.NotifyAFKResult:
		majsoul.IFReceive.NotifyAFKResult(notify)
	case *message.NotifyGameFinishRewardV2:
		majsoul.IFReceive.NotifyGameFinishRewardV2(notify)
	case *message.NotifyActivityRewardV2:
		majsoul.IFReceive.NotifyActivityRewardV2(notify)
	case *message.NotifyActivityPointV2:
		majsoul.IFReceive.NotifyActivityPointV2(notify)
	case *message.NotifyLeaderboardPointV2:
		majsoul.IFReceive.NotifyLeaderboardPointV2(notify)
	case *message.NotifyNewGame:
		majsoul.IFReceive.NotifyNewGame(notify)
	case *message.NotifyPlayerLoadGameReady:
		majsoul.IFReceive.NotifyPlayerLoadGameReady(notify)
	case *message.NotifyGameBroadcast:
		majsoul.IFReceive.NotifyGameBroadcast(notify)
	case *message.NotifyGameEndResult:
		majsoul.IFReceive.NotifyGameEndResult(notify)
	case *message.NotifyGameTerminate:
		majsoul.IFReceive.NotifyGameTerminate(notify)
	case *message.NotifyPlayerConnectionState:
		majsoul.IFReceive.NotifyPlayerConnectionState(notify)
	case *message.NotifyAccountLevelChange:
		majsoul.IFReceive.NotifyAccountLevelChange(notify)
	case *message.NotifyGameFinishReward:
		majsoul.IFReceive.NotifyGameFinishReward(notify)
	case *message.NotifyActivityReward:
		majsoul.IFReceive.NotifyActivityReward(notify)
	case *message.NotifyActivityPoint:
		majsoul.IFReceive.NotifyActivityPoint(notify)
	case *message.NotifyLeaderboardPoint:
		majsoul.IFReceive.NotifyLeaderboardPoint(notify)
	case *message.NotifyGamePause:
		majsoul.IFReceive.NotifyGamePause(notify)
	case *message.NotifyEndGameVote:
		majsoul.IFReceive.NotifyEndGameVote(notify)
	case *message.NotifyObserveData:
		majsoul.IFReceive.NotifyObserveData(notify)
	case *message.NotifyRoomPlayerReady_AccountReadyState:
		majsoul.IFReceive.NotifyRoomPlayerReady_AccountReadyState(notify)
	case *message.NotifyRoomPlayerDressing_AccountDressingState:
		majsoul.IFReceive.NotifyRoomPlayerDressing_AccountDressingState(notify)
	case *message.NotifyAnnouncementUpdate_AnnouncementUpdate:
		majsoul.IFReceive.NotifyAnnouncementUpdate_AnnouncementUpdate(notify)
	case *message.NotifyActivityUpdate_FeedActivityData:
		majsoul.IFReceive.NotifyActivityUpdate_FeedActivityData(notify)
	case *message.NotifyActivityUpdate_FeedActivityData_CountWithTimeData:
		majsoul.IFReceive.NotifyActivityUpdate_FeedActivityData_CountWithTimeData(notify)
	case *message.NotifyActivityUpdate_FeedActivityData_GiftBoxData:
		majsoul.IFReceive.NotifyActivityUpdate_FeedActivityData_GiftBoxData(notify)
	case *message.NotifyPayResult_ResourceModify:
		majsoul.IFReceive.NotifyPayResult_ResourceModify(notify)
	case *message.NotifyGameFinishRewardV2_LevelChange:
		majsoul.IFReceive.NotifyGameFinishRewardV2_LevelChange(notify)
	case *message.NotifyGameFinishRewardV2_MatchChest:
		majsoul.IFReceive.NotifyGameFinishRewardV2_MatchChest(notify)
	case *message.NotifyGameFinishRewardV2_MainCharacter:
		majsoul.IFReceive.NotifyGameFinishRewardV2_MainCharacter(notify)
	case *message.NotifyGameFinishRewardV2_CharacterGift:
		majsoul.IFReceive.NotifyGameFinishRewardV2_CharacterGift(notify)
	case *message.NotifyActivityRewardV2_ActivityReward:
		majsoul.IFReceive.NotifyActivityRewardV2_ActivityReward(notify)
	case *message.NotifyActivityPointV2_ActivityPoint:
		majsoul.IFReceive.NotifyActivityPointV2_ActivityPoint(notify)
	case *message.NotifyLeaderboardPointV2_LeaderboardPoint:
		majsoul.IFReceive.NotifyLeaderboardPointV2_LeaderboardPoint(notify)
	case *message.NotifyGameFinishReward_LevelChange:
		majsoul.IFReceive.NotifyGameFinishReward_LevelChange(notify)
	case *message.NotifyGameFinishReward_MatchChest:
		majsoul.IFReceive.NotifyGameFinishReward_MatchChest(notify)
	case *message.NotifyGameFinishReward_MainCharacter:
		majsoul.IFReceive.NotifyGameFinishReward_MainCharacter(notify)
	case *message.NotifyGameFinishReward_CharacterGift:
		majsoul.IFReceive.NotifyGameFinishReward_CharacterGift(notify)
	case *message.NotifyActivityReward_ActivityReward:
		majsoul.IFReceive.NotifyActivityReward_ActivityReward(notify)
	case *message.NotifyActivityPoint_ActivityPoint:
		majsoul.IFReceive.NotifyActivityPoint_ActivityPoint(notify)
	case *message.NotifyLeaderboardPoint_LeaderboardPoint:
		majsoul.IFReceive.NotifyLeaderboardPoint_LeaderboardPoint(notify)
	case *message.NotifyEndGameVote_VoteResult:
		majsoul.IFReceive.NotifyEndGameVote_VoteResult(notify)
	case *message.ActionPrototype:
		majsoul.IFReceive.ActionPrototype(notify)
	}
}
func (majsoul *Majsoul) NotifyCaptcha(notify *message.NotifyCaptcha) {}
func (majsoul *Majsoul) NotifyRoomGameStart(notify *message.NotifyRoomGameStart) {
	majsoul.game = NewClientConn(majsoul.Ctx, majsoul.Config.GameAddress)
	majsoul.FastTestClient = message.NewFastTestClient(majsoul.game)
	go majsoul.receiveGame()
	var err error
	majsoul.GameInfo, err = majsoul.AuthGame(majsoul.Ctx, &message.ReqAuthGame{
		AccountId: majsoul.Account.AccountId,
		Token:     notify.ConnectToken,
		GameUuid:  notify.GameUuid,
	})
	if err != nil {
		logger.Panic("NotifyRoomGameStart", zap.Error(err))
		return
	}
	fields := make([]zap.Field, 0)
	for _, player := range majsoul.GameInfo.Players {
		fields = append(fields, zap.Uint32(player.Nickname, player.AccountId))
	}
	logger.Debug("AuthGame", fields...)
	_, err = majsoul.CheckNetworkDelay(majsoul.Ctx, &message.ReqCommon{})
	if err != nil {
		logger.Panic("NotifyRoomGameStart", zap.Error(err))
		return
	}
	_, err = majsoul.EnterGame(majsoul.Ctx, &message.ReqCommon{})
	if err != nil {
		logger.Panic("NotifyRoomGameStart", zap.Error(err))
		return
	}
	logger.Debug("EnterGame")
}
func (majsoul *Majsoul) NotifyMatchGameStart(notify *message.NotifyMatchGameStart)         {}
func (majsoul *Majsoul) NotifyRoomPlayerReady(notify *message.NotifyRoomPlayerReady)       {}
func (majsoul *Majsoul) NotifyRoomPlayerDressing(notify *message.NotifyRoomPlayerDressing) {}
func (majsoul *Majsoul) NotifyRoomPlayerUpdate(notify *message.NotifyRoomPlayerUpdate)     {}
func (majsoul *Majsoul) NotifyRoomKickOut(notify *message.NotifyRoomKickOut)               {}
func (majsoul *Majsoul) NotifyFriendStateChange(notify *message.NotifyFriendStateChange)   {}
func (majsoul *Majsoul) NotifyFriendViewChange(notify *message.NotifyFriendViewChange)     {}
func (majsoul *Majsoul) NotifyFriendChange(notify *message.NotifyFriendChange)             {}
func (majsoul *Majsoul) NotifyNewFriendApply(notify *message.NotifyNewFriendApply)         {}
func (majsoul *Majsoul) NotifyClientMessage(notify *message.NotifyClientMessage)           {}
func (majsoul *Majsoul) NotifyAccountUpdate(notify *message.NotifyAccountUpdate)           {}
func (majsoul *Majsoul) NotifyAnotherLogin(notify *message.NotifyAnotherLogin)             {}
func (majsoul *Majsoul) NotifyAccountLogout(notify *message.NotifyAccountLogout)           {}
func (majsoul *Majsoul) NotifyAnnouncementUpdate(notify *message.NotifyAnnouncementUpdate) {}
func (majsoul *Majsoul) NotifyNewMail(notify *message.NotifyNewMail)                       {}
func (majsoul *Majsoul) NotifyDeleteMail(notify *message.NotifyDeleteMail)                 {}
func (majsoul *Majsoul) NotifyReviveCoinUpdate(notify *message.NotifyReviveCoinUpdate)     {}
func (majsoul *Majsoul) NotifyDailyTaskUpdate(notify *message.NotifyDailyTaskUpdate)       {}
func (majsoul *Majsoul) NotifyActivityTaskUpdate(notify *message.NotifyActivityTaskUpdate) {}
func (majsoul *Majsoul) NotifyActivityPeriodTaskUpdate(notify *message.NotifyActivityPeriodTaskUpdate) {
}
func (majsoul *Majsoul) NotifyAccountRandomTaskUpdate(notify *message.NotifyAccountRandomTaskUpdate) {
}
func (majsoul *Majsoul) NotifyActivitySegmentTaskUpdate(notify *message.NotifyActivitySegmentTaskUpdate) {
}
func (majsoul *Majsoul) NotifyActivityUpdate(notify *message.NotifyActivityUpdate) {}
func (majsoul *Majsoul) NotifyAccountChallengeTaskUpdate(notify *message.NotifyAccountChallengeTaskUpdate) {
}
func (majsoul *Majsoul) NotifyNewComment(notify *message.NotifyNewComment)           {}
func (majsoul *Majsoul) NotifyRollingNotice(notify *message.NotifyRollingNotice)     {}
func (majsoul *Majsoul) NotifyGiftSendRefresh(notify *message.NotifyGiftSendRefresh) {}
func (majsoul *Majsoul) NotifyShopUpdate(notify *message.NotifyShopUpdate)           {}
func (majsoul *Majsoul) NotifyVipLevelChange(notify *message.NotifyVipLevelChange)   {}
func (majsoul *Majsoul) NotifyServerSetting(notify *message.NotifyServerSetting)     {}
func (majsoul *Majsoul) NotifyPayResult(notify *message.NotifyPayResult)             {}
func (majsoul *Majsoul) NotifyCustomContestAccountMsg(notify *message.NotifyCustomContestAccountMsg) {
}
func (majsoul *Majsoul) NotifyCustomContestSystemMsg(notify *message.NotifyCustomContestSystemMsg) {}
func (majsoul *Majsoul) NotifyMatchTimeout(notify *message.NotifyMatchTimeout)                     {}
func (majsoul *Majsoul) NotifyCustomContestState(notify *message.NotifyCustomContestState)         {}
func (majsoul *Majsoul) NotifyActivityChange(notify *message.NotifyActivityChange)                 {}
func (majsoul *Majsoul) NotifyAFKResult(notify *message.NotifyAFKResult)                           {}
func (majsoul *Majsoul) NotifyGameFinishRewardV2(notify *message.NotifyGameFinishRewardV2)         {}
func (majsoul *Majsoul) NotifyActivityRewardV2(notify *message.NotifyActivityRewardV2)             {}
func (majsoul *Majsoul) NotifyActivityPointV2(notify *message.NotifyActivityPointV2)               {}
func (majsoul *Majsoul) NotifyLeaderboardPointV2(notify *message.NotifyLeaderboardPointV2)         {}
func (majsoul *Majsoul) NotifyNewGame(notify *message.NotifyNewGame)                               {}
func (majsoul *Majsoul) NotifyPlayerLoadGameReady(notify *message.NotifyPlayerLoadGameReady)       {}
func (majsoul *Majsoul) NotifyGameBroadcast(notify *message.NotifyGameBroadcast)                   {}
func (majsoul *Majsoul) NotifyGameEndResult(notify *message.NotifyGameEndResult)                   {}
func (majsoul *Majsoul) NotifyGameTerminate(notify *message.NotifyGameTerminate)                   {}
func (majsoul *Majsoul) NotifyPlayerConnectionState(notify *message.NotifyPlayerConnectionState)   {}
func (majsoul *Majsoul) NotifyAccountLevelChange(notify *message.NotifyAccountLevelChange)         {}
func (majsoul *Majsoul) NotifyGameFinishReward(notify *message.NotifyGameFinishReward)             {}
func (majsoul *Majsoul) NotifyActivityReward(notify *message.NotifyActivityReward)                 {}
func (majsoul *Majsoul) NotifyActivityPoint(notify *message.NotifyActivityPoint)                   {}
func (majsoul *Majsoul) NotifyLeaderboardPoint(notify *message.NotifyLeaderboardPoint)             {}
func (majsoul *Majsoul) NotifyGamePause(notify *message.NotifyGamePause)                           {}
func (majsoul *Majsoul) NotifyEndGameVote(notify *message.NotifyEndGameVote)                       {}
func (majsoul *Majsoul) NotifyObserveData(notify *message.NotifyObserveData)                       {}
func (majsoul *Majsoul) NotifyRoomPlayerReady_AccountReadyState(notify *message.NotifyRoomPlayerReady_AccountReadyState) {
}
func (majsoul *Majsoul) NotifyRoomPlayerDressing_AccountDressingState(notify *message.NotifyRoomPlayerDressing_AccountDressingState) {
}
func (majsoul *Majsoul) NotifyAnnouncementUpdate_AnnouncementUpdate(notify *message.NotifyAnnouncementUpdate_AnnouncementUpdate) {
}
func (majsoul *Majsoul) NotifyActivityUpdate_FeedActivityData(notify *message.NotifyActivityUpdate_FeedActivityData) {
}
func (majsoul *Majsoul) NotifyActivityUpdate_FeedActivityData_CountWithTimeData(notify *message.NotifyActivityUpdate_FeedActivityData_CountWithTimeData) {
}
func (majsoul *Majsoul) NotifyActivityUpdate_FeedActivityData_GiftBoxData(notify *message.NotifyActivityUpdate_FeedActivityData_GiftBoxData) {
}
func (majsoul *Majsoul) NotifyPayResult_ResourceModify(notify *message.NotifyPayResult_ResourceModify) {
}
func (majsoul *Majsoul) NotifyGameFinishRewardV2_LevelChange(notify *message.NotifyGameFinishRewardV2_LevelChange) {
}
func (majsoul *Majsoul) NotifyGameFinishRewardV2_MatchChest(notify *message.NotifyGameFinishRewardV2_MatchChest) {
}
func (majsoul *Majsoul) NotifyGameFinishRewardV2_MainCharacter(notify *message.NotifyGameFinishRewardV2_MainCharacter) {
}
func (majsoul *Majsoul) NotifyGameFinishRewardV2_CharacterGift(notify *message.NotifyGameFinishRewardV2_CharacterGift) {
}
func (majsoul *Majsoul) NotifyActivityRewardV2_ActivityReward(notify *message.NotifyActivityRewardV2_ActivityReward) {
}
func (majsoul *Majsoul) NotifyActivityPointV2_ActivityPoint(notify *message.NotifyActivityPointV2_ActivityPoint) {
}
func (majsoul *Majsoul) NotifyLeaderboardPointV2_LeaderboardPoint(notify *message.NotifyLeaderboardPointV2_LeaderboardPoint) {
}
func (majsoul *Majsoul) NotifyGameFinishReward_LevelChange(notify *message.NotifyGameFinishReward_LevelChange) {
}
func (majsoul *Majsoul) NotifyGameFinishReward_MatchChest(notify *message.NotifyGameFinishReward_MatchChest) {
}
func (majsoul *Majsoul) NotifyGameFinishReward_MainCharacter(notify *message.NotifyGameFinishReward_MainCharacter) {
}
func (majsoul *Majsoul) NotifyGameFinishReward_CharacterGift(notify *message.NotifyGameFinishReward_CharacterGift) {
}
func (majsoul *Majsoul) NotifyActivityReward_ActivityReward(notify *message.NotifyActivityReward_ActivityReward) {
}
func (majsoul *Majsoul) NotifyActivityPoint_ActivityPoint(notify *message.NotifyActivityPoint_ActivityPoint) {
}
func (majsoul *Majsoul) NotifyLeaderboardPoint_LeaderboardPoint(notify *message.NotifyLeaderboardPoint_LeaderboardPoint) {
}
func (majsoul *Majsoul) NotifyEndGameVote_VoteResult(notify *message.NotifyEndGameVote_VoteResult) {}
func (majsoul *Majsoul) ActionPrototype(notify *message.ActionPrototype)                           {}
