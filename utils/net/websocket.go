package net

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"utils/logger"
)

type WSClient struct {
	connAddr string // 连接地址
	token    string // 身份认证token
	mu       *sync.Mutex
	header   http.Header
	conn     *websocket.Conn // websocket连接
}

func NewWSClient(addr string, header http.Header) *WSClient {
	return &WSClient{
		connAddr: addr,
		header:   header,
		mu:       &sync.Mutex{},
	}
}

func (client *WSClient) Connect() error {
	logger.Info("connecting to websocket server", zap.String("addr", client.connAddr))
	conn, response, err := websocket.DefaultDialer.Dial(client.connAddr, client.header)

	if err != nil {
		logger.Fatal("connect to websocket server failed", zap.Error(err))
		return err
	}

	defer response.Body.Close()
	client.conn = conn
	logger.Info("connect to websocket server success", zap.String("addr", client.connAddr), zap.String("response", response.Status))

	return nil
}

func (client *WSClient) Read() []byte {
	t, payload, err := client.conn.ReadMessage()

	if err != nil {
		logger.DPanic("read message failed", zap.Error(err))
		return []byte{}
	}

	if t != websocket.BinaryMessage {
		logger.DPanic("read message failed", zap.Error(err))
		return []byte{}
	}
	return payload
}

func (client *WSClient) Send(body []byte) {
	client.mu.Lock()
	defer client.mu.Unlock()

	err := client.conn.WriteMessage(websocket.BinaryMessage, body)
	if err != nil {
		logger.DPanic("send message failed", zap.Error(err))
	}
}

func (client *WSClient) Close() {
	if client.conn != nil {
		client.conn.Close()
	}
}
