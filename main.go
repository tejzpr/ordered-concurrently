package orderedconcurrently

import (
	"sync"
)

// OrderedInput input for Processing
type OrderedInput struct {
	WorkFn WorkFunction
}

// OrderedOutput is the output channel type from Process
type OrderedOutput struct {
	Value interface{}
}

// Options options for Process
type Options struct {
	PoolSize         int
	OutChannelBuffer int
}

// WorkFunction interface
type WorkFunction interface {
	Run() interface{}
}

// Process processes work function based on input.
// It Accepts an OrderedInput read channel, work function and concurrent go routine pool size.
// It Returns an OrderedOutput channel.
func Process(inputChan <-chan WorkFunction, options *Options) <-chan OrderedOutput {
	type processInput struct {
		workFn WorkFunction
		order  uint64
		value  interface{}
	}
	outputChan := make(chan OrderedOutput, options.OutChannelBuffer)

	go func() {
		if options.PoolSize < 1 {
			// Set a minimum number of processors
			options.PoolSize = 1
		}
		processChan := make(chan processInput, options.PoolSize)
		aggregatorChan := make(chan processInput, options.PoolSize)

		// Go routine to print data in order
		go func() {
			var current uint64
			var tmp *processInput
			outputMap := make(map[uint64]*processInput, options.PoolSize)
			defer func() {
				close(outputChan)
			}()
			for item := range aggregatorChan {
				tmp = &item
				if item.order != current {
					outputMap[item.order] = tmp
					continue
				}
				for {
					if tmp == nil {
						break
					}
					outputChan <- OrderedOutput{Value: tmp.value}
					delete(outputMap, current)
					current++
					tmp = outputMap[current]
				}
			}
		}()

		poolWg := sync.WaitGroup{}
		poolWg.Add(options.PoolSize)
		// Create a goroutine pool
		for i := 0; i < options.PoolSize; i++ {
			go func(worker int) {
				defer func() {
					poolWg.Done()
				}()
				for input := range processChan {
					input.value = input.workFn.Run()
					input.workFn = nil
					aggregatorChan <- input
				}
			}(i)
		}

		go func() {
			poolWg.Wait()
			close(aggregatorChan)
		}()

		go func() {
			defer func() {
				close(processChan)
			}()
			var order uint64
			for input := range inputChan {
				processChan <- processInput{workFn: input, order: order}
				order++
			}
		}()
	}()
	return outputChan
}
