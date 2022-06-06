package majsoul

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"majsoul/message"
	"math/rand"
	"net/http"
	"strings"
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

type ClientConn struct {
	Ctx context.Context
	*net.WSClient
	msgIndex uint8
}

func NewClientConn(ctx context.Context, addr string) *ClientConn {
	header := http.Header{}
	header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44")
	cConn := &ClientConn{
		Ctx:      ctx,
		WSClient: net.NewWSClient(addr, header),
	}
	return cConn
}

func (c *ClientConn) Start() {
	err := c.WSClient.Connect()
	if err != nil {
		logger.Fatal("connect to websocket server failed", zap.Error(err))
	}
	go c.Loop()
}

func (c *ClientConn) Loop() {
receive:
	for {
		msg := c.WSClient.Read()
		fmt.Printf("%s", msg)
		//bot.handle(msg)
		select {
		case <-c.Ctx.Done():
			break receive
		default:
		}
	}
}

func (c *ClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	tokens := strings.Split(method, "/")
	api := strings.Join(tokens, ".")
	_ = api

	body, err := proto.Marshal(args.(proto.Message))
	if err != nil {
		logger.DPanic("marshal message failed", zap.Error(err))
		return fmt.Errorf("marshal message failed")
	}

	wrapper := &message.Wrapper{
		Name: api,
		Data: body,
	}

	body, err = proto.Marshal(wrapper)
	if err != nil {
		logger.DPanic("marshal message failed", zap.Error(err))
		return fmt.Errorf("marshal message failed")
	}
	buff := new(bytes.Buffer)
	c.msgIndex %= 255
	buff.WriteByte(MsgTypeRequest)
	buff.WriteByte(c.msgIndex - (c.msgIndex >> 7 << 7))
	buff.WriteByte(c.msgIndex >> 7)
	//buff.Write([]byte{2, 0, 0, 10, 18, 46, 108, 113, 46, 76, 111, 98, 98, 121, 46, 104, 101, 97, 116, 98, 101, 97, 116, 18, 0})
	buff.Write(body)
	c.Send(buff.Bytes())
	c.msgIndex++
	//c.Send(wrapper)
	return nil
}

func (c *ClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	//TODO implement me
	panic("implement me")
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
	cConn.Start()
	return &Majsoul{
		Ctx:         ctx,
		Cancel:      cancel,
		request:     net.NewRequest(c.ServerAddress),
		LobbyClient: message.NewLobbyClient(cConn),
	}
}

func Hash(data string) string {
	hash := hmac.New(sha256.New, []byte("lailai"))
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
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
