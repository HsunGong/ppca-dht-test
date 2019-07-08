## prepare

go get -u -v github.com/fatih/color

## Why `Run` need `WaitGroup`

> Give a time to free up space

```go
func (n *dhtNode) Run(wg *sync.WaitGroup) {
	defer func() {
		// your code
		wg.Done()
    }()
```