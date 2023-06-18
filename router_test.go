package gocommons

import (
	"fmt"
	"testing"
)

type Integer int

func (x Integer) Hash() uint32 {
	return uint32(x % 5)
}

func Test_Router(t *testing.T) {
	parallelism := 5
	router := NewRouter[Integer](uint32(parallelism))
	count := []int{0, 0, 0, 0, 0}
	for i := 0; i < parallelism; i++ {
		j := i
		callback := func(x Integer) {
			fmt.Println("got ", x)
			count[j]++
		}
		router.Subscribe(callback)
	}

	fmt.Println("before sending")
	for i := 0; i < 10; i++ {
		router.Route(Integer(i))
	}
	router.Close()
	for i := 0; i < 5; i++ {

		if count[i] != 2 {
			t.Error("not as expected")
			t.Fail()
		}
	}

}
