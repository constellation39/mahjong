package majsoul

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"majsoul/message"
	"math/rand"
	"strings"
	"utils/logger"
	"utils/net"
)

type Config struct {
	ServerAddress string `json:"serverAddress,omitempty"`
}

type ClientConn struct {
}

func (c ClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	tokens := strings.Split("/lq.Lobby/login", "/")
	api := strings.Join(tokens, ".")

}

func (c ClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	//TODO implement me
	panic("implement me")
}

var cConn = &ClientConn{}

type Majsoul struct {
	request     *net.Request
	lobbyClient message.LobbyClient
}

func New(c *Config) *Majsoul {
	return &Majsoul{
		request:     net.NewRequest(c.ServerAddress),
		lobbyClient: message.NewLobbyClient(cConn),
	}
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
