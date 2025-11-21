package xiter

import "iter"

// Before runs the before function on each item before yielding it.
func Before[T any](iterator iter.Seq[T], before func(T)) iter.Seq[T] {
	return func(yield func(T) bool) {
		for item := range iterator {
			before(item)
			if !yield(item) {
				break
			}
		}
	}
}
