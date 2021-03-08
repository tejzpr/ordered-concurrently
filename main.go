package orderedconcurrently

import (
	"sync"
)

// OrderedInput input for Processing
type OrderedInput struct {
	Value interface{}
}

// OrderedOutput is the output channel type from Process
type OrderedOutput struct {
	Value interface{}
}

// Options options for Process
type Options struct {
	PoolSize int
}

// WorkFunction the function which performs work
type WorkFunction func(interface{}) interface{}

// Process processes work function based on input.
// It Accepts an OrderedInput read channel, work function and concurrent go routine pool size.
// It Returns an OrderedOutput channel.
func Process(inputChan <-chan *OrderedInput, wf WorkFunction, options *Options) <-chan *OrderedOutput {
	outputChan := make(chan *OrderedOutput)
	type processInput struct {
		value interface{}
		order uint64
		wg    *sync.WaitGroup
	}
	go func() {
		processors := options.PoolSize
		if processors == 0 {
			// Set a minimum number of processors
			processors = 1
		}
		wg := sync.WaitGroup{}
		processChan := make(chan *processInput)

		aggregatorChan := make(chan *processInput)

		// Go routine to print data in order
		go func() {
			var current uint64
			outputMap := make(map[uint64]*processInput)
			for item := range aggregatorChan {
				if item.order != current {
					outputMap[item.order] = item
					continue
				}
				for {
					if item == nil {
						break
					}
					outputChan <- &OrderedOutput{Value: item.value}
					item.wg.Done()
					delete(outputMap, current)
					current++
					item = outputMap[current]
				}
			}
		}()

		// Create a goroutine pool
		for i := 0; i < processors; i++ {
			go func() {
				for input := range processChan {
					// Process work
					input.value = wf(input.value)
					aggregatorChan <- input
				}
			}()
		}

		var order uint64
		for input := range inputChan {
			wg.Add(1)
			processChan <- &processInput{input.Value, order, &wg}
			order++
		}
		// Wait till execution is complete
		wg.Wait()
		close(outputChan)
	}()
	return outputChan
}
