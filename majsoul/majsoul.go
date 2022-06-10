package majsoul

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

type Config struct {
	ServerAddress  string `json:"serverAddress"`
	GatewayAddress string `json:"gatewayAddress"`
}

type Majsoul struct {
	Ctx    context.Context
	Cancel context.CancelFunc
	message.LobbyClient
	request *net.Request
}

func New(c *Config) *Majsoul {
	ctx, cancel := context.WithCancel(context.Background())
	cConn := NewClientConn(ctx, c.GatewayAddress)
	m := &Majsoul{
		Ctx:         ctx,
		Cancel:      cancel,
		request:     net.NewRequest(c.ServerAddress),
		LobbyClient: message.NewLobbyClient(cConn),
	}
	m.check()
	go m.heatbeat()
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

func (majsoul *Majsoul) GetVersion() *Version {
	body, err := majsoul.request.Get(fmt.Sprintf("1/version.json?randv=%d", rand.Intn(1000000000)))
	if err != nil {
		logger.Error("GetVersion", zap.Error(err))
	}
	version := new(Version)
	err = json.Unmarshal(body, version)
	if err != nil {
		logger.Error("GetVersion", zap.Error(err))
	}
	return version
}

func (majsoul *Majsoul) check() {
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
