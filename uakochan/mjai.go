package uakochan

import (
	"go.uber.org/zap"
	"utils/logger"
)

// MJAI is a message type.
type Type struct {
	Type string `json:"type"`
}

// NewMJAI 握手
type Hello struct {
	Type            string `json:"type"`
	Protocol        string `json:"protocol"`
	ProtocolVersion int    `json:"protocol_version"`
}

// MJAI 握手
type Join struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Room string `json:"room"`
}

// MJAI 游戏开始
type StartGame struct {
	Type  string   `json:"type"`
	ID    int      `json:"id"`
	Names []string `json:"names"`
}

// MJAI 无
type None struct {
	Type string `json:"type"`
}

// MJAI 新的一轮开始.
type StartKyoku struct {
	Type       string     `json:"type"`
	Bakaze     string     `json:"bakaze"`
	Kyoku      uint32     `json:"kyoku"`
	Honba      uint32     `json:"honba"`
	Kyotaku    uint32     `json:"kyotaku"`
	Oya        uint32     `json:"oya"`
	DoraMarker string     `json:"dora_marker"`
	Tehais     [][]string `json:"tehais"`
}

// MJAI 摸牌
type Tsumo struct {
	Type  string `json:"type"`
	Actor int    `json:"actor"`
	Pai   string `json:"pai"`
}

// MJAI 打牌
type Dahai struct {
	Type      string `json:"type"`
	Actor     int    `json:"actor"`
	Pai       string `json:"pai"`
	Tsumogiri bool   `json:"tsumogiri"`
}

// MJAI 碰
type Pon struct {
	Type     string   `json:"type"`
	Actor    int      `json:"actor"`
	Target   int      `json:"target"`
	Pai      string   `json:"pai"`
	Consumed []string `json:"consumed"`
}

// MJAI 吃
type Chi struct {
	Type     string   `json:"type"`
	Actor    int      `json:"actor"`
	Target   int      `json:"target"`
	Pai      string   `json:"pai"`
	Consumed []string `json:"consumed"`
}

// MJAI 杠
type Kakan struct {
	Type     string   `json:"type"`
	Actor    int      `json:"actor"`
	Pai      string   `json:"pai"`
	Consumed []string `json:"consumed"`
}

// MJAI 明杠
type Daiminkan struct {
	Type     string   `json:"type"`
	Actor    int      `json:"actor"`
	Target   int      `json:"target"`
	Pai      string   `json:"pai"`
	Consumed []string `json:"consumed"`
}

// MJAI 暗杠
type Ankan struct {
	Type     string   `json:"type"`
	Actor    int      `json:"actor"`
	Consumed []string `json:"consumed"`
}

// MJAI 立直
type Reach struct {
	Type  string `json:"type"`
	Actor int    `json:"actor"`
}

// MJAI 立直的结果
type ReachAccepted struct {
	Type   string `json:"type"`
	Actor  int    `json:"actor"`
	Deltas []int  `json:"deltas"`
	Scores []int  `json:"scores"`
}

// MJAI 和牌
type Hora struct {
	Type           string          `json:"type"`
	Actor          int             `json:"actor"`
	Target         int             `json:"target"`
	Pai            string          `json:"pai"`
	UradoraMarkers []string        `json:"uradora_markers"`
	HoraTehais     []string        `json:"hora_tehais"`
	Yakus          [][]interface{} `json:"yakus"`
	Fu             int             `json:"fu"`
	Fan            int             `json:"fan"`
	HoraPoints     int             `json:"hora_points"`
	Deltas         []int           `json:"deltas"`
	Scores         []int           `json:"scores"`
}

// MJAI 结束一轮
type EndKyoku struct {
	Type string `json:"type"`
}

// MJAI 流局
type Ryukyoku struct {
	Type    string     `json:"type"`
	Reason  string     `json:"reason"`
	Tehais  [][]string `json:"tehais"`
	Tenpais []bool     `json:"tenpais"`
	Deltas  []int      `json:"deltas"`
	Scores  []int      `json:"scores"`
}

func GetMessageType(name string) (ret interface{}) {
	switch name {
	case "hello":
		ret = new(Hello)
	case "join":
		ret = new(Join)
	case "start_game":
		ret = new(StartGame)
	case "none":
		ret = new(None)
	case "start_kyoku":
		ret = new(StartKyoku)
	case "tsumo":
		ret = new(Tsumo)
	case "dahai":
		ret = new(Dahai)
	case "pon":
		ret = new(Pon)
	case "chi":
		ret = new(Chi)
	case "kakan":
		ret = new(Kakan)
	case "daiminkan":
		ret = new(Daiminkan)
	case "ankan":
		ret = new(Ankan)
	case "reach":
		ret = new(Reach)
	case "reach_accepted":
		ret = new(ReachAccepted)
	case "hora":
		ret = new(Hora)
	case "end_kyoku":
		ret = new(EndKyoku)
	case "ryukyoku":
		ret = new(Ryukyoku)
	default:
		logger.Error("GetMessageType not found", zap.String("name", name))
	}
	return
}
