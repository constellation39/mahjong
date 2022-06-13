package main

import (
	"majsoul"
	"uakochan"
)

type Majsoul struct {
	*majsoul.Majsoul
	*uakochan.UAkochan
	Seat int
}

func NewMajsoul() *Majsoul {
	cfg := majsoul.LoadConfig()
	//m := &Majsoul{Majsoul: majsoul.New(cfg)}
	m := &Majsoul{Majsoul: majsoul.New(cfg), UAkochan: uakochan.New()}
	m.IFReceive = m
	return m
}

func main() {
	mSoul := NewMajsoul()
	mSoul.Login("1601198895@qq.com", "miku39..")
	select {
	case <-mSoul.Ctx.Done():
	}
}
