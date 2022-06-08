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
	"io/ioutil"
	"majsoul/message"
	"math/rand"
	"net/http"
	"strings"
	"sync"
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

type ClientConn struct {
	Ctx context.Context
	*net.WSClient
	msgIndex uint8
	replys   sync.Map // 回复消息 map[uint8]*Reply
}

type Reply struct {
	out  proto.Message
	wait chan struct{}
	hook func(*message.Wrapper)
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
		switch msg[0] {
		case MsgTypeNotify:
			c.HandleNotify(msg)
		case MsgTypeResponse:
			c.HandleResponse(msg)
		default:
			logger.Info("message does not have any path", zap.ByteString("msg", msg))
		}
		select {
		case <-c.Ctx.Done():
			break receive
		default:
		}
	}
}

func (c *ClientConn) HandleNotify(msg []byte) {
	wrapper := new(message.Wrapper)
	err := proto.Unmarshal(msg[1:], wrapper)
	if err != nil {
		return
	}
}

func (c *ClientConn) HandleResponse(msg []byte) {
	key := (msg[2] << 8) + msg[1]
	v, ok := c.replys.Load(key)
	if !ok {
		logger.Error("Response not found", zap.Uint8("key", key))
		return
	}
	reply, ok := v.(*Reply)
	if !ok {
		logger.Error("rv not proto.Message", zap.Reflect("rv", reply))
		return
	}
	wrapper := new(message.Wrapper)
	err := proto.Unmarshal(msg[3:], wrapper)
	reply.hook(wrapper)
	if err != nil {
		logger.Error("proto.Unmarshal failed", zap.Error(err))
		return
	}
	err = proto.Unmarshal(wrapper.Data, reply.out)
	if err != nil {
		logger.Error("proto.Unmarshal failed", zap.Error(err))
		return
	}
	close(reply.wait)
}

func (c *ClientConn) Invoke(ctx context.Context, method string, in interface{}, out interface{}, opts ...grpc.CallOption) error {
	tokens := strings.Split(method, "/")
	api := strings.Join(tokens, ".")

	body, err := proto.Marshal(in.(proto.Message))
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
	buff.WriteByte(c.msgIndex - (c.msgIndex >> 8 << 8))
	buff.WriteByte(c.msgIndex >> 8)
	buff.Write(body)
	c.Send(buff.Bytes())
	reply := &Reply{
		out:  out.(proto.Message),
		wait: make(chan struct{}),
		hook: func(wrapper *message.Wrapper) {
			body, err := proto.Marshal(wrapper)
			if err != nil {
				logger.Debug("SaveRecord", zap.Error(err))
			}
			err = ioutil.WriteFile(fmt.Sprintf("./record/%s-%d", api, c.msgIndex), body, 0666)
			if err != nil {
				logger.Debug("SaveRecord", zap.Error(err))
			}
		},
	}
	if _, ok := c.replys.LoadOrStore(c.msgIndex, reply); ok {
		logger.DPanic("c.msgIndex exists", zap.Uint8("msgIndex", c.msgIndex))
	}
	defer c.replys.Delete(c.msgIndex)
	c.msgIndex++
	<-reply.wait
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
