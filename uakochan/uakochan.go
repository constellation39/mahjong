package uakochan

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net"
	"sync"
	"utils/logger"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

var Conns sync.Map // map[string]*Conn

func init() {
	conn, err := net.Listen("tcp", ":11600")
	if err != nil {
		logger.Panic("listen error:", zap.Error(err))
	}
	logger.Info("listen on port 11600")
	go listen(conn)
}

func listen(conn net.Listener) {
	for {
		conn, err := conn.Accept()
		if err != nil {
			logger.Debug("accept error:", zap.Error(err))
			continue
		}
		logger.Debug("accept a new connection from", zap.String("addr", conn.RemoteAddr().String()))
		New(conn).handle()
	}
}

type Conn struct {
	net.Conn
}

func New(conn net.Conn) *Conn {
	return &Conn{
		Conn: conn,
	}
}

func Service() {

}

func (conn *Conn) Send(msg proto.Message) {
	buff := new(bytes.Buffer)
	data, _ := json.Marshal(msg)
	buff.Write(data)
	buff.WriteByte('\n')
	conn.Write(buff.Bytes())
}

func (conn *Conn) handle() {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			logger.Debug("read error:", zap.Error(err))
			break
		}
		data = data[:len(data)-1]
		logger.Debug("recv:", zap.ByteString("data", data))
	}
}
