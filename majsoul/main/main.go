package main

import (
	"config"
	"go.uber.org/zap"
	"logger"
	"majsoul/client"
)

func main() {
	cfg := new(client.Config)
	err := config.Read("majsoul.json", cfg)
	if err != nil {
		logger.Panic("init client fail", zap.Error(err))
	}

	majsoul := client.New(cfg)

	version := majsoul.GetVersion()
	if version.Version != "0.10.105.w" {
		logger.Info("雀魂当前版本为", zap.String("Version", version.Version))
	}

}