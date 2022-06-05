package main

import (
	"config"
	"fmt"
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
	fmt.Printf("%+v", version)
}
