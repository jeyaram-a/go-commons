package gocommons

import (
	"testing"
)

func contains(haystack []int, needle int) bool {
	for _, x := range haystack {
		if x == needle {
			return true
		}
	}
	return false
}

func TestTransformer(t *testing.T) {
	inArr := []int{1, 2, 3, 4, 5}
	transformer := NewTransformer[int, int](uint(4), inArr, func(x int) (int, error) {
		return x * 10, nil
	})
	out := transformer.Transform(nil)
	for _, x := range inArr {
		if !contains(out, x*10) {
			t.Errorf("%d missing ", (x * 10))
			t.Fail()
		}
	}

}

func TestTransformerWithFilter(t *testing.T) {
	inArr := []int{1, 2, 3, 4, 5}
	transformer := NewTransformer[int, int](uint(4), inArr, func(x int) (int, error) {
		return x * 10, nil
	})
	out := transformer.Transform(func(x Result[int]) bool {
		if x.Err != nil {
			return false
		}
		return x.ReturnVal%20 == 0
	})
	for _, x := range out {
		if x%20 != 0 {
			t.Errorf("%d missing ", (x * 10))
			t.Fail()
		}
	}

}

func TestTransformerWithEmpty(t *testing.T) {
	inArr := []int{}
	transformer := NewTransformer[int, int](uint(4), inArr, func(x int) (int, error) {
		return x * 10, nil
	})
	out := transformer.Transform(nil)
	for _, x := range inArr {
		if !contains(out, x*10) {
			t.Errorf("%d missing ", (x * 10))
			t.Fail()
		}
	}

}
