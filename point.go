package rollingwindow

import "sync"

// PointPolicy is a rolling window policy that tracks the last N
// values inserted regardless of insertion time.
type PointPolicy[T any] struct {
	windowSize int
	window     Window[T]
	offset     int
	lock       *sync.RWMutex
}

// NewPointPolicy generates a Policy that operates on a rolling set of
// input points. The number of points is determined by the size of the given
// window. Each bucket will contain, at most, one data point when the window
// is full.
func NewPointPolicy[T any](window Window[T]) *PointPolicy[T] {
	var p = &PointPolicy[T]{
		windowSize: len(window),
		window:     window,
		lock:       &sync.RWMutex{},
	}
	for offset, bucket := range window {
		if len(bucket) < 1 {
			window[offset] = make([]T, 1)
		}
	}
	return p
}

// Window returns the current window.
func (w *PointPolicy[T]) Window() [][]T {
	return w.window
}

// Append a value to the window.
func (w *PointPolicy[T]) Append(value T) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.window[w.offset][0] = value
	w.offset = (w.offset + 1) % w.windowSize
}

// Reduce the window to a single value using a reduction function.
func (w *PointPolicy[T]) Reduce(f func(Window[T]) T) T {
	w.lock.Lock()
	defer w.lock.Unlock()

	return f(w.window)
}

// Count returns counts of values in the window.
func (w *PointPolicy[T]) Count() int {
	w.lock.Lock()
	defer w.lock.Unlock()

	var result int
	for _, bucket := range w.window {
		result += len(bucket)
	}

	return result
}
