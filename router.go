package gocommons

import (
	"sync"
)

type Router[T Hashable] struct {
	subscriptionIndex int
	parallelism       uint32
	routes            []chan T
	wg                *sync.WaitGroup
}

type Hashable interface {
	Hash() uint32
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

func (router *Router[T]) Route(in T) {
	index := in.Hash() % router.parallelism
	router.routes[index] <- in
}

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

func (router *Router[T]) Close() {
	for _, route := range router.routes {
		close(route)
	}
	router.wg.Wait()
}
