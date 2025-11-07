// Package utils contains the helper function/struct
package utils

// Option is an generic configuration helper.
type Option[T any] interface {
	// Apply applies this configuration to the given option
	Apply(*T)
}

// FnOption is an generic function helper for Option.
type FnOption[T any] struct {
	Fn func(*T)
}

// Apply will invoke fn on provide option.
func (o *FnOption[T]) Apply(opt *T) {
	o.Fn(opt)
}

// NewFnOption will instant an Option with f.
func NewFnOption[T any](f func(*T)) Option[T] {
	return &FnOption[T]{
		Fn: f,
	}
}
