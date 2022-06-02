@echo off
set GOOS=windows
set GOARCH=amd64
go mod tidy
go build -ldflags "-s -w -X main.version=$(git describe --abbrev=0 --tags)" -o mahjong.exe mahjong/main
pause