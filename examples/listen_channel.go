package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	concurrently "github.com/tejzpr/ordered-concurrently"
)

type loadWorker int

func (w loadWorker) Run() interface{} {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	return w * 2
}

func main() {
	inputChan := make(chan concurrently.WorkFunction, 10)
	output := concurrently.Process(inputChan, &concurrently.Options{PoolSize: 10, OutChannelBuffer: 10})

	ticker := time.NewTicker(100 * time.Millisecond)
	done := make(chan bool)
	wg := &sync.WaitGroup{}
	go func() {
		input := 0
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				inputChan <- loadWorker(input)
				wg.Add(1)
				input++
			default:
			}
		}
	}()

	go func() {
		for out := range output {
			log.Println(out.Value.(loadWorker))
			wg.Done()
		}
	}()

	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	close(inputChan)
	done <- true
	wg.Wait()
}
