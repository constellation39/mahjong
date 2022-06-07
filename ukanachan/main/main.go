package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"io/ioutil"
	"majsoul"
	"majsoul/message"
	"time"
	"ukanachan/paipu"
	_ "ukanachan/paipu"
	"utils/logger"
)

//func init() {
//	go UpdateToDay()
//	LoadConfig()
//	count, err := GetCount(StartTimestamp)
//	if err != nil {
//		logger.Error("GetCount", zap.Error(err))
//	}
//	logger.Debug("GetCount", zap.Int("count", count))
//	body, err := GetRecord(StartTimestamp, count, 12)
//	if err != nil {
//		logger.Error("GetCount", zap.Error(err))
//	}
//	logger.Debug("GetRecord", zap.Reflect("body", body))
//}

func main() {
	mCfg := majsoul.LoadConfig()
	m := majsoul.New(mCfg)

	version := m.GetVersion()
	if version.Version != "0.10.105.w" {
		logger.Info("liqi.json的版本为0.10.105.w,雀魂当前版本为", zap.String("Version", version.Version))
	}

	paipu.LoadConfig()

	go func() {
		t := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-t.C:
				_, err := m.Heatbeat(m.Ctx, &message.ReqHeatBeat{})
				if err != nil {
					logger.Info("Heatbeat", zap.Error(err))
					return
				}
			}
		}
	}()

	t := time.NewTicker(time.Millisecond)

	for {
		select {
		case <-t.C:
			records := GetPaiPu()
			GetRecord(m, records)
			paipu.StartTimestamp = paipu.NextDayTimestamp(paipu.StartTimestamp)
			paipu.SaveConfig()
			return
		}
	}

}

func GetPaiPu() paipu.Record {
	count, err := paipu.GetCount()
	if err != nil {
		logger.Error("GetPaiPu", zap.Error(err))
	}
	if count == 0 {
		return nil
	}
	logger.Debug("GetPaiPu", zap.String("StartTimestamp", paipu.StartTimestamp.String()), zap.Int("count", count))
	ret := make(paipu.Record, 0)
	for _, v := range paipu.Mode {
		body, err := paipu.GetRecord(count, v)
		if err != nil {
			logger.Error("GetPaiPu", zap.Error(err))
			continue
		}
		ret = append(ret, body...)
	}
	return ret
}

func GetRecord(m *majsoul.Majsoul, records paipu.Record) {
	for _, record := range records {
		fetchGameRecord, err := m.FetchGameRecord(m.Ctx, &message.ReqGameRecord{
			GameUuid:            record.UUID,
			ClientVersionString: "web-0.10.105",
		})
		if err != nil {
			logger.Info("GetRecord", zap.Error(err))
			time.Sleep(time.Millisecond)
			continue
		}
		body, err := proto.Marshal(fetchGameRecord)
		if err != nil {
			logger.Info("GetRecord", zap.Error(err))
			time.Sleep(time.Millisecond)
			continue
		}
		err = ioutil.WriteFile(fmt.Sprintf("record/%s", record.UUID), body, 0666)
		if err != nil {
			logger.Info("GetRecord", zap.Error(err))
		}
		time.Sleep(time.Second * 10)
	}
}
