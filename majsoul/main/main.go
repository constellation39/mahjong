package main

import (
	"encoding/json"
	"go.uber.org/zap"
	"majsoul"
	"majsoul/message"
	"utils/logger"
)

func main() {
	cfg := majsoul.LoadConfig()
	m := majsoul.New(cfg)

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
			Resource: "0.10.105.w",
			Package:  "",
		},
		GenAccessToken:    true,
		CurrencyPlatforms: []uint32{2, 6, 8, 10, 11},
		// 电话1 邮箱0
		Type:                0,
		Version:             0,
		ClientVersionString: "web-0.10.105",
	})
	if err != nil {
		return
	}
	logger.Debug("Login", zap.Reflect("Res", loginRes))

	fetchFriendList, err := m.FetchFriendList(m.Ctx, &message.ReqCommon{})
	if err != nil {
		return
	}
	logger.Debug("FetchFriendList", zap.Reflect("Res", fetchFriendList))

	loopNotify(m)

	select {
	case <-m.Ctx.Done():
	}
}

func loopNotify(m *majsoul.Majsoul) {

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

func NotifyClientMessage(notify *message.NotifyClientMessage) {
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
