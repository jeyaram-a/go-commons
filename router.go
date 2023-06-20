package gocommons

import (
	"sync"
)

// Hashable. implement Hash()
// Usage:
// 		type Integer int
// 		func (x Integer) Hash() uint32 {
// 			return uint32(x % 5)
// 		}

type Hashable interface {
	Hash() uint32
}

// Router routes Hashable input to a route.
// Parallelism is the number of unique routes.
// Each of these unique routes can be consumed
// parallely
//
// Route logic:
//
//	route = Hashable.hash() % parallelism
//
// Usage:
//
//	router := NewRouter(uint32(2))
//	router.Subscribe(func(in Hashable) {
//		fmt.Printf("%v", in)
//	} )
//	for i:=0; i<10; i++ {
//		router.Route(Integer(i))
//	}
//	router.Close()
type Router[T Hashable] struct {
	subscriptionIndex int
	parallelism       uint32
	routes            []chan T
	wg                *sync.WaitGroup
	closed            bool
}

func NewRouter[T Hashable](parallelism uint32) *Router[T] {
	routes := make([]chan T, parallelism)

	for i := 0; i < int(parallelism); i++ {
		routes[i] = make(chan T)
	}

	return &Router[T]{
		subscriptionIndex: 0,
		parallelism:       parallelism,
		routes:            routes,
		wg:                &sync.WaitGroup{},
	}
}

// Routes the input to the designated route
func (router *Router[T]) Route(in T) {
	if router.closed {
		panic("sending to a closed router")
	}
	index := in.Hash() % router.parallelism
	router.routes[index] <- in
}

// Subscribe a callback for one of the routes.
// Routes are subscribed in round-robin fashion
func (router *Router[T]) Subscribe(callback func(in T)) {
	if router.subscriptionIndex == int(router.parallelism) {
		router.subscriptionIndex = 0
	}
	routerIndex := router.subscriptionIndex
	router.wg.Add(1)
	go func() {
		for input := range router.routes[routerIndex] {
			callback(input)
		}
		router.wg.Done()
	}()
	router.subscriptionIndex += 1
}

// Closes all routes. Also waits for all subscriptions to end
func (router *Router[T]) Close() {
	if router.closed {
		panic("closing an already closed router")
	}
	router.closed = true
	for _, route := range router.routes {
		close(route)
	}
	router.wg.Wait()
}
