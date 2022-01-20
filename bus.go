package EventBus

import (
	"context"
	"sync"
)

type Bus interface {
	Close()
	Emit()
	EmitAsync()
	Subscribe()
	Unsubscribe()
}

type EventBus struct {
	Channels  map[string][]chan interface{}
	Ctx       context.Context
	CtxCancel context.CancelFunc
	WaitGroup *sync.WaitGroup
}

func New() *EventBus {
	m := make(map[string][]chan interface{})
	ctx, ctxCancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	return &EventBus{
		Channels:  m,
		Ctx:       ctx,
		CtxCancel: ctxCancel,
		WaitGroup: wg,
	}
}

// Emit emits a value onto every channel. Will block until value is received
func (e *EventBus) Emit(topic string, args ...interface{}) {
	for _, channel := range e.Channels[topic] {
		channel <- args
	}
}

// Emit emits a value onto every channel in a fan-out fashion. Non-blocking
func (e *EventBus) EmitAsync(topic string, args ...interface{}) {
	go func() {
		for i := 0; i < len(e.Channels[topic]); i++ {
			e.WaitGroup.Add(1)
			go func(i int) {
				defer e.WaitGroup.Done()
				channel := e.Channels[topic][i]
				select {
				case <-e.Ctx.Done():
					return
				case channel <- args:
				}
			}(i)
		}
	}()
}

func (e *EventBus) Subscribe(topic string, channels ...chan interface{}) {
	e.Channels[topic] = append(e.Channels[topic], channels...)
}

func (e *EventBus) Unsubscribe(topic string, channel chan interface{}) {
	for i := 0; i < len(e.Channels[topic]); i++ {
		if e.Channels[topic][i] == channel {
			copy(e.Channels[topic][i:], e.Channels[topic][i+1:])
			e.Channels[topic] = e.Channels[topic][:len(e.Channels[topic])-1]
			break
		}
	}
}

func (e *EventBus) Close() {
	e.CtxCancel()
	e.WaitGroup.Wait()
	for _, topic := range e.Channels {
		for _, channel := range topic {
			close(channel)
		}
	}
}
