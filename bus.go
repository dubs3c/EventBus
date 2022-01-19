package main

import (
	"context"
)

type Bus interface {
	Close()
	Emit()
	EmitAsync()
	Subscribe()
	Unsubscribe()
}

type EventBus struct {
	Channels map[string][]chan interface{}
	Ctx      context.Context
}

func NewBus() *EventBus {
	m := make(map[string][]chan interface{})
	return &EventBus{
		Channels: m,
		Ctx:      context.Context(context.Background()),
	}
}

// Emit emits a value onto every channel
func (e *EventBus) Emit(topic string, args ...interface{}) {
	for _, channel := range e.Channels[topic] {
		channel <- args
	}
}

// Emit emits a value onto every channel
func (e *EventBus) EmitAsync(topic string, args interface{}) {
	for i := 0; i < len(e.Channels[topic]); i++ {
		go func(i int) {
			channel := e.Channels[topic][i]
			// There is an issue here
			// If channel is not receiving, it will block for ever
			// If the channel is closed, a panic will arise
			// maybe use context to check if we should abort - did not work, or use it internally??
			// use wait groups?? maybe together with a done channel or ctx <----
			// https://www.leolara.me/blog/closing_a_go_channel_written_by_several_goroutines/
			channel <- args
		}(i)
	}
}

func (e *EventBus) Subscribe(topic string, channels ...chan interface{}) {
	e.Channels[topic] = append(e.Channels[topic], channels...)
}

func (e *EventBus) Unsubscribe(topic string, channel chan interface{}) {

}

func (e *EventBus) Close() {
	for _, topic := range e.Channels {
		for _, channel := range topic {
			close(channel)
		}
	}
}
