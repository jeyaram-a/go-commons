package gocommons

import "sync"

// Parallely applies a transformation to input slice
// and emits the transformed objects as a slice.
// Since these transformations are applied parallely
// order is not guarenteed
type Transformer[T any, R any] struct {
	parallelism uint
	inArr       []T
	in          chan T
	out         chan R
	wg          *sync.WaitGroup
}

func NewTransformer[T any, R any](parallelism uint, inArr []T, transformFunc func(T) (R, error)) *Transformer[T, R] {
	in := make(chan T)
	out := make(chan R)
	var wg sync.WaitGroup

	for i := 0; i < int(parallelism); i++ {
		wg.Add(1)
		go func() {
			for t := range in {
				result, err := transformFunc(t)
				if err == nil {
					out <- result
				}
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

// Applies the tranform func for each element in input slice.
// Blocks until all transformation are completed
func (transformer *Transformer[T, R]) Transform(filter func(R) bool) []R {
	collected := make([]R, 0)

	go func() {
		for outElement := range transformer.out {
			if filter == nil || filter(outElement) {
				collected = append(collected, outElement)
			}
		}

	}()

	for _, element := range transformer.inArr {
		transformer.in <- element
	}

	close(transformer.in)
	transformer.wg.Wait()
	close(transformer.out)

	return collected
}
