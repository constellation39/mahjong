package main

import (
	"go.uber.org/zap"
	"majsoul"
	"majsoul/message"
	"utils/logger"
)

func main() {
	cfg := majsoul.LoadConfig()
	m := majsoul.New(cfg)

	version := m.GetVersion()
	if version.Version != "0.10.105.w" {
		logger.Info("liqi.json的版本为0.10.105.w,雀魂当前版本为", zap.String("Version", version.Version))
	}

	heatbeatRes, err := m.Heatbeat(m.Ctx, &message.ReqHeatBeat{})
	if err != nil {
		return
	}
	logger.Debug("Heatbeat", zap.Reflect("Res", heatbeatRes))

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

	fetchGameRecordList, err := m.FetchGameRecordList(m.Ctx, &message.ReqGameRecordList{
		Start: 0,
		Count: 10,
		Type:  0,
	})
	if err != nil {
		return
	}
	logger.Debug("FetchGameRecordList", zap.Reflect("Res", fetchGameRecordList))
	//220606-86339438-f468-432a-9123-87f75dd129a9_a414418030
	fetchGameRecord, err := m.FetchGameRecord(m.Ctx, &message.ReqGameRecord{
		GameUuid:            "220606-86339438-f468-432a-9123-87f75dd129a9",
		ClientVersionString: "web-0.10.105",
	})
	if err != nil {
		return
	}
	logger.Debug("FetchGameRecord", zap.Reflect("Data", fetchGameRecord.Data))

}
