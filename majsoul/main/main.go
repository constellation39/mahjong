package main

import (
	"go.uber.org/zap"
	"majsoul"
	"majsoul/message"
	"time"
	"utils/config"
	"utils/logger"
)

func main() {
	cfg := new(majsoul.Config)
	err := config.Read("majsoul.json", cfg)
	if err != nil {
		logger.Panic("init client fail", zap.Error(err))
	}

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

	//loginRes, err := m.Login(m.Ctx, &message.ReqLogin{
	//	Account:   "1601198895@qq.com",
	//	Password:  majsoul.Hash("miku39.."),
	//	Reconnect: false,
	//	//win10(.2.Chro
	//	// me:.webP..X..*$f
	//	// cf915be-f559-41b
	//	// c-9765-c4e1248dc
	//	// 6a12...0.10.105.
	//	// w8.B......H.Z.we
	//	// b-0.10.105
	//	Device: &message.ClientDeviceInfo{
	//		Platform:       "pc",
	//		Hardware:       "pc",
	//		Os:             "windows",
	//		OsVersion:      "win10",
	//		IsBrowser:      true,
	//		Software:       "",
	//		SalePlatform:   "",
	//		HardwareVendor: "",
	//		ModelNumber:    "",
	//		ScreenWidth:    0,
	//		ScreenHeight:   0,
	//	},
	//	RandomKey:           "cfc35be-f519-4cbc-9765-c4c124cdc6a16",
	//	ClientVersion:       nil,
	//	GenAccessToken:      false,
	//	CurrencyPlatforms:   nil,
	//	Type:                0,
	//	Version:             0,
	//	ClientVersionString: version.Version,
	//})
	//if err != nil {
	//	return
	//}
	//logger.Debug("Login", zap.Reflect("Res", loginRes))
	time.Sleep(time.Hour)
}
