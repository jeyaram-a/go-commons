# go commons

Go-commons provides a range of handy tools designed to tackle prevalent concurrency and file operation patterns.

### Router:
Distributes input across provided number of routes. To consume events one has to subscribe to a route. Only one subscriber per route. Input should implement Hashable interface.
```go
// Router example
// Router which routes integers

// Implement Hashable for int
type Integer int
func (x Integer) Hash() uint32 {
	return uint32(x)
}

// Creating a router with 2 routes
router := NewRouter(uint32(2))
// Adding subscriber1
router.Subscribe(func(in Hashable) {
    fmt.Printf("subscriber 1 %v", in)
} )
// Adding subscriber2
router.Subscribe(func(in Hashable) {
    fmt.Printf("subscriber 2 %v", in)
} )
for i:=0; i<10; i++ {
    router.Route(Integer(i))
}
router.Close()
```

### Transformer:
Transformer is a utility to achieve structured concurreny with goroutines. Applies transformation to input slice with a specified parallelism. Results can also be filtered by passing a filter (like eliminating all errors)

```go
// Transformer example

// Transformer downloads input files in 10 parallel goroutines
// and all the erroneous attempts to download are ignored

func download(path string) ([]byte, error) {
    // download
}

filesToBeDownloaded = ["one", "two"]

transformer := NewTransformer(unint(10), filesToBeDownloaded, func(path string) ([]byte, error) {
    return download(path)
})

// start downloading
job := transformer.Transform(func(r gocommons.Result[[]byte]) bool {
        // filtering out all errored download
        if r.Err != nil {
            log.Error(r.Err)
        }
		return r.Err == nil
	})

// blocks until entire input is processed
result := job.Get()

```

### File utils:
```go
// returns true if folder exists else returns false
exists := FolderExists(path)

// err != nil if some error is encountered in stat
isEmpty, err := IsFolderEmpty(path)

```

### Collections:
```go
// Filter => Filters the input arr on the predicate proided
arr := make([]int, 10)
filtered := Filter(arr, func(x int) {
    x%2 == 0
})

// Map
arr := make([]string, 10)
mapped := Map_(arr, func(x string) {
    return len(x)
})

// Contains
arr := make([]string, 10)
present := Contains(arr, "hello")
```
