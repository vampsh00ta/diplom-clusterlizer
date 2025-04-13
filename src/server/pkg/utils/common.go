package utils

import "time"

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for range attempts {
		if err = fn(); err != nil {
			time.Sleep(delay)
			continue
		}
		return nil
	}
	return nil
}

func PythonMap[T any, R any](input []T, fn func(T) R) []R {
	result := make([]R, len(input))
	for i, v := range input {
		result[i] = fn(v)
	}
	return result
}

func SafeNil[T any](t *T) T {
	var res T
	if t == nil {
		return res
	}
	return *t
}
