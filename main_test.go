package orderedconcurrently

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

// The work that needs to be performed
func workFn(val interface{}) WorkFunction {
	return func() interface{} {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
		return val.(int) * 2
	}
}

// The work that needs to be performed
func zeroLoadWorkFn(val interface{}) WorkFunction {
	return func() interface{} {
		return val.(int)
	}
}

func Test1(t *testing.T) {
	t.Run("Test with Preset Pool Size", func(t *testing.T) {
		max := 10
		inputChan := make(chan OrderedInput)
		wg := &sync.WaitGroup{}

		outChan := Process(inputChan, &Options{PoolSize: 10})
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
			input := OrderedInput{WorkFn: workFn(work)}
			inputChan <- input
		}
		close(inputChan)
		wg.Wait()
		if counter != max {
			t.Error("Input count does not match output count")
		}
		t.Log("Test with Preset Pool Size Completed")
	})
}

func Test2(t *testing.T) {
	t.Run("Test with default Pool Size", func(t *testing.T) {
		max := 10
		inputChan := make(chan OrderedInput)
		wg := &sync.WaitGroup{}

		outChan := Process(inputChan, &Options{OutChannelBuffer: 2})
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
			input := OrderedInput{WorkFn: workFn(work)}
			inputChan <- input
		}
		close(inputChan)
		wg.Wait()
		if counter != max {
			t.Error("Input count does not match output count")
		}
		t.Log("Test with Default Pool Size Completed")
	})
}

func Test3(t *testing.T) {
	t.Run("Test Zero Load", func(t *testing.T) {
		max := 10
		inputChan := make(chan OrderedInput)
		wg := &sync.WaitGroup{}

		outChan := Process(inputChan, &Options{OutChannelBuffer: 2})
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
			input := OrderedInput{WorkFn: zeroLoadWorkFn(work)}
			inputChan <- input
		}
		close(inputChan)
		wg.Wait()
		if counter != max {
			t.Error("Input count does not match output count")
		}
		t.Log("Test with Default Pool Size and Zero Load Completed")
	})

}

func Test4(t *testing.T) {
	t.Run("Test without workgroup", func(t *testing.T) {
		max := 10
		inputChan := make(chan OrderedInput)
		output := Process(inputChan, &Options{PoolSize: 10, OutChannelBuffer: 10})
		go func() {
			for work := 0; work < max; work++ {
				inputChan <- OrderedInput{WorkFn: zeroLoadWorkFn(work)}
			}
			close(inputChan)
		}()
		counter := 0
		for out := range output {
			if _, ok := out.Value.(int); !ok {
				t.Error("Invalid output")
			} else {
				counter++
			}
		}
		if counter != max {
			t.Error("Input count does not match output count")
		}
		t.Log("Test without workgroup Completed")
	})
}

func BenchmarkOC(b *testing.B) {
	max := 1000000
	inputChan := make(chan OrderedInput)
	output := Process(inputChan, &Options{PoolSize: 10, OutChannelBuffer: 10})
	go func() {
		for work := 0; work < max; work++ {
			inputChan <- OrderedInput{WorkFn: zeroLoadWorkFn(work)}
		}
		close(inputChan)
	}()
	for out := range output {
		_ = out
	}
}
