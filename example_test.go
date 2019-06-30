package awaiter_test

import (
	"fmt"
	"time"

	"github.com/v-braun/awaiter"
)

func ExampleAwaiter_Go() {
	awaiter := awaiter.New()

	counter := 0

	awaiter.Go(func() {
		time.Sleep(time.Second * 1)
		counter += 1
	})

	awaiter.Go(func() {
		time.Sleep(time.Second * 1)
		counter += 1
	})

	awaiter.AwaitSync()

	fmt.Printf("counter: %d", counter)
	// Output: counter: 2
}

func ExampleAwaiter_Cancel() {
	awaiter := awaiter.New()

	counter := 0

	awaiter.Go(func() {
	loop:
		for {
			select {
			case <-awaiter.CancelRequested():
				break loop
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}

		counter += 1
	})

	awaiter.Go(func() {
	loop:
		for {
			if awaiter.IsCancelRequested() {
				break loop
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}

		counter += 1
	})

	time.Sleep(time.Second * 1)

	awaiter.Cancel()
	awaiter.Cancel() // call Cancel multiple times is ok
	awaiter.AwaitSync()

	fmt.Printf("counter: %d", counter)
	// Output: counter: 2
}
