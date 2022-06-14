package main

import (
	"go.uber.org/zap"
	"majsoul"
	"majsoul/message"
	"uakochan"
	"utils/logger"
)

func (m *Majsoul) Tsumo(actor int, tile string, optionalOperationList *message.OptionalOperationList) {
	reply := m.UAkochan.Tsumo(actor, GetUAkochanTile(tile))

	if reply == nil {
		logger.Error("tsumo reply is nil")
		return
	}

	switch reply.(type) {
	case *uakochan.None:
		logger.Debug("tsumo reply is None", zap.Int("actor", actor), zap.String("tile", tile))
	case *uakochan.Dahai:
		logger.Debug("tsumo reply is Dahai", zap.Int("actor", actor), zap.String("tile", tile))
		option := GetOptionalOperation(majsoul.DISCARD, optionalOperationList)
		if option == nil {
			logger.Error("tsumo reply is Dahai, but option is nil", zap.Reflect("optionalOperationList", optionalOperationList))
			return
		}
		in := reply.(*uakochan.Dahai)
		_, err := m.InputOperation(m.Ctx, &message.ReqSelfOperation{
			Type:      majsoul.DISCARD,
			Tile:      GetSoulTile(in.Pai),
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
		option := GetOptionalOperation(majsoul.KAKAN, optionalOperationList)
		if option == nil {
			logger.Error("tsumo reply is Kakan, but option is nil", zap.Reflect("optionalOperationList", optionalOperationList))
			return
		}
		in := reply.(*uakochan.Kakan)
		comb := GetSoulComb(in.Consumed)
		index := GetOperationIndex(comb, option)
		_, err := m.InputOperation(m.Ctx, &message.ReqSelfOperation{
			Type:    majsoul.KAKAN,
			Index:   index,
			Timeuse: 1,
		})
		if err != nil {
			logger.Error("tsumo error:", zap.Error(err))
			return
		}
	case *uakochan.Ankan:
		logger.Debug("tsumo reply is Ankan", zap.Int("actor", actor), zap.String("tile", tile))
		option := GetOptionalOperation(majsoul.ANKAN, optionalOperationList)
		if option == nil {
			logger.Error("tsumo reply is Ankan, but option is nil")
			return
		}
		in := reply.(*uakochan.Ankan)
		comb := GetSoulComb(in.Consumed)
		index := GetOperationIndex(comb, option)
		_, err := m.InputOperation(m.Ctx, &message.ReqSelfOperation{
			Type:    majsoul.ANKAN,
			Index:   index,
			Timeuse: 1,
		})
		if err != nil {
			logger.Error("tsumo error:", zap.Error(err))
			return
		}
	case *uakochan.Reach:
		logger.Debug("tsumo reply is Reach", zap.Int("actor", actor), zap.String("tile", tile))
		option := GetOptionalOperation(majsoul.RIICHI, optionalOperationList)
		if option == nil {
			logger.Error("tsumo reply is Reach, but option is nil")
			return
		}
		in := reply.(*uakochan.Reach)
		r := m.UAkochan.Reach(in.Actor)
		i := r.(*uakochan.Dahai)
		_, err := m.InputOperation(m.Ctx, &message.ReqSelfOperation{
			Type:    majsoul.RIICHI,
			Tile:    GetSoulTile(i.Pai),
			Moqie:   i.Tsumogiri,
			Timeuse: 1,
		})
		logger.Debug("Reach", zap.Reflect("ReqSelfOperation", &message.ReqSelfOperation{
			Type:    majsoul.RIICHI,
			Tile:    GetSoulTile(i.Pai),
			Moqie:   i.Tsumogiri,
			Timeuse: 1,
		}))
		if err != nil {
			logger.Error("tsumo error:", zap.Error(err))
			return
		}
	case *uakochan.Hora:
		logger.Debug("tsumo reply is Hora", zap.Int("actor", actor), zap.String("tile", tile))
		option := GetOptionalOperation(majsoul.TSUMO, optionalOperationList)
		if option == nil {
			logger.Error("tsumo reply is Hora, but option is nil")
			return
		}
		_, err := m.InputOperation(m.Ctx, &message.ReqSelfOperation{
			Type:  majsoul.TSUMO,
			Index: 0,
		})
		if err != nil {
			logger.Error("tsumo error:", zap.Error(err))
			return
		}
	default:
		logger.Error("tsumo reply is unknown", zap.Int("actor", actor), zap.String("tile", tile), zap.Any("reply", reply))
	}
}

func (m *Majsoul) Dahai(actor int, tile string, tsumogiri bool, optionalOperationList *message.OptionalOperationList) {
	reply := m.UAkochan.Dahai(actor, GetUAkochanTile(tile), tsumogiri)

	if reply == nil {
		logger.Error("dahai reply is nil")
		return
	}

	switch reply.(type) {
	case *uakochan.None:
		logger.Debug("dahai reply is None", zap.Int("actor", actor), zap.String("tile", tile))
		if optionalOperationList != nil && len(optionalOperationList.OperationList) > 0 {
			_, err := m.InputChiPengGang(m.Ctx, &message.ReqChiPengGang{
				CancelOperation: true,
				Timeuse:         1,
			})
			if err != nil {
				logger.Error("none error:", zap.Error(err))
				return
			}
		}
	case *uakochan.Pon:
		logger.Debug("dahai reply is Pon", zap.Int("actor", actor), zap.String("tile", tile))
		option := GetOptionalOperation(majsoul.PON, optionalOperationList)
		if option == nil {
			logger.Error("dahai reply is Pon, but option is nil")
			return
		}
		in := reply.(*uakochan.Pon)
		comb := GetSoulComb(in.Consumed)
		index := GetOperationIndex(comb, option)
		_, err := m.InputChiPengGang(m.Ctx, &message.ReqChiPengGang{
			Type:    majsoul.PON,
			Index:   index,
			Timeuse: 1,
		})
		if err != nil {
			logger.Error("dahai error:", zap.Error(err))
			return
		}
	case *uakochan.Chi:
		logger.Debug("dahai reply is Chi", zap.Int("actor", actor), zap.String("tile", tile))
		option := GetOptionalOperation(majsoul.CHI, optionalOperationList)
		if option == nil {
			logger.Error("dahai reply is Chi, but option is nil")
			return
		}
		in := reply.(*uakochan.Chi)
		comb := GetSoulComb(in.Consumed)
		index := GetOperationIndex(comb, option)
		_, err := m.InputChiPengGang(m.Ctx, &message.ReqChiPengGang{
			Type:    majsoul.CHI,
			Index:   index,
			Timeuse: 1,
		})
		if err != nil {
			logger.Error("dahai error:", zap.Error(err))
			return
		}
	case *uakochan.Daiminkan:
		logger.Debug("dahai reply is Daiminkan", zap.Int("actor", actor), zap.String("tile", tile))
		option := GetOptionalOperation(majsoul.MINKAN, optionalOperationList)
		if option == nil {
			logger.Error("dahai reply is Daiminkan, but option is nil")
			return
		}
		in := reply.(*uakochan.Daiminkan)
		comb := GetSoulComb(in.Consumed)
		index := GetOperationIndex(comb, option)
		_, err := m.InputChiPengGang(m.Ctx, &message.ReqChiPengGang{
			Type:    majsoul.MINKAN,
			Index:   index,
			Timeuse: 1,
		})
		if err != nil {
			logger.Error("dahai error:", zap.Error(err))
			return
		}
	case *uakochan.Hora:
		logger.Debug("dahai reply is Hora", zap.Int("actor", actor), zap.String("tile", tile))
		option := GetOptionalOperation(majsoul.RON, optionalOperationList)
		if option == nil {
			logger.Error("dahai reply is Hora, but option is nil")
			return
		}
		_, err := m.InputOperation(m.Ctx, &message.ReqSelfOperation{
			Type:    majsoul.RON,
			Timeuse: 1,
		})
		if err != nil {
			logger.Error("dahai error:", zap.Error(err))
			return
		}
	default:
		logger.Error("dahai reply is unknown", zap.Int("actor", actor), zap.String("tile", tile), zap.Any("reply", reply))
	}
}

func GetOptionalOperation(option int, optionalOperationList *message.OptionalOperationList) (ret *message.OptionalOperation) {
	if optionalOperationList == nil {
		return nil
	}
	for _, operation := range optionalOperationList.OperationList {
		if operation.Type == uint32(option) {
			ret = operation
		}
	}
	return
}

func GetOperationIndex(need string, option *message.OptionalOperation) (ret uint32) {
	for i, c := range option.Combination {
		if c == need {
			ret = uint32(i)
			break
		}
	}
	return
}
