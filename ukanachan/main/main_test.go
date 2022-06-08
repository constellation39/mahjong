package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"majsoul/message"
	"testing"
)

func TestSaveRecord(t *testing.T) {
	body, _ := ioutil.ReadFile("191201-ae8c22ad-8b71-4544-9ec2-76d7a4f80ff4")

	fmt.Printf("%s \n", body)

	wrapper := new(message.Wrapper)

	proto.Unmarshal(body, wrapper)

	fmt.Println(wrapper.Name)

	ResGameRecord := new(message.ResGameRecord)

	proto.Unmarshal(wrapper.Data, ResGameRecord)

	if ResGameRecord.Head.Config.Category != 2 {
		panic("Category != 2")
	}

	fmt.Println(ResGameRecord.Head.Config.Category)
}
