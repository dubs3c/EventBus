package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {

	test1 := func(ctx context.Context, wg *sync.WaitGroup) chan interface{} {
		c := make(chan interface{})
		wg.Add(1)
		go func() {
			defer wg.Done()
			d := make(chan interface{})
			// simluate some kind of work
			go func() {
				v := <-c
				d <- v
			}()

			select {
			case <-ctx.Done():
				fmt.Println("[CANCELED] this is test one ")
			case v := <-d:
				fmt.Println("[SUCCESS] this is test one, got:", v.(string))
			}

		}()
		return c
	}

	test2 := func(ctx context.Context, wg *sync.WaitGroup) chan interface{} {
		c := make(chan interface{})
		wg.Add(1)
		go func() {
			defer wg.Done()
			d := make(chan interface{}, 1)
			// simluate some kind of work
			go func() {
				time.Sleep(3 * time.Second)
				// this goroutine will leak :(
				d <- <-c
			}()

			select {
			case <-ctx.Done():
				fmt.Println("[CANCELED] this is test two ")
				// is this an ugly hack?
				// If context timeout, bug.go:32 will block and crash when trying to send on a closed channel
				// by reading the value below, we drain the channel and make sure no blocks can happen
				// and that the goroutine wont try to send on a blocked channel
				<-c
			case v := <-d:
				fmt.Println("[SUCCESS] this is test two, got:", v.(string))
			}

		}()
		return c
	}

	test3 := func(ctx context.Context, wg *sync.WaitGroup) chan interface{} {
		c := make(chan interface{})
		wg.Add(1)
		go func() {
			defer wg.Done()
			d := make(chan interface{})
			// simluate some kind of work
			go func() {
				v := <-c
				d <- v
			}()

			select {
			case <-ctx.Done():
				fmt.Println("[CANCELED] this is test three ")
			case v := <-d:
				fmt.Println("[SUCCESS] this is test three, got:", v.(string))
			}

		}()
		return c
	}

	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	bus := NewBus()
	defer bus.Close()
	bus.Ctx = ctx

	bus.Subscribe("cool", test1(ctx, wg))
	bus.Subscribe("cool", test2(ctx, wg))
	bus.Subscribe("cool", test3(ctx, wg))

	bus.EmitAsync("cool", "hey")

	wg.Wait()

	fmt.Println("Finished work")

}
