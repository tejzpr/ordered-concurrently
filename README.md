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
go get github.com/tejzpr/ordered-concurrently
```
## Import Module in your source code
```go
import concurrently "github.com/tejzpr/ordered-concurrently" 
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
		log.Println(out)
	}
}
```
### Example - 2
```go
func main() {
	max := 100
	// Can be a non blocking channel as well
	inputChan := make(chan concurrently.WorkFunction)
	wg := &sync.WaitGroup{}

	outChan := concurrently.Process(inputChan, &concurrently.Options{OutChannelBuffer: 2})
	go func() {
		for out := range outChan {
			log.Println(out)
			wg.Done()
		}
	}()

	// Create work and sent to input channel
	// Output will be in the order of input
	for work := 0; work < max; work++ {
		wg.Add(1)
		input := loadWorker(work)
		inputChan <- input
	}
	close(inputChan)
	wg.Wait()
}
```
# Credits
1.  [u/justinisrael](https://www.reddit.com/user/justinisrael/) for inputs on improving resource usage.
2.  [mh-cbon](https://github.com/mh-cbon) for identifying potential [deadlocks](https://github.com/tejzpr/ordered-concurrently/issues/2).
