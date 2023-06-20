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
	transformer := NewTransformer[int, int](uint(4), inArr, func(x int) int {
		return x * 10
	})
	out := transformer.Transform()
	for _, x := range inArr {
		if !contains(out, x*10) {
			t.Errorf("%d missing ", (x * 10))
			t.Fail()
		}
	}

}

func TestTransformerWithEmpty(t *testing.T) {
	inArr := []int{}
	transformer := NewTransformer[int, int](uint(4), inArr, func(x int) int {
		return x * 10
	})
	out := transformer.Transform()
	for _, x := range inArr {
		if !contains(out, x*10) {
			t.Errorf("%d missing ", (x * 10))
			t.Fail()
		}
	}

}
