package utils

import (
	"fmt"
	"time"
)

var Env string
var Version string
var BuildTime string
var CommitID string
var Branch string
var BuildUser string

func init() {
	fmt.Printf(`%+v	Build	{"Env": %s, "Version": %s, "BuildTime": %s, "CommitID": %s, "Branch": %s, "BuildUser": %s} %s`, time.Now().Format("2006-01-02T03:04:05.999+0800"), Env, Version, BuildTime, CommitID, Branch, BuildUser, "\n")
}
