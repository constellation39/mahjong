package main

import (
	"go.uber.org/zap"
	"majsoul"
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

}
