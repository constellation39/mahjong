package client

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"logger"
	"math/rand"
	"request"
)

type Config struct {
	ServerAddress string `json:"serverAddress,omitempty"`
}

type Version struct {
	Version      string `json:"version"`
	ForceVersion string `json:"force_version"`
	Code         string `json:"code"`
}

// Client 实现了与服务器的交互
type Client struct {
	r *request.Request
}

func New(c *Config) *Client {
	return &Client{
		r: request.New(c.ServerAddress),
	}
}

func (c *Client) GetVersion() *Version {
	body, err := c.r.Get(fmt.Sprintf("1/version.json?randv=%d", rand.Intn(1000000000)))
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
