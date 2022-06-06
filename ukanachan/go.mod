module ukanachan

require (
	go.uber.org/zap v1.21.0
	utils v0.0.0
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
)

replace utils => ./../utils

go 1.18
