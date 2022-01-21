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
	channels  map[string][]chan interface{}
	ctx       context.Context
	ctxCancel context.CancelFunc
	waitGroup *sync.WaitGroup
	mutex     *sync.Mutex
}

func New() *EventBus {
	ctx, ctxCancel := context.WithCancel(context.Background())

	return &EventBus{
		channels:  make(map[string][]chan interface{}),
		ctx:       ctx,
		ctxCancel: ctxCancel,
		waitGroup: &sync.WaitGroup{},
		mutex:     &sync.Mutex{},
	}
}

// Emit emits a value onto every channel. Will block until value is received
func (e *EventBus) Emit(topic string, args ...interface{}) {
	for _, channel := range e.channels[topic] {
		channel <- args
	}
}

// Emit emits a value onto every channel in a fan-out fashion. Non-blocking
func (e *EventBus) EmitAsync(topic string, args ...interface{}) {
	go func() {
		for i := 0; i < len(e.channels[topic]); i++ {
			e.waitGroup.Add(1)
			go func(i int) {
				defer e.waitGroup.Done()
				channel := e.channels[topic][i]
				select {
				case <-e.ctx.Done():
					return
				case channel <- args:
				}
			}(i)
		}
	}()
}

// Subscribe a channel or multiple channels to a given topic
func (e *EventBus) Subscribe(topic string, channels ...chan interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.channels[topic] = append(e.channels[topic], channels...)
}

// Unsubscribe a channel from a given topic
func (e *EventBus) Unsubscribe(topic string, channel chan interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for i := 0; i < len(e.channels[topic]); i++ {
		if e.channels[topic][i] == channel {
			copy(e.channels[topic][i:], e.channels[topic][i+1:])
			e.channels[topic] = e.channels[topic][:len(e.channels[topic])-1]
			break
		}
	}
}

// CountSubscribers returns the amount of subscribers for a given topic
func (e *EventBus) CountSubscribers(topic string) int {
	e.mutex.Lock()
	c := len(e.channels[topic])
	e.mutex.Unlock()
	return c
}

// GetSubscribers returns a list of channels subscribed to a topic
func (e *EventBus) GetSubscribers(topic string) []chan interface{} {
	e.mutex.Lock()
	c := e.channels[topic]
	e.mutex.Unlock()
	return c
}

// RemoveTopic will remove a given topic and its subscribers from the topic list
func (e *EventBus) RemoveTopic(topic string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	delete(e.channels, topic)

}

// Close closes all channels subscribed to a topic
func (e *EventBus) Close() {
	e.ctxCancel()
	e.waitGroup.Wait()
	for _, topic := range e.channels {
		for _, channel := range topic {
			close(channel)
		}
	}
}
