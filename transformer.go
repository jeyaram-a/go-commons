package gocommons

import (
	"sync"
)

// Parallely applies a transformation to input slice
// and emits the transformed objects as a slice.
// Since these transformations are applied parallely
// order is not guarenteed
type Transformer[T any, R any] struct {
	parallelism uint
	inArr       []T
	in          chan T
	out         chan Result[R]
	wg          *sync.WaitGroup
}

func NewTransformer[T any, R any](parallelism uint, inArr []T, transformFunc func(T) (R, error)) *Transformer[T, R] {
	in := make(chan T)
	out := make(chan Result[R])
	var wg sync.WaitGroup

	for i := 0; i < int(parallelism); i++ {
		wg.Add(1)
		go func() {
			for t := range in {
				result, err := transformFunc(t)
				out <- Result[R]{result, err}
			}
			wg.Done()
		}()
	}

	return &Transformer[T, R]{
		parallelism: parallelism,
		inArr:       inArr,
		in:          in,
		out:         out,
		wg:          &wg,
	}
}

type Result[R any] struct {
	ReturnVal R
	Err       error
}

// Applies the tranform func for each element in input slice.
// Blocks until all transformation are completed
func (transformer *Transformer[T, R]) Transform(filter func(Result[R]) bool) []R {
	collected := make([]R, 0)
	var collectWg sync.WaitGroup
	collectWg.Add(1)
	go func() {
		for outElement := range transformer.out {
			if filter == nil || filter(outElement) {
				collected = append(collected, outElement.ReturnVal)
			}
		}
		collectWg.Done()
	}()

	for _, element := range transformer.inArr {
		transformer.in <- element
	}
	close(transformer.in)
	transformer.wg.Wait()
	close(transformer.out)
	collectWg.Wait()

	return collected
}
