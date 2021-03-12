package orderedconcurrently

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

// The work that needs to be performed
func workFn(val interface{}) interface{} {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	return val.(int) * 2
}

// The work that needs to be performed
func zeroLoadWorkFn(val interface{}) interface{} {
	return val.(int)
}

func Test(t *testing.T) {
	t.Run("Test with Preset Pool Size", func(t *testing.T) {
		max := 10
		inputChan := make(chan *OrderedInput)
		wg := &sync.WaitGroup{}

		outChan := Process(inputChan, workFn, &Options{PoolSize: 10})
		counter := 0
		go func(t *testing.T) {
			for out := range outChan {
				if _, ok := out.Value.(int); !ok {
					t.Error("Invalid output")
				} else {
					counter++
				}
				wg.Done()
			}
		}(t)

		// Create work and the associated order
		for work := 0; work < max; work++ {
			wg.Add(1)
			input := &OrderedInput{work}
			inputChan <- input
		}
		wg.Wait()
		if counter != max {
			t.Error("Input count does not match output count")
		}
		t.Log("Test with Preset Pool Size Completed")
	})
	t.Run("Test with default Pool Size", func(t *testing.T) {
		max := 10
		inputChan := make(chan *OrderedInput)
		wg := &sync.WaitGroup{}

		outChan := Process(inputChan, workFn, &Options{OutChannelBuffer: 2})
		counter := 0
		go func(t *testing.T) {
			for out := range outChan {
				if _, ok := out.Value.(int); !ok {
					t.Error("Invalid output")
				} else {
					counter++
				}
				wg.Done()
			}
		}(t)

		// Create work and the associated order
		for work := 0; work < max; work++ {
			wg.Add(1)
			input := &OrderedInput{work}
			inputChan <- input
		}
		wg.Wait()
		if counter != max {
			t.Error("Input count does not match output count")
		}
		t.Log("Test with Default Pool Size Completed")
	})
	t.Run("Test Zero Load", func(t *testing.T) {
		max := 10
		inputChan := make(chan *OrderedInput)
		wg := &sync.WaitGroup{}

		outChan := Process(inputChan, zeroLoadWorkFn, &Options{OutChannelBuffer: 2})
		counter := 0
		go func(t *testing.T) {
			for out := range outChan {
				if _, ok := out.Value.(int); !ok {
					t.Error("Invalid output")
				} else {
					counter++
				}
				wg.Done()
			}
		}(t)

		// Create work and the associated order
		for work := 0; work < max; work++ {
			wg.Add(1)
			input := &OrderedInput{work}
			inputChan <- input
		}
		wg.Wait()
		if counter != max {
			t.Error("Input count does not match output count")
		}
		t.Log("Test with Default Pool Size and Zero Load Completed")
	})
}
