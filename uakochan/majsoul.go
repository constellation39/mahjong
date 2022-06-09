package uakochan

import "majsoul"

func NewMajsoul() *majsoul.Majsoul {
	cfg := majsoul.LoadConfig()
	m := majsoul.New(cfg)
	return m
}
