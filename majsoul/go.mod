module majsoul

require (
	github.com/golang/protobuf v1.5.2
	go.uber.org/zap v1.21.0
	google.golang.org/grpc v1.47.0
	google.golang.org/protobuf v1.28.0
	utils v0.0.0
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)

replace utils => ./../utils

go 1.18
