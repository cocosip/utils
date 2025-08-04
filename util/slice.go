// Package util provides a collection of utility functions.
package util

import "slices"

// Distinct returns a new slice with all unique elements from the input slice.
// The order of elements in the result is the order of their first appearance in the input.
// It uses the standard library's slices.Contains function for checking existence.
func Distinct[T comparable](items []T) []T {
	var result []T
	for _, item := range items {
		if !slices.Contains(result, item) {
			result = append(result, item)
		}
	}
	return result
}

// FirstOrDefault returns the first element of a slice and true, or a zero value and false if the slice is empty.
func FirstOrDefault[T any](items []T) (T, bool) {
	var zero T
	if len(items) > 0 {
		return items[0], true
	}
	return zero, false
}

// LastOrDefault returns the last element of a slice and true, or a zero value and false if the slice is empty.
func LastOrDefault[T any](items []T) (T, bool) {
	var zero T
	if len(items) > 0 {
		return items[len(items)-1], true
	}
	return zero, false
}