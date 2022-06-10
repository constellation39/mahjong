package uakochan

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net"
	"time"
	"utils/logger"
)

const (
	E = "E"
	S = "S"
	W = "W"
	N = "N"
	//(東), S(南), W(西), N(北)
)

var port = ":11600"
var waits []*UAkochan

func init() {
	listen()
}

func listen() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		logger.Panic("listen error:", zap.Error(err))
	}
	logger.Info("uakochan listen on", zap.String("port", port))
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Debug("accept error:", zap.Error(err))
				continue
			}
			logger.Debug("accept a new connection from", zap.String("addr", conn.RemoteAddr().String()))
			if len(waits) == 0 {
				logger.Error("len(waits) == 0", zap.String("listener", conn.RemoteAddr().String()))
				conn.Close()
				break
			}
			wait := waits[0]
			wait.conn = conn
			go wait.receive()
			wait.wait <- struct{}{}
			waits = waits[1:]
		}
	}()
}

type UAkochan struct {
	conn net.Conn
	out  interface{}
	wait chan struct{}
}

func New() *UAkochan {
	uAkochan := &UAkochan{wait: make(chan struct{})}
	waits = append(waits, uAkochan)
	select {
	case <-time.After(time.Minute):
		logger.Panic("wait Akochan conn timeout")
	case <-uAkochan.wait:
	}
	uAkochan.Hello()
	uAkochan.StartGame([]string{"1", "2", "3", "4"})
	return uAkochan
}

func (uAkochan *UAkochan) Hello() string {
	res := new(Join)
	err := uAkochan.invoke(&Hello{
		Type:            "hello",
		Protocol:        "mjsonp",
		ProtocolVersion: 1,
	}, res)
	if err != nil {
		logger.Panic("hello error:", zap.Error(err))
	}
	logger.Debug("hello success", zap.Reflect("Name", res.Name))
	return res.Name
}

func (uAkochan *UAkochan) StartGame(names []string) {
	res := new(None)
	err := uAkochan.invoke(&StartGame{
		Type:  "start_kyoku",
		ID:    1,
		Names: names,
	}, res)
	if err != nil {
		logger.Error("start game error:", zap.Error(err))
		return
	}
}

func (uAkochan *UAkochan) StartKyoku(bakaze string, kyoku int, honba int, kyotaku int, oya int, dora_marker string, tehais [][]string) {
	res := new(None)
	err := uAkochan.invoke(&StartKyoku{
		Type:       "start_kyoku",
		Bakaze:     bakaze,
		Kyoku:      kyoku,
		Honba:      honba,
		Kyotaku:    kyotaku,
		Oya:        oya,
		DoraMarker: dora_marker,
		Tehais:     tehais,
	}, res)
	if err != nil {
		logger.Error("start kyoku error:", zap.Error(err))
		return
	}
}

func (uAkochan *UAkochan) Tsumo(actor int, pai string) {
	res := new(None)
}

func (uAkochan *UAkochan) invoke(in interface{}, out interface{}) error {
	uAkochan.send(in)
	uAkochan.out = out

	select {
	case <-time.After(time.Second * 3):
		return fmt.Errorf("invoke timeout")
	case <-uAkochan.wait:
	}
	return nil
}

func (uAkochan *UAkochan) send(in interface{}) {
	logger.Debug("send:", zap.Reflect("in", in))
	buff := new(bytes.Buffer)
	data, err := json.Marshal(in)
	if err != nil {
		logger.Error("json.Marshal error:", zap.Error(err))
	}
	buff.Write(data)
	buff.WriteByte('\n')
	_, err = uAkochan.conn.Write(buff.Bytes())
	if err != nil {
		logger.Debug("write error:", zap.Error(err))
		return
	}
}

func (uAkochan *UAkochan) receive() {
	defer uAkochan.conn.Close()
	reader := bufio.NewReader(uAkochan.conn)
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			logger.Debug("read error:", zap.Error(err))
			break
		}
		data = data[:len(data)-1]
		logger.Debug("recv:", zap.ByteString("raw", data))
		err = json.Unmarshal(data, uAkochan.out)
		if err != nil {
			logger.Error("json.Unmarshal error:", zap.Error(err))
		}
		uAkochan.wait <- struct{}{}
	}
}

func (uAkochan *UAkochan) close() {
	err := uAkochan.conn.Close()
	if err != nil {
		logger.Error("close conn", zap.Error(err))
	}
}

//func New

//
//type Conn struct {
//	net.Conn
//}
//
//func New(conn net.Conn) *Conn {
//	return &Conn{
//		Conn: conn,
//	}
//}
//
//func (conn *Conn) send(msg proto.Message) {
//	buff := new(bytes.Buffer)
//	data, _ := json.Marshal(msg)
//	buff.Write(data)
//	buff.WriteByte('\n')
//	conn.Write(buff.Bytes())
//}
//
//func (conn *Conn) handle() {
//	defer conn.close()
//	reader := bufio.NewReader(conn)
//	for {
//		data, err := reader.ReadBytes('\n')
//		if err != nil {
//			logger.Debug("read error:", zap.Error(err))
//			break
//		}
//		data = data[:len(data)-1]
//		logger.Debug("recv:", zap.ByteString("data", data))
//	}
//}
