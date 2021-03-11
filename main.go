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
	PoolSize         int
	OutChannelBuffer int
}

// WorkFunction the function which performs work
type WorkFunction func(interface{}) interface{}

// Process processes work function based on input.
// It Accepts an OrderedInput read channel, work function and concurrent go routine pool size.
// It Returns an OrderedOutput channel.
func Process(inputChan <-chan *OrderedInput, wf WorkFunction, options *Options) <-chan *OrderedOutput {
	outputChan := make(chan *OrderedOutput, options.OutChannelBuffer)
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
		processChan := make(chan *processInput)
		aggregatorChan := make(chan *processInput)
		wg := sync.WaitGroup{}
		doneSemaphoreChan := make(chan bool)
		// Go routine to print data in order
		go func() {
			var current uint64
			outputMap := make(map[uint64]*processInput)
			for {
				select {
				case item, ok := <-aggregatorChan:
					if ok {
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
					} else {
						aggregatorChan = nil
					}
				}
				if aggregatorChan == nil {
					close(outputChan)
					doneSemaphoreChan <- true
				}
			}
		}()

		inputClosedSemaphoreChan := make(chan bool)
		// Create a goroutine pool
		for i := 0; i < processors; i++ {
			go func() {
				for input := range processChan {
					wg.Add(1)
					input.value = wf(input.value)
					input.wg = &wg
					aggregatorChan <- input
					select {
					case <-inputClosedSemaphoreChan:
						wg.Wait()
						close(aggregatorChan)
					default:
						continue
					}
				}
			}()
		}

		var order uint64
		for {
			select {
			case input, ok := <-inputChan:
				if ok {
					processChan <- &processInput{input.Value, order, nil}
					order++
				} else {
					inputChan = nil
				}
			}
			if inputChan == nil {
				close(processChan)
				inputClosedSemaphoreChan <- true
				break
			}
		}
		<-doneSemaphoreChan
	}()
	return outputChan
}
