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

var port = ":11600"
var waits []*UAkochan

func init() {
	listen()
}

func listen() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		logger.Panic("uakochan.listen", zap.Error(err))
	}
	logger.Info("uakochan.listen", zap.String("port", port))
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Debug("uakochan.listen Accept", zap.Error(err))
				continue
			}
			logger.Debug("uakochan.listen Accept", zap.String("addr", conn.RemoteAddr().String()))
			if len(waits) == 0 {
				logger.Error("uakochan.listen no wait", zap.String("listener", conn.RemoteAddr().String()))
				err := conn.Close()
				if err != nil {
					return
				}
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
	now := time.Now()
	uAkochan := &UAkochan{wait: make(chan struct{})}
	waits = append(waits, uAkochan)
	logger.Info("uakochan.New Wait")
	select {
	case <-time.After(time.Minute * 60):
		logger.Panic("uakochan.New timeout")
	case <-uAkochan.wait:
	}
	uAkochan.Hello()
	logger.Info("uakochan.New", zap.String("time", time.Since(now).String()))
	return uAkochan
}

func (uAkochan *UAkochan) Hello() string {
	err := uAkochan.invoke(&Hello{
		Type:            "hello",
		Protocol:        "mjsonp",
		ProtocolVersion: 1,
	})
	if err != nil {
		logger.Panic("UAkochan.Hello", zap.Error(err))
	}
	join := uAkochan.out.(*Join)
	logger.Debug("UAkochan.Hello", zap.Reflect("Name", join.Name))
	return join.Name
}

func (uAkochan *UAkochan) StartGame(id int, names []string) {
	err := uAkochan.invoke(&StartGame{
		Type:  "start_game",
		ID:    id,
		Names: names,
	})
	if err != nil {
		logger.Error("UAkochan.StartGame", zap.Error(err))
		return
	}
	logger.Debug("UAkochan.StartGame")
}

func (uAkochan *UAkochan) StartKyoku(bakaze string, kyoku uint32, honba uint32, kyotaku uint32, oya uint32, dora_marker string, tehais [][]string) {
	err := uAkochan.invoke(&StartKyoku{
		Type:       "start_kyoku",
		Bakaze:     bakaze,
		Kyoku:      kyoku,
		Honba:      honba,
		Kyotaku:    kyotaku,
		Oya:        oya,
		DoraMarker: dora_marker,
		Tehais:     tehais,
	})
	if err != nil {
		logger.Error("UAkochan.StartKyoku", zap.Error(err))
		return
	}
	logger.Debug("StartKyoku")
}

func (uAkochan *UAkochan) Tsumo(actor int, pai string) interface{} {
	err := uAkochan.invoke(&Tsumo{
		Type:  "tsumo",
		Actor: actor,
		Pai:   pai,
	})
	if err != nil {
		logger.Error("UAkochan.Tsumo", zap.Error(err))
		return nil
	}
	logger.Debug("UAkochan.Tsumo")
	return uAkochan.out
}

func (uAkochan *UAkochan) Dahai(actor int, pai string, tsumogiri bool) interface{} {
	err := uAkochan.invoke(&Dahai{
		Type:      "dahai",
		Actor:     actor,
		Pai:       pai,
		Tsumogiri: tsumogiri,
	})
	if err != nil {
		logger.Error("UAkochan.Dahai", zap.Error(err))
	}
	logger.Debug("UAkochan.Dahai")
	return uAkochan.out
}

func (uAkochan *UAkochan) Pon(actor, target int, pai string, consumed []string) interface{} {
	err := uAkochan.invoke(&Pon{
		Type:     "pon",
		Target:   target,
		Actor:    actor,
		Pai:      pai,
		Consumed: consumed,
	})
	if err != nil {
		logger.Error("UAkochan.Pon", zap.Error(err))
	}
	logger.Debug("UAkochan.Pon")
	return uAkochan.out
}

func (uAkochan *UAkochan) Chi(actor, target int, pai string, consumed []string) interface{} {
	err := uAkochan.invoke(&Chi{
		Type:     "chi",
		Target:   target,
		Actor:    actor,
		Pai:      pai,
		Consumed: consumed,
	})
	if err != nil {
		logger.Error("UAkochan.Chi", zap.Error(err))
	}
	logger.Debug("UAkochan.Chi")
	return uAkochan.out
}

func (uAkochan *UAkochan) Kakan(actor int, pai string, consumed []string) {
	err := uAkochan.invoke(&Kakan{
		Type:     "kakan",
		Actor:    actor,
		Pai:      pai,
		Consumed: consumed,
	})
	if err != nil {
		logger.Error("UAkochan.Kakan", zap.Error(err))
	}
	logger.Debug("UAkochan.Kakan")
}

func (uAkochan *UAkochan) Daiminkan(actor int, target int, pai string, consumed []string) {
	err := uAkochan.invoke(&Daiminkan{
		Type:     "daiminkan",
		Actor:    actor,
		Target:   target,
		Pai:      pai,
		Consumed: consumed,
	})
	if err != nil {
		logger.Error("UAkochan.Daiminkan", zap.Error(err))
	}
	logger.Debug("UAkochan.Daiminkan")
}

func (uAkochan *UAkochan) Ankan(actor int, consumed []string) {
	err := uAkochan.invoke(&Ankan{
		Type:     "ankan",
		Actor:    actor,
		Consumed: consumed,
	})
	if err != nil {
		logger.Error("UAkochan.Ankan", zap.Error(err))
	}
	logger.Debug("UAkochan.Ankan")
}

