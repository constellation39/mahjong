package proc

import (
	"context"
	"log"
	"testing"
	"time"
)

type Base struct {
}

func (b Base) Proc(i interface{}) {
	log.Printf("Call(%+v)\n", i)
}

func (b Base) Tick(t time.Time) {
	log.Printf("Time(%+v)\n", t)
}

func TestMain(m *testing.M) {

	proc := NewProc(context.Background(), &Base{})
	go proc.Listen()

	data := make(chan interface{}, 100)

	for i := 0; i < 10; i++ {
		data <- i
	}
	data <- "end"

	for datum := range data {
		proc.msgCh <- datum
		time.Sleep(time.Second)
	}

	m.Run()
}

func TestNewProc(t *testing.T) {

}

func TestProc_Listen(t *testing.T) {

}
