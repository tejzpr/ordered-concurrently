package orderedconcurrently

import (
	"container/heap"
	"sync"
)

// Options options for Process
type Options struct {
	PoolSize         int
	OutChannelBuffer int
}

// OrderedOutput is the output channel type from Process
type OrderedOutput struct {
	Value     interface{}
	Remaining func() int
}

// WorkFunction interface
type WorkFunction func() interface{}

// Process processes work function based on input.
// It Accepts an WorkFunction read channel, work function and concurrent go routine pool size.
// It Returns an interface{} channel.
func Process(inputChan <-chan WorkFunction, options *Options) <-chan OrderedOutput {

	outputChan := make(chan OrderedOutput, options.OutChannelBuffer)

	go func() {
		if options.PoolSize < 1 {
			// Set a minimum number of processors
			options.PoolSize = 1
		}
		processChan := make(chan *processInput, options.PoolSize)
		aggregatorChan := make(chan *processInput, options.PoolSize)

		// Go routine to print data in order
		go func() {
			var current uint64
			outputHeap := &processInputHeap{}
			defer func() {
				close(outputChan)
			}()
			remaining := func() int {
				return outputHeap.Len()
			}
			for item := range aggregatorChan {
				heap.Push(outputHeap, item)
				for top, ok := outputHeap.Peek(); ok && top.order == current; {
					outputChan <- OrderedOutput{Value: heap.Pop(outputHeap).(*processInput).value, Remaining: remaining}
					current++
				}
			}

			for outputHeap.Len() > 0 {
				outputChan <- OrderedOutput{Value: heap.Pop(outputHeap).(*processInput).value, Remaining: remaining}
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
					input.value = input.workFn()
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
				processChan <- &processInput{workFn: input, order: order}
				order++
			}
		}()
	}()
	return outputChan
}
