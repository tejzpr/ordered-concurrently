<a href="https://github.com/tejzpr/ordered-concurrently/actions/workflows/tests.yml"><img src="https://github.com/tejzpr/ordered-concurrently/actions/workflows/tests.yml/badge.svg" alt="Tests"/></a>
[![codecov](https://codecov.io/gh/tejzpr/ordered-concurrently/branch/master/graph/badge.svg?token=6WIXWRO3EW)](https://codecov.io/gh/tejzpr/ordered-concurrently)
[![Go Reference](https://pkg.go.dev/badge/github.com/tejzpr/ordered-concurrently.svg)](https://pkg.go.dev/github.com/tejzpr/ordered-concurrently)
[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/tejzpr/ordered-concurrently)
[![Go Report Card](https://goreportcard.com/badge/github.com/tejzpr/ordered-concurrently)](https://goreportcard.com/report/github.com/tejzpr/ordered-concurrently)

# Ordered Concurrently
A library for parallel processing with ordered output in Go. This module processes work concurrently / in parallel and returns output in a channel in the order of input. It is useful in concurrently / parallelly processing items in a queue, and get output in the order provided by the queue.

# Usage 
## Get Module
```go
go get github.com/tejzpr/ordered-concurrently/v2
```
## Import Module in your source code
```go
import concurrently "github.com/tejzpr/ordered-concurrently/v2" 
```
## Create a work function by implementing WorkFunction interface
```go
// Create a type based on your input to the work function
type loadWorker int

// The work that needs to be performed
// The input type should implement the WorkFunction interface
func (w loadWorker) Run() interface{} {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	return w * 2
}
```
## Run
### Example - 1
```go
func main() {
	max := 10
	inputChan := make(chan concurrently.WorkFunction)
	output := concurrently.Process(inputChan, &concurrently.Options{PoolSize: 10, OutChannelBuffer: 10})
	go func() {
		for work := 0; work < max; work++ {
			inputChan <- loadWorker(work)
		}
		close(inputChan)
	}()
	for out := range output {
		log.Println(out.Value)
	}
}
```
### Example - 2 - Process unknown number of inputs
```go
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

	var res []loadWorker
	go func() {
		for out := range output {
			res = append(res, out.Value.(loadWorker))
			wg.Done()
		}
	}()

	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	done <- true
	close(inputChan)
	wg.Wait()

	// Check if output is sorted
	isSorted := sort.SliceIsSorted(res, func(i, j int) bool {
		return res[i] < res[j]
	})
	if !isSorted {
		log.Println("output is not sorted")
	}
}
```
# Credits
1.  [u/justinisrael](https://www.reddit.com/user/justinisrael/) for inputs on improving resource usage.
2.  [mh-cbon](https://github.com/mh-cbon) for identifying potential [deadlocks](https://github.com/tejzpr/ordered-concurrently/issues/2).
