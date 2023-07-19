package gocommons

func GetIntArr(from, end int) []int {
	arr := make([]int, end-from)
	for i := from; i < end; i++ {
		arr[i-from] = i
	}
	return arr
}

func Sum(arr []int) int {
	sum := 0
	for _, x := range arr {
		sum += x
	}
	return sum
}

func Filter[T any](arr []T, filter func(T) bool) []T {
	results := make([]T, 0)
	for _, x := range arr {
		if filter == nil || filter(x) {
			results = append(results, x)
		}
	}
	return results
}

func Map_[T any, R any](in []T, mapFunc func(in T) R) []R {
	out := make([]R, 0)
	for _, x := range in {
		out = append(out, mapFunc(x))
	}

	return out
}

// all elements in other is in me
func Contains[T comparable](me []T, other []T) bool {
	for _, o := range other {
		present := false
		for _, m := range me {
			if o == m {
				present = true
				break
			}
		}
		if !present {
			return false
		}
	}
	return true
}
