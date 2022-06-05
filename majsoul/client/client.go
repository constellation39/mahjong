package client

import "request"

type Config struct {
	ServerAddress string `json:"serverAddress,omitempty"`
}

// Client 实现了与服务器的交互
type Client struct {
	r *request.Request
}

func New(c *Config) {
	//r := request.New(c.ServerAddress)
}
