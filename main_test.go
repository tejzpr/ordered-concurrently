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
	t.Run("Basic Tests", func(t *testing.T) {
		max := 10
		inputChan := make(chan *OrderedInput)
		wg := &sync.WaitGroup{}
		go func(t *testing.T) {
			outChan := Process(inputChan, workFn, 10)
			for out := range outChan {
				if _, ok := out.Value.(int); !ok {
					t.Error("Invalid output")
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
	})
}
