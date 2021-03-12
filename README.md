<a href="https://github.com/tejzpr/ordered-concurrently/actions/workflows/tests.yml"><img src="https://github.com/tejzpr/ordered-concurrently/actions/workflows/tests.yml/badge.svg" alt="Tests"/></a>
[![codecov](https://codecov.io/gh/tejzpr/ordered-concurrently/branch/master/graph/badge.svg?token=6WIXWRO3EW)](https://codecov.io/gh/tejzpr/ordered-concurrently)
[![Go Reference](https://pkg.go.dev/badge/github.com/tejzpr/ordered-concurrently.svg)](https://pkg.go.dev/github.com/tejzpr/ordered-concurrently)
[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/tejzpr/ordered-concurrently)
[![Go Report Card](https://goreportcard.com/badge/github.com/tejzpr/ordered-concurrently)](https://goreportcard.com/report/github.com/tejzpr/ordered-concurrently)

# Ordered Concurrently
A ligrary for parallel processing with ordered output in Go. This module processes work concurrently / in parallel and returns output in a channel in the order of input. It is useful in concurrently / parallelly processing items in a queue, and get output in the order provided by the queue.

# Usage 
## Get Module
```go
go get github.com/tejzpr/ordered-concurrently
```
## Import Module in your source code
```go
import concurrently "github.com/tejzpr/ordered-concurrently" 
```
## Create a work function
```go
// The work that needs to be performed
func workFn(val interface{}) interface{} {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	return val.(int) * 2
}
```
## Run
```go
func main() {
	max := 100
	// Can be a non blocking channel as well
	inputChan := make(chan *concurrently.OrderedInput)
	wg := &sync.WaitGroup{}

	outChan := concurrently.Process(inputChan, workFn, &concurrently.Options{PoolSize: 10})
	go func() {
		for out := range outChan {
			log.Println(out.Value)
			wg.Done()
		}
	}()

	// Create work and sent to input channel
	// Output will be in the order of input
	for work := 0; work < max; work++ {
		wg.Add(1)
		input := &concurrently.OrderedInput{work}
		inputChan <- input
	}
	close(inputChan)
	wg.Wait()
}
```
# Credits
Thanks to [u/justinisrael](https://www.reddit.com/user/justinisrael/) for inputs on improving resource usage.

