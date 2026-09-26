package common

func Ptr[T any](value T) *T {
	return &value
}

func ValueByPtr[T any](ptr *T) T {
	var value T
	if ptr != nil {
		value = *ptr
	}

	return value
}
