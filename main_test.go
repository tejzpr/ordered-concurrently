package orderedconcurrently

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

type zeroLoadWorker int

func (w zeroLoadWorker) Run() interface{} {
	return w * 2
}

type loadWorker int

func (w loadWorker) Run() interface{} {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	return w * 2
}

func Test1(t *testing.T) {
	t.Run("Test with Preset Pool Size", func(t *testing.T) {
		max := 10
		inputChan := make(chan WorkFunction)
		wg := &sync.WaitGroup{}

		outChan := Process(inputChan, &Options{PoolSize: 10})
		counter := 0
		go func(t *testing.T) {
			for out := range outChan {
				if _, ok := out.Value.(loadWorker); !ok {
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
			inputChan <- loadWorker(work)
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
		inputChan := make(chan WorkFunction)
		wg := &sync.WaitGroup{}

		outChan := Process(inputChan, &Options{OutChannelBuffer: 2})
		counter := 0
		go func(t *testing.T) {
			for out := range outChan {
				if _, ok := out.Value.(loadWorker); !ok {
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
			inputChan <- loadWorker(work)
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
		inputChan := make(chan WorkFunction)
		wg := &sync.WaitGroup{}

		outChan := Process(inputChan, &Options{OutChannelBuffer: 2})
		counter := 0
		go func(t *testing.T) {
			for out := range outChan {
				if _, ok := out.Value.(zeroLoadWorker); !ok {
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
			inputChan <- zeroLoadWorker(work)
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
		inputChan := make(chan WorkFunction)
		output := Process(inputChan, &Options{PoolSize: 10, OutChannelBuffer: 10})
		go func() {
			for work := 0; work < max; work++ {
				inputChan <- zeroLoadWorker(work)
			}
			close(inputChan)
		}()
		counter := 0
		for out := range output {
			if _, ok := out.Value.(zeroLoadWorker); !ok {
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
	inputChan := make(chan WorkFunction)
	output := Process(inputChan, &Options{PoolSize: 10, OutChannelBuffer: 10})
	go func() {
		for work := 0; work < max; work++ {
			inputChan <- zeroLoadWorker(work)
		}
		close(inputChan)
	}()
	for out := range output {
		_ = out
	}
}
