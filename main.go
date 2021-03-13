package orderedconcurrently

import (
	"container/heap"
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

type processInput struct {
	workFn WorkFunction
	order  uint64
	value  interface{}
}

// WorkFunction the function which performs work
type WorkFunction func() interface{}

// Process processes work function based on input.
// It Accepts an OrderedInput read channel, work function and concurrent go routine pool size.
// It Returns an OrderedOutput channel.
func Process(inputChan <-chan OrderedInput, options *Options) <-chan OrderedOutput {

	outputChan := make(chan OrderedOutput, options.OutChannelBuffer)

	go func() {
		if options.PoolSize < 1 {
			// Set a minimum number of processors
			options.PoolSize = 1
		}
		processChan := make(chan processInput, options.PoolSize)
		aggregatorChan := make(chan processInput, options.PoolSize)
		doneSemaphoreChan := make(chan bool)

		// Go routine to print data in order
		go func() {
			var current uint64
			outputHeap := &processInputHeap{}
			defer func() {
				close(outputChan)
				doneSemaphoreChan <- true
			}()

			for item := range aggregatorChan {
				tmp := &item
				heap.Push(outputHeap, tmp)
				top := outputHeap.Peek()
				if top != nil && top.(*processInput).order == current {
					outputChan <- OrderedOutput{Value: heap.Pop(outputHeap).(*processInput).value}
					current++
					continue
				}
			}

			for outputHeap.Len() > 0 {
				outputChan <- OrderedOutput{Value: heap.Pop(outputHeap).(*processInput).value}
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
				processChan <- processInput{workFn: input.WorkFn, order: order}
				order++
			}
		}()

		<-doneSemaphoreChan
	}()
	return outputChan
}
