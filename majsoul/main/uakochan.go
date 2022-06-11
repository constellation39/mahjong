package main

import (
	"go.uber.org/zap"
	"majsoul"
	"majsoul/message"
	"uakochan"
	"utils/logger"
)

func (m *Majsoul) Tsumo(actor int, tile string) {
	reply := m.UAkochan.Tsumo(actor, uakochan.GetAkoChanTile(tile))

	if reply == nil {
		logger.Error("tsumo reply is nil")
		return
	}

	switch reply.(type) {
	case *uakochan.None:
		logger.Debug("tsumo reply is None", zap.Int("actor", actor), zap.String("tile", tile))
	case *uakochan.Dahai:
		logger.Debug("tsumo reply is Dahai", zap.Int("actor", actor), zap.String("tile", tile))
		in := reply.(*uakochan.Dahai)
		_, err := m.InputOperation(m.Ctx, &message.ReqSelfOperation{
			Type:      majsoul.DISCARD,
			Tile:      uakochan.GetSoulTile(in.Pai),
			Moqie:     in.Tsumogiri,
			Timeuse:   1,
			TileState: 0,
		})
		if err != nil {
			logger.Error("tsumo error:", zap.Error(err))
			return
		}
	case *uakochan.Kakan:
		logger.Debug("tsumo reply is Kakan", zap.Int("actor", actor), zap.String("tile", tile))
		in := reply.(*uakochan.Kakan)
		_, err := m.InputOperation(m.Ctx, &message.ReqSelfOperation{
			Type:            majsoul.KAKAN,
			Index:           0,
			Tile:            "",
			CancelOperation: false,
			Moqie:           false,
			Timeuse:         0,
			TileState:       0,
			ChangeTiles:     nil,
			TileStates:      nil,
			GapType:         0,
		})
		if err != nil {
			logger.Error("tsumo error:", zap.Error(err))
			return
		}
	case *uakochan.Ankan:
		logger.Debug("tsumo reply is Ankan", zap.Int("actor", actor), zap.String("tile", tile))
	case *uakochan.Reach:
		logger.Debug("tsumo reply is Reach", zap.Int("actor", actor), zap.String("tile", tile))
	case *uakochan.Hora:
		logger.Debug("tsumo reply is Hora", zap.Int("actor", actor), zap.String("tile", tile))
	default:
		logger.Error("tsumo reply is unknown", zap.Int("actor", actor), zap.String("tile", tile), zap.Any("reply", reply))
	}
}

func (m *Majsoul) Dahai(actor int, tile string, tsumogiri bool) {
	reply := m.UAkochan.Dahai(actor, uakochan.GetAkoChanTile(tile), tsumogiri)

	if reply == nil {
		logger.Error("dahai reply is nil")
		return
	}

	switch reply.(type) {
	case *uakochan.None:
		logger.Debug("dahai reply is None", zap.Int("actor", actor), zap.String("tile", tile))
	case *uakochan.Pon:
		logger.Debug("dahai reply is Pon", zap.Int("actor", actor), zap.String("tile", tile))
	case *uakochan.Chi:
		logger.Debug("dahai reply is Chi", zap.Int("actor", actor), zap.String("tile", tile))
	case *uakochan.Daiminkan:
		logger.Debug("dahai reply is Daiminkan", zap.Int("actor", actor), zap.String("tile", tile))
	case *uakochan.Hora:
		logger.Debug("dahai reply is Hora", zap.Int("actor", actor), zap.String("tile", tile))
	default:
		logger.Error("dahai reply is unknown", zap.Int("actor", actor), zap.String("tile", tile), zap.Any("reply", reply))
	}
}
