module ukanachan

require (
	go.uber.org/zap v1.21.0
	logger v0.0.0
)

require (
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
)

replace (
	logger => ./../utils/logger
	utils => ./../utils
)

go 1.18
