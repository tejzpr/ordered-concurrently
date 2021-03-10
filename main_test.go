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
	return val
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
		doneChan := make(chan bool)
		outChan := Process(inputChan, workFn, &Options{})
		go func(t *testing.T) {
			counter := 0
			for {
				select {
				case out, chok := <-outChan:
					if chok {
						if _, ok := out.Value.(int); !ok {
							t.Error("Invalid output")
						} else {
							counter++
						}
					} else {
						if counter != max {
							t.Error("Input count does not match output count")
						}
						doneChan <- true
					}
				}
			}
		}(t)

		// Create work and the associated order
		for work := 0; work < max; work++ {

			input := &OrderedInput{work}
			inputChan <- input
		}
		close(inputChan)
		<-doneChan
		t.Log("Test with default Pool Size Completed")
	})
}
