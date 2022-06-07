@echo off
set projectName=paipu
set GOOS=linux
set outFile=bin\%projectName%.linux
go mod tidy -compat=1.17
set GOARCH=amd64
go build %buildFlags% -o ../%outFile%  .