package main

import (
	"config"
	"fmt"
	"go.uber.org/zap"
	"logger"
	"majsoul/client"
)

func main() {
	c := new(client.Config)
	err := config.Read("majsoul.json", c)

	if err != nil {
		logger.Panic("init client fail", zap.Error(err))
	}

	fmt.Printf("%+v", c)
}
