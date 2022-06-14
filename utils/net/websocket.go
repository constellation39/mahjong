package net

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
	"net/http"
	"sync"
	"utils/logger"
)

type WSClient struct {
	connAddr string // 连接地址
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
	conn, response, err := websocket.DefaultDialer.Dial(client.connAddr, client.header)
	if err != nil {
		logger.DPanic("WSClient.Connect failed", zap.String("addr", client.connAddr), zap.Error(err))
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.DPanic("WSClient.Connect close failed", zap.Error(err))
		}
	}(response.Body)

	client.conn = conn
	logger.Debug("WSClient.Connect success", zap.String("addr", client.connAddr))
	return nil
}

func (client *WSClient) Read() []byte {
	t, payload, err := client.conn.ReadMessage()
	if err != nil {
		logger.DPanic("WSClient.Read", zap.Error(err))
		return []byte{}
	}

	if t != websocket.BinaryMessage {
		logger.DPanic("WSClient.Read type not is websocket.BinaryMessage")
		return []byte{}
	}
	return payload
}

func (client *WSClient) Send(body []byte) {
	client.mu.Lock()
	defer client.mu.Unlock()

	err := client.conn.WriteMessage(websocket.BinaryMessage, body)
	if err != nil {
		logger.DPanic("WSClient.Send", zap.Reflect("body", body), zap.Error(err))
	}
}

func (client *WSClient) Close() {
	if client.conn == nil {
		logger.DPanic("WSClient.Close conn is nil")
		return
	}
	err := client.conn.Close()
	if err != nil {
		logger.DPanic("WSClient.Close conn close failed", zap.Error(err))
		return
	}
}
