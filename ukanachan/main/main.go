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
	"utils/net"
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

	loginRes, err := m.Login(m.Ctx, &message.ReqLogin{
		Account:   "1601198895@qq.com",
		Password:  majsoul.Hash("miku39.."),
		Reconnect: false,
		Device: &message.ClientDeviceInfo{
			Platform:       "pc",
			Hardware:       "pc",
			Os:             "windows",
			OsVersion:      "win10",
			IsBrowser:      true,
			Software:       "Chrome",
			SalePlatform:   "web",
			HardwareVendor: "",
			ModelNumber:    "",
			ScreenWidth:    914,
			ScreenHeight:   1316,
		},
		RandomKey: "cfc35be-f519-4cbc-9765-c4c124cdc6a16",
		ClientVersion: &message.ClientVersionInfo{
			Resource: "0.10.105.w",
			Package:  "",
		},
		GenAccessToken:    true,
		CurrencyPlatforms: []uint32{2, 6, 8, 10, 11},
		// 电话1 邮箱0
		Type:                0,
		Version:             0,
		ClientVersionString: "web-0.10.105",
	})

	if err != nil {
		logger.Panic("Login failed", zap.Error(err))
	}

	if loginRes.Error != nil {
		logger.Panic("Login failed", zap.Reflect("Error", loginRes.Error))
	}

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
	logger.Debug("GetPaiPu")
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
	logger.Debug("GetPaiPu return", zap.Int("count", len(ret)))
	return ret
}

func GetRecord(m *majsoul.Majsoul, records paipu.Record) {
	logger.Debug("GetRecord")
	for _, record := range records {
		logger.Debug("FetchGameRecord", zap.String("uuid", record.UUID))
		fetchGameRecord, err := m.FetchGameRecord(m.Ctx, &message.ReqGameRecord{
			GameUuid:            record.UUID,
			ClientVersionString: "web-0.10.105",
		})
		if err != nil {
			logger.Info("GetRecord", zap.Error(err))
			time.Sleep(time.Millisecond)
			continue
		}
		filename := fmt.Sprintf("record/%s", record.UUID)
		err = SaveRecord(filename, fetchGameRecord)
		if err != nil {
			logger.Info("GetRecord", zap.Error(err))
			time.Sleep(time.Second * 3)
			continue
		}
		time.Sleep(time.Second * 10)
	}
}

func SaveRecord(filename string, record *message.ResGameRecord) error {
	var err error
	if len(record.Data) == 0 {
		record.Data, err = net.Get(record.DataUrl)
		if err != nil {
			return err
		}
	}
	l := len(record.Head.Result.Players)
	if l != 4 {
		return fmt.Errorf("记录%s(%d)不是4人麻将", filename, l)
	}
	logger.Debug("SaveRecord", zap.String("filename", filename))
	wrapper := new(message.Wrapper)
	wrapper.Data, err = proto.Marshal(record)
	if err != nil {
		logger.Debug("SaveRecord", zap.Error(err))
	}
	body, err := proto.Marshal(wrapper)
	if err != nil {
		logger.Debug("SaveRecord", zap.Error(err))
	}
	err = ioutil.WriteFile(filename, body, 0666)
	if err != nil {
		return err
	}
	return nil
}
