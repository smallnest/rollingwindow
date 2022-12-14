package rollingwindow

import (
	"golang.org/x/exp/constraints"
)

// Sum the values within the float64 window.
func Sum[T constraints.Integer | constraints.Float](w Window[T]) T {
	var result T
	for _, bucket := range w {
		for _, p := range bucket {
			result = result + p
		}
	}
	return result
}

// Avg the values within the float64 window.
func Avg[T constraints.Integer | constraints.Float](w Window[T]) T {
	var result T
	var count int
	for _, bucket := range w {
		for _, p := range bucket {
			result = result + p
			count = count + 1
		}
	}
	return result / T(count)
}

// Min the values within the float64 window.
func Min[T constraints.Integer | constraints.Float](w Window[T]) T {
	var result T
	var started = true
	for _, bucket := range w {
		for _, p := range bucket {
			if started {
				result = p
				started = false
				continue
			}
			if p < result {
				result = p
			}
		}
	}
	return result
}

// Max the values within the float64 window.
func Max[T constraints.Integer | constraints.Float](w Window[T]) T {
	var result T
	var started = true
	for _, bucket := range w {
		for _, p := range bucket {
			if started {
				result = p
				started = false
				continue
			}
			if p > result {
				result = p
			}
		}
	}
	return result
}
