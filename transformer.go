package gocommons

import (
	"sync"
)

// Parallely applies a transformation to input slice
// and emits the transformed objects as a slice.
// Since these transformations are applied parallely
// order is not guarenteed
type Transformer[T any, R any] struct {
	parallelism     uint
	inArr           []T
	in              chan T
	out             chan Result[R]
	transformWg     *sync.WaitGroup
	resultContainer *Container[R]
	consumerWg      *sync.WaitGroup
	supplierWg      *sync.WaitGroup
}

func NewTransformer[T any, R any](parallelism uint, inArr []T, transformFunc func(T) (R, error)) *Transformer[T, R] {
	in := make(chan T)
	out := make(chan Result[R])
	var supplierWg sync.WaitGroup
	var transformWg sync.WaitGroup
	var collectWg sync.WaitGroup

	for i := 0; i < int(parallelism); i++ {
		transformWg.Add(1)
		go func() {
			for t := range in {
				result, err := transformFunc(t)
				out <- Result[R]{result, err}
			}
			transformWg.Done()
		}()
	}
	resultContainer := &Container[R]{
		arr: make([]Result[R], 0),
	}
	return &Transformer[T, R]{
		parallelism:     parallelism,
		inArr:           inArr,
		in:              in,
		out:             out,
		resultContainer: resultContainer,
		supplierWg:      &supplierWg,
		transformWg:     &transformWg,
		consumerWg:      &collectWg,
	}
}

type Result[R any] struct {
	ReturnVal R
	Err       error
}

type TransformJob[T any, R any] struct {
	resultContainer *Container[R]
	transformWg     *sync.WaitGroup
	consumerWg      *sync.WaitGroup
	supplierWg      *sync.WaitGroup
	in              chan T
	out             chan Result[R]
}

// Blocks until all transformation are completed and returns the results
// Errors are available using GetErrors
func (job *TransformJob[T, R]) Get() []Result[R] {
	job.supplierWg.Wait()
	job.transformWg.Wait()
	close(job.out)
	job.consumerWg.Wait()
	return job.resultContainer.arr
}

type Container[R any] struct {
	arr []Result[R]
}

// Applies the tranform func for each element in input slice.
func (transformer *Transformer[T, R]) Transform(filter func(Result[R]) bool) *TransformJob[T, R] {
	transformer.consumerWg.Add(1)
	go func() {
		for outElement := range transformer.out {
			if filter == nil || filter(outElement) {
				transformer.resultContainer.arr = append(transformer.resultContainer.arr, outElement)
			}
		}
		transformer.consumerWg.Done()
	}()
	transformer.supplierWg.Add(1)
	go func() {
		for _, element := range transformer.inArr {
			transformer.in <- element
		}
		close(transformer.in)
		transformer.supplierWg.Done()
	}()

	return &TransformJob[T, R]{
		resultContainer: transformer.resultContainer,
		transformWg:     transformer.transformWg,
		consumerWg:      transformer.consumerWg,
		supplierWg:      transformer.supplierWg,
		in:              transformer.in,
		out:             transformer.out,
	}
}
