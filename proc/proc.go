package proc

import (
	"context"
	"time"
)

type IFService interface {
	Proc(interface{})
	Tick(time.Time)
}

type Proc struct {
	ctx     context.Context
	cancel  context.CancelFunc
	service IFService
	msgCh   chan interface{}
}

func NewProc(ctx context.Context, service IFService) *Proc {
	proc := &Proc{
		ctx:     nil,
		cancel:  nil,
		service: service,
		msgCh:   make(chan interface{}, 32),
	}

	proc.ctx, proc.cancel = context.WithCancel(ctx)
	return proc
}

func (proc *Proc) Listen() {
	tk := time.NewTicker(time.Second)
	for {
		select {
		case msg := <-proc.msgCh:
			proc.service.Proc(msg)
		case t := <-tk.C:
			proc.service.Tick(t)
		}
	}
}
