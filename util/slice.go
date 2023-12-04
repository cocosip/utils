package util

type SliceComparableType[T any] comparable

func Contains[T SliceComparableType[T]](items []T, value T) bool {
	for i := range items {
		if items[i] == value {
			return true
		}
	}
	return false
}

func Distinct[T SliceComparableType[T]](items []T) []T {
	var r []T
	for i := range items {
		if !Contains(r, items[i]) {
			r = append(r, items[i])
		}
	}
	return r
}
