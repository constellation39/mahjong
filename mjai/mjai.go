package proxy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net"
	"utils/logger"

	"go.uber.org/zap"
)

type Proxy struct {
}

func New() *Proxy {
	return &Proxy{}
}

func Service() {
	conn := listen()
	for {
		conn, err := conn.Accept()
		if err != nil {
			logger.Debug("accept error:", zap.Error(err))
			continue
		}
		logger.Debug("accept a new connection from", zap.String("addr", conn.RemoteAddr().String()))
		go handle(conn)
	}
}

type Hello struct {
	Type            string `json:"type"`
	Protocol        string `json:"protocol"`
	ProtocolVersion int    `json:"protocol_version"`
}

var hello = &Hello{
	Type:            "hello",
	Protocol:        "mjsonp",
	ProtocolVersion: 1,
}

func handle(conn net.Conn) {
	defer conn.Close()
	buff := new(bytes.Buffer)
	data, _ := json.Marshal(hello)
	buff.Write(data)
	buff.WriteByte('\n')
	conn.Write(buff.Bytes())
	reader := bufio.NewReader(conn)
	for {
		data, err := reader.ReadBytes('\n')
		data = data[:len(data)-1]
		if err != nil {
			logger.Debug("read error:", zap.Error(err))
			break
		}
		logger.Debug("recv:", zap.ByteString("data", data))
	}
}

func listen() net.Listener {
	conn, err := net.Listen("tcp", ":11600")
	if err != nil {
		logger.Panic("listen error:", zap.Error(err))
	}
	logger.Info("listen on port 11600")
	return conn
}
