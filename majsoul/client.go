package majsoul

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"majsoul/message"
	"net/http"
	"strings"
	"sync"
	"utils/logger"
	"utils/net"
)

type ClientConn struct {
	ctx context.Context
	*net.WSClient
	msgIndex uint8
	replys   sync.Map // 回复消息 map[uint8]*Reply
	notify   chan proto.Message
}

type Reply struct {
	out  proto.Message
	wait chan struct{}
}

func NewClientConn(ctx context.Context, addr string) *ClientConn {
	header := http.Header{}
	header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44")
	cConn := &ClientConn{
		ctx:      ctx,
		WSClient: net.NewWSClient(addr, header),
		notify:   make(chan proto.Message, 32),
	}
	err := cConn.WSClient.Connect()
	if err != nil {
		logger.Fatal("connect to websocket server failed", zap.Error(err))
	}
	go cConn.loop()
	return cConn
}

func (c *ClientConn) loop() {
receive:
	for {
		msg := c.WSClient.Read()
		switch msg[0] {
		case MsgTypeNotify:
			c.handleNotify(msg)
		case MsgTypeResponse:
			c.handleResponse(msg)
		default:
			logger.Info("message does not have any path", zap.ByteString("msg", msg))
		}
		select {
		case <-c.ctx.Done():
			break receive
		default:
		}
	}
}

func (c *ClientConn) handleNotify(msg []byte) {
	wrapper := new(message.Wrapper)
	err := proto.Unmarshal(msg[1:], wrapper)
	if err != nil {
		return
	}
	logger.Debug("notify", zap.String("name", wrapper.Name))
}

func (c *ClientConn) handleResponse(msg []byte) {
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

func (c *ClientConn) Receive() proto.Message {
	return nil
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
