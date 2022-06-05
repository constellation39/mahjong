package net

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"logger"
	"net/http"
	"sync"
)

type WSClient struct {
	connAddr string // 连接地址
	token    string // 身份认证token
	mu       *sync.Mutex
	conn     *websocket.Conn // websocket连接
}

func NewWSClient(addr, token string) *WSClient {
	return &WSClient{
		connAddr: addr,
		token:    token,
		mu:       &sync.Mutex{},
	}
}

func (client *WSClient) Connect() error {
	header := http.Header{}
	header.Add("Authorization", "Bearer "+client.token)

	logger.Info("connecting to websocket server", zap.String("addr", client.connAddr))
	conn, response, err := websocket.DefaultDialer.Dial(client.connAddr, header)

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

	if t != websocket.TextMessage {
		return []byte{}
	}

	logger.Debug("receive message", zap.String("payload", string(payload)))
	return payload
}

func (client *WSClient) Send(data interface{}) {
	client.mu.Lock()
	defer client.mu.Unlock()

	logger.Debug("send message", zap.Reflect("data", data))

	body, err := json.Marshal(data)

	if err != nil {
		logger.DPanic("marshal message failed", zap.Error(err))
		return
	}

	err = client.conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		logger.DPanic("send message failed", zap.Error(err))
	}
}

func (client *WSClient) Close() {
	if client.conn != nil {
		client.conn.Close()
	}
}
