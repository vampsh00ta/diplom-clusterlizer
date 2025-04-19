package utils

import "errors"

var ErrMissing = errors.New("optional has no value")

const mustErrorMsg = "requested must-have optional value but its missing"

type Optional[T any] struct {
	Value T
	Valid bool
}

func NewOptional[T any](v T) Optional[T] {
	return Optional[T]{Value: v, Valid: true}
}

func NewEmptyOptional[T any]() Optional[T] {
	var empty T

	return Optional[T]{Value: empty, Valid: false}
}

func (o Optional[T]) From(v T) Optional[T] {
	var res Optional[T]
	res.Set(v)

	return res
}

func (o *Optional[T]) Set(v T) {
	o.Value = v
	o.Valid = true
}

func (o Optional[T]) Get() (T, error) {
	if !o.Valid {
		return *new(T), ErrMissing
	}

	return o.Value, nil
}

func OptionalFromPointer[T any](v *T) Optional[T] {
	if v == nil {
		return NewEmptyOptional[T]()
	}
	return NewOptional(*v)
}
