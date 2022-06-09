package sapk

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"
	"utils/config"
	"utils/logger"
	"utils/net"
)

// startTimestamp 开始时间戳
var startTimestamp = time.Date(2019, time.November, 29, 0, 0, 0, 0, time.Local)
var toDay = time.Now()
var r = net.NewRequest("https://ak-data-1.sapk.ch")
var mode = []int{8, 9, 11, 12, 15, 16}

//"金东":  8,
//"金":   9,
//"玉东":  11,
//"玉":   12,
//"王座东": 15,
//"王座":  16,
//"三金":   22,
//"三玉":   24,
//"三王座":  26,
//"三金东":  21,
//"三玉东":  23,
//"三王座东": 25,

func init() {
	loadConfig()
}

type c struct {
	StartTimestamp int64 `json:"start_timestamp"`
	LastTimestamp  int64 `json:"last_timestamp"`
}

func loadConfig() {
	var pCfg c
	err := config.Read("sapk.json", &pCfg)
	if err != nil {
		logger.Error("loadConfig", zap.Error(err))
		time.Sleep(time.Second * 10)
		return
	}
	startTimestamp = time.UnixMilli(pCfg.StartTimestamp)
}

func saveConfig() {
	err := config.Write("sapk.json", &c{
		StartTimestamp: startTimestamp.UnixMilli(),
		LastTimestamp:  toDay.UnixMilli(),
	})
	if err != nil {
		logger.Panic("saveConfig", zap.Error(err))
	}
}

//func UpdateToDay() {
//	day := time.NewTicker(time.Second * 60 * 60 * 24)
//	toDay = <-day.C
//}

// toDayEndTimestamp 得到当前时间戳的终结时间戳
func toDayEndTimestamp(timestamp time.Time) time.Time {
	return startTimestamp.AddDate(0, 0, 1).Add(-time.Second)
}

// nextDayTimestamp 得到当前时间戳的下一天
func nextDayTimestamp(timestamp time.Time) time.Time {
	return timestamp.AddDate(0, 0, 1)
}

func nextDay() {
	startTimestamp = nextDayTimestamp(startTimestamp)
	saveConfig()
}

func GetRecord(ctx context.Context, ticker *time.Ticker) <-chan Record {
	out := make(chan Record)
	go loop(ctx, ticker, out)
	return out
}

func loop(ctx context.Context, ticker *time.Ticker, in chan<- Record) {
	for {
		select {
		case <-ctx.Done():
			break
		default:
			requestSapk(ctx, ticker, in)
			nextDay()
			time.Sleep(time.Second * 10)
		}
	}
}

func requestSapk(ctx context.Context, ticker *time.Ticker, in chan<- Record) {
	logger.Info("request Sapk", zap.Reflect("time", startTimestamp))
	cnt, err := getTotal()
	if err != nil {
		logger.Error("get total failed", zap.Time("startTimestamp", startTimestamp))
		return
	}
	if cnt == 0 {
		logger.Error("total is 0")
		return
	}
start:
	for _, mod := range mode {
		select {
		case <-ticker.C:
			record, err := getRecord(cnt, mod)
			if err != nil {
				logger.Info("get record failed", zap.Time("startTimestamp", startTimestamp), zap.Int("cnt", cnt), zap.Int("mod", mod))
				continue
			}
			in <- record
		case <-ctx.Done():
			break start
		}
	}
}

type total struct {
	Count int
}

// getTotal 得到当前时间戳的对局数统计
func getTotal() (int, error) {
	body, err := r.Get(fmt.Sprintf("api/count/%d", startTimestamp.UnixMilli()))
	if err != nil {
		return 0, err
	}
	var tal total
	err = json.Unmarshal(body, &tal)
	if err != nil {
		return 0, err
	}
	return tal.Count, nil
}

type Record []struct {
	ID        string     `json:"_id"`
	ModeID    int        `json:"modeId"`
	UUID      string     `json:"uuid"`
	StartTime int        `json:"startTime"`
	EndTime   int        `json:"endTime"`
	Players   []*Players `json:"Players"`
}
type Players struct {
	AccountID    int    `json:"accountId"`
	Nickname     string `json:"nickname"`
	Level        int    `json:"level"`
	Score        int    `json:"score"`
	GradingScore int    `json:"gradingScore"`
}

func getRecord(count, mode int) (Record, error) {
	tde := toDayEndTimestamp(startTimestamp)
	body, err := r.Get(fmt.Sprintf("api/v2/pl4/games/%d/%d?limit=%d&descending=true&mode=%d", tde.Unix(), startTimestamp.Unix(), count, mode))
	if err != nil {
		return nil, err
	}
	rcd := Record{}
	err = json.Unmarshal(body, &rcd)
	if err != nil {
		return nil, err
	}
	return rcd, nil
}