func (uAkochan *UAkochan) Reach(actor int) interface{} {
	err := uAkochan.invoke(&Reach{
		Type:  "reach",
		Actor: actor,
	})
	if err != nil {
		logger.Error("UAkochan.Reach", zap.Error(err))
	}
	logger.Debug("UAkochan.Reach")
	return uAkochan.out
}

func (uAkochan *UAkochan) ReachAccepted(actor int, deltas []int, scores []int) {
	err := uAkochan.invoke(&ReachAccepted{
		Type:   "reach_accepted",
		Actor:  actor,
		Deltas: deltas,
		Scores: scores,
	})
	if err != nil {
		logger.Error("UAkochan.ReachAccepted", zap.Error(err))
	}
	logger.Debug("UAkochan.ReachAccepted")
}

func (uAkochan *UAkochan) ReachAccepted_(actor int) {
	err := uAkochan.invoke(&ReachAccepted{
		Type:  "reach_accepted",
		Actor: actor,
	})
	if err != nil {
		logger.Error("UAkochan.ReachAccepted_", zap.Error(err))
	}
	logger.Debug("UAkochan.ReachAccepted_")
}

func (uAkochan *UAkochan) Hora(actor, target int, pai string, uradoraMarkers, horaTehais []string, yakus [][]interface{}, fu, fan, horaPoints int, deltas []int, scores []int) {
	err := uAkochan.invoke(&Hora{
		Type:           "hora",
		Actor:          actor,
		Target:         target,
		Pai:            "",
		UradoraMarkers: uradoraMarkers,
		HoraTehais:     horaTehais,
		Yakus:          yakus,
		Fu:             fu,
		Fan:            fan,
		HoraPoints:     horaPoints,
		Deltas:         deltas,
		Scores:         scores,
	})
	if err != nil {
		logger.Error("UAkochan.Hora", zap.Error(err))
	}
	logger.Debug("UAkochan.Hora")
}

func (uAkochan *UAkochan) Hora_(actor, target int) {
	err := uAkochan.invoke(&Hora{
		Type:   "hora",
		Actor:  actor,
		Target: target,
	})
	if err != nil {
		logger.Error("UAkochan.Hora", zap.Error(err))
	}
	logger.Debug("UAkochan.Hora")
}

func (uAkochan *UAkochan) EndKyoku() {
	err := uAkochan.invoke(&EndKyoku{Type: "end_kyoku"})
	if err != nil {
		logger.Error("UAkochan.EndKyoku", zap.Error(err))
		return
	}
	logger.Debug("UAkochan.EndKyoku")
}

func (uAkochan *UAkochan) Ryukyoku(reason string, tehais [][]string, tenpais []bool, deltas, scores []int) {
	err := uAkochan.invoke(&Ryukyoku{
		Type:    "ryukyoku",
		Reason:  reason,
		Tehais:  tehais,
		Tenpais: tenpais,
		Deltas:  deltas,
		Scores:  scores,
	})
	if err != nil {
		logger.Error("UAkochan.Ryukyoku", zap.Error(err))
	}
	logger.Debug("UAkochan.Ryukyoku")
}

func (uAkochan *UAkochan) Ryukyoku_() {
	err := uAkochan.invoke(&Ryukyoku{
		Type: "ryukyoku",
	})
	if err != nil {
		logger.Error("UAkochan.Ryukyoku_", zap.Error(err))
	}
	logger.Debug("UAkochan.Ryukyoku_")
}

//func (uAkochan *UAkochan) invoke(in interface{}, out interface{}) {
func (uAkochan *UAkochan) invoke(in interface{}) error {
	uAkochan.send(in)
	now := time.Now()
	defer func() {
		logger.Debug("UAkochan.invoke", zap.String("time", time.Since(now).String()))
	}()
	select {
	case <-time.After(time.Second * 10):
		return fmt.Errorf("invoke timeout")
	case <-uAkochan.wait:
	}
	return nil
}

func (uAkochan *UAkochan) send(in interface{}) {
	logger.Debug("UAkochan.send", zap.Reflect("in", in))
	buff := new(bytes.Buffer)
	data, err := json.Marshal(in)
	if err != nil {
		logger.Error("json.Marshal", zap.Error(err))
	}
	buff.Write(data)
	buff.WriteByte('\n')
	_, err = uAkochan.conn.Write(buff.Bytes())
	if err != nil {
		logger.Debug("write", zap.Error(err))
		return
	}
}

func (uAkochan *UAkochan) receive() {
	defer uAkochan.conn.Close()
	reader := bufio.NewReader(uAkochan.conn)
	t := new(Type)
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			logger.Debug("UAkochan.receive", zap.Error(err))
			break
		}
		data = data[:len(data)-1]
		logger.Debug("UAkochan.receive", zap.ByteString("raw", data))
		err = json.Unmarshal(data, t)
		if err != nil {
			logger.Error("UAkochan.receive", zap.Error(err))
		}
		v := GetMessageType(t.Type)
		err = json.Unmarshal(data, v)
		if err != nil {
			logger.Error("UAkochan.receive", zap.Error(err))
		}
		uAkochan.out = v
		uAkochan.wait <- struct{}{}
	}
}

func (uAkochan *UAkochan) close() {
	err := uAkochan.conn.Close()
	if err != nil {
		logger.Error("UAkochan.Close conn", zap.Error(err))
	}
}
