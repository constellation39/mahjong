package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"io/ioutil"
	"majsoul"
	"majsoul/message"
	"time"
	"ukanachan/sapk"
	_ "ukanachan/sapk"
	"utils/logger"
	"utils/net"
)

func main() {
	mCfg := majsoul.LoadConfig()
	m := majsoul.New(mCfg)

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

	handle(m)
}

func handle(m *majsoul.Majsoul) {
	ticker := time.NewTicker(time.Second * 5)
	for records := range sapk.GetRecord(m.Ctx, ticker) {
		handleRecords(m, ticker, records)
	}
}

func handleRecords(m *majsoul.Majsoul, ticker *time.Ticker, records sapk.Record) {
	wrapper := new(message.Wrapper)
main:
	for _, record := range records {
		select {
		case <-ticker.C:
			fetchGameRecord, err := m.FetchGameRecord(m.Ctx, &message.ReqGameRecord{
				GameUuid:            record.UUID,
				ClientVersionString: "web-0.10.105",
			})
			if err != nil {
				logger.Info("FetchGameRecord", zap.Error(err))
				time.Sleep(time.Millisecond)
				continue
			}
			filename := fmt.Sprintf("record/%s", record.UUID)
			if len(fetchGameRecord.Data) == 0 {
				fetchGameRecord.Data, err = net.Get(fetchGameRecord.DataUrl)
				if err != nil {
					logger.Info("net.Get(fetchGameRecord.DataUrl)", zap.Error(err))
					continue
				}
			}
			l := len(fetchGameRecord.Head.Result.Players)
			if l != 4 {
				logger.Info("Players.len != 4", zap.String("filename", filename), zap.Int("l", l))
			}
			wrapper.Data, err = proto.Marshal(fetchGameRecord)
			if err != nil {
				logger.Info("Marshal", zap.Error(err))
			}
			body, err := proto.Marshal(wrapper)
			if err != nil {
				logger.Info("Marshal", zap.Error(err))
			}
			//buffer := new(bytes.Buffer)
			//buffer.Write(make([]byte, 3))
			//buffer.Write(body)
			err = ioutil.WriteFile(filename, body, 0666)
			if err != nil {
				logger.Info("WriteFile", zap.Error(err), zap.String("filename", filename), zap.ByteString("content", body))
			}
			logger.Info("WriteFile", zap.String("filename", filename))
		case <-m.Ctx.Done():
			break main
		}
	}
}
