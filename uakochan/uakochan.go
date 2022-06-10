package uakochan

import (
	"go.uber.org/zap"
	"net"
	"time"
	"utils/logger"
)

var waits []*UAkochan

func init() {
	listen()
}

type UAkochan struct {
	conn net.Conn
	wait chan struct{}
}

func listen() {
	listener, err := net.Listen("tcp", ":11600")
	if err != nil {
		logger.Panic("listen error:", zap.Error(err))
	}
	logger.Info("listen on port 11600")
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Debug("accept error:", zap.Error(err))
				continue
			}
			logger.Debug("accept a new connection from", zap.String("addr", conn.RemoteAddr().String()))
			if len(waits) == 0 {
				logger.DPanic("len(waits) == 0", zap.String("listener", conn.RemoteAddr().String()))
				err := conn.Close()
				if err != nil {
					logger.Error("listener.Close()", zap.Error(err))
					return
				}
				return
			}
			wait := waits[0]
			wait.conn = conn
			close(wait.wait)
			waits = waits[1:]
		}
	}()
}

func New() *UAkochan {
	uAkochan := &UAkochan{wait: make(chan struct{})}
	waits = append(waits, uAkochan)
	select {
	case <-time.After(time.Minute):
		logger.Panic("wait Akochan conn timeout")
	case <-uAkochan.wait:
	}
	return uAkochan
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
//func (conn *Conn) Send(msg proto.Message) {
//	buff := new(bytes.Buffer)
//	data, _ := json.Marshal(msg)
//	buff.Write(data)
//	buff.WriteByte('\n')
//	conn.Write(buff.Bytes())
//}
//
//func (conn *Conn) handle() {
//	defer conn.Close()
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
