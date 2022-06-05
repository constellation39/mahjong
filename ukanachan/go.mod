module ukanachan

require (
	config v0.0.0
	go.uber.org/zap v1.21.0
	logger v0.0.0
	request v0.0.0
)

require (
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
)

replace (
	config => ../config
	logger => ../logger
	request => ../request
)

go 1.18
