package orderedconcurrently

import (
	"sync"
)

// OrderedInput input for Process
type OrderedInput struct {
	Value interface{}
	Order int
}

// OrderedOutput is the output channel type from Process
type OrderedOutput struct {
	Value interface{}
}

// WorkFunction the function which performs work
type WorkFunction func(interface{}) interface{}

// Process processes work function based on input
func Process(inputChan <-chan *OrderedInput, wf WorkFunction, poolSize int) <-chan *OrderedOutput {
	outputChan := make(chan *OrderedOutput)
	type processInput struct {
		value interface{}
		order int
		wg    *sync.WaitGroup
	}
	go func() {
		processors := poolSize
		wg := sync.WaitGroup{}
		processChan := make(chan *processInput)

		aggregatorChan := make(chan *processInput)

		// Go routine to print data in order
		go func() {
			current := 0
			outputMap := make(map[int]*processInput)
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

		for input := range inputChan {
			wg.Add(1)
			processChan <- &processInput{input.Value, input.Order, &wg}
		}
		// Wait till execution is complete
		wg.Wait()
		close(outputChan)
	}()
	return outputChan
}
