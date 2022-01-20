package EventBus

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

func Test_RaceCondtion(t *testing.T) {

	consumeMath := func(wg *sync.WaitGroup) chan interface{} {
		wg.Add(1)
		c := make(chan interface{})

		go func() {
			defer wg.Done()
			<-c
		}()

		return c
	}

	consumeMathAgain := func(wg *sync.WaitGroup) chan interface{} {
		wg.Add(1)
		c := make(chan interface{})

		go func() {
			defer wg.Done()
			<-c
		}()

		return c
	}

	consumeText := func(wg *sync.WaitGroup) chan interface{} {
		wg.Add(1)
		c := make(chan interface{})

		go func() {
			defer wg.Done()
			<-c
		}()

		return c
	}

	bus := New()
	defer bus.Close()

	wg := &sync.WaitGroup{}

	bus.Subscribe("math", consumeMath(wg))
	bus.Subscribe("math", consumeMathAgain(wg))
	bus.Subscribe("text", consumeText(wg))

	bus.EmitAsync("math", rand.Int())
	bus.EmitAsync("text", strconv.Itoa(rand.Int()))

	wg.Wait()
}

func Test_Bus(t *testing.T) {

	consumeText := func(result chan<- string) chan interface{} {
		c := make(chan interface{})

		go func() {
			m := <-c
			received := m.([]interface{})
			v := received[0].(string)
			result <- v
		}()

		return c
	}

	bus := New()
	defer bus.Close()

	result := make(chan string)
	test := "Go is awesome"

	bus.Subscribe("text", consumeText(result))
	bus.Emit("text", test)

	r := <-result

	if r != test {
		t.Errorf("expected %s, got %s", test, r)
	}
}

func Test_Subscribe(t *testing.T) {
	bus := New()
	defer bus.Close()

	bus.Subscribe("text", make(chan interface{}), make(chan interface{}))
	bus.Subscribe("awesome", make(chan interface{}), make(chan interface{}), make(chan interface{}))

	if len(bus.Channels["text"]) != 2 {
		t.Errorf("expected 2, got %d", len(bus.Channels["text"]))
	}

	if len(bus.Channels["awesome"]) != 3 {
		t.Errorf("expected 3, got %d", len(bus.Channels["awesome"]))
	}
}

func Test_UnSubscribe(t *testing.T) {
	bus := New()
	defer bus.Close()

	x := make(chan interface{})
	y := make(chan interface{})

	bus.Subscribe("text", x, y)
	bus.Unsubscribe("text", y)

	if len(bus.Channels["text"]) != 1 {
		t.Errorf("expected 1, got %d", len(bus.Channels["text"]))
	}

	if bus.Channels["text"][0] != x {
		t.Errorf("expected %x, got %x", x, bus.Channels["text"][0])
	}
}
