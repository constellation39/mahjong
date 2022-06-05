package paipu

import (
	"config"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"logger"
	"time"
)

// 设置一个开始时间戳
var StartTimestamp = time.Date(2019, time.November, 29, 0, 0, 0, 0, time.Local)
var ToDay = time.Now()
var r = request.New("https://ak-data-1.sapk.ch")
var Mode = map[string]int{
	"王座":   16,
	"玉":    12,
	"金":    9,
	"王座东":  15,
	"玉东":   11,
	"金东":   8,
	"三金":   22,
	"三玉":   24,
	"三王座":  26,
	"三金东":  21,
	"三玉东":  23,
	"三王座东": 25,
}

type C struct {
	StartTimestamp int64 `json:"start_timestamp"`
	LastTimestamp  int64 `json:"last_timestamp"`
}

func init() {
	go UpdateToDay()
	LoadConfig()
	count, err := GetCount(StartTimestamp)
	if err != nil {
		logger.Error("GetCount", zap.Error(err))
	}
	logger.Debug("GetCount", zap.Int("count", count))
	body, err := GetRecord(StartTimestamp, count, 12)
	if err != nil {
		logger.Error("GetCount", zap.Error(err))
	}
	logger.Debug("GetRecord", zap.Reflect("body", body))
}

func LoadConfig() {
	var c C
	err := config.Read("config/config.json", &c)
	if err != nil {
		logger.Info("LoadConfig", zap.Error(err))
		return
	}
	StartTimestamp = time.Unix(c.StartTimestamp, 0)
}

func SaveConfig() {
	err := config.Write("config/config.json", &C{
		StartTimestamp: StartTimestamp.Unix(),
		LastTimestamp:  ToDay.Unix(),
	})
	if err != nil {
		logger.Error("SaveConfig", zap.Error(err))
	}
}

func UpdateToDay() {
	day := time.NewTicker(time.Second * 60 * 60 * 24)
	ToDay = <-day.C
}

// 得到当前时间戳的终结时间戳
func ToDayEndTimestamp(timestamp time.Time) time.Time {
	return StartTimestamp.AddDate(0, 0, 1).Add(-time.Second)
}

// 得到当前时间戳的下一天
func NextDayTimestamp(timestamp time.Time) time.Time {
	return timestamp.AddDate(0, 0, 1)
}

type Count struct {
	Count int
}

// 得到当前时间戳的对局数统计
func GetCount(timestamp time.Time) (int, error) {
	body, err := r.Get(fmt.Sprintf("api/count/%d", timestamp.Unix()))

	if err != nil {
		logger.Error("GetCount", zap.Error(err))
		return 0, err
	}

	var count Count
	err = json.Unmarshal(body, &count)
	if err != nil {
		logger.Error("GetCount", zap.Error(err))
		return 0, err
	}
	return count.Count, nil
}

type Record []struct {
	ID        string     `json:"_id"`
	ModeID    int        `json:"modeId"`
	UUID      string     `json:"uuid"`
	StartTime int        `json:"startTime"`
	EndTime   int        `json:"endTime"`
	Players   []*Players `json:"players"`
}
type Players struct {
	AccountID    int    `json:"accountId"`
	Nickname     string `json:"nickname"`
	Level        int    `json:"level"`
	Score        int    `json:"score"`
	GradingScore int    `json:"gradingScore"`
}

func GetRecord(timestamp time.Time, count, mode int) (Record, error) {
	tde := ToDayEndTimestamp(timestamp)
	body, err := r.Get(fmt.Sprintf("api/v2/pl4/games/%d/%d?limit=%d&descending=true&mode=%d", tde.Unix(), timestamp.Unix(), count, mode))
	if err != nil {
		logger.Error("GetRecord", zap.Error(err))
		return nil, err
	}

	var record Record
	err = json.Unmarshal(body, &record)
	if err != nil {
		logger.Error("GetRecord", zap.Error(err))
		return nil, err
	}
	return record, nil
}
