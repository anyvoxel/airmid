// Package xerrors defines the helpful errors & function
package xerrors

import (
	"errors"

	pkgerrors "github.com/pkg/errors"
)

// Is alias the errors.Is.
var Is = errors.Is

// As alias the errors.As.
var As = errors.As

// New alias the errors.New.
var New = errors.New

// Wrapf alias the Wrapf.
var Wrapf = pkgerrors.Wrapf

// Errorf alias the Errorf.
var Errorf = pkgerrors.Errorf

var (
	// ErrNotImplement defines the function is not impelement.
	ErrNotImplement = errors.New("NotImplement")

	// ErrNotFound defines the object is not found.
	ErrNotFound = errors.New("ObjectNotFound")

	// ErrDuplicate defines the object is duplicate.
	ErrDuplicate = errors.New("DuplicateObject")

	// ErrContinue defines the err can continue.
	ErrContinue = errors.New("Continue")

	// ErrRetryable defines the err can retry.
	ErrRetryable = errors.New("Retryable")

	// ErrNonRetryable defines the err cann't retry.
	ErrNonRetryable = errors.New("NonRetryable")
)

// WrapNotFound return the wraped not found error.
func WrapNotFound(format string, args ...any) error {
	return Wrapf(ErrNotFound, format, args...)
}

// IsNotFound return true if the error is object not found.
func IsNotFound(err error) bool {
	return Is(err, ErrNotFound)
}

// WrapDuplicate return the wraped duplicate error.
func WrapDuplicate(format string, args ...any) error {
	return Wrapf(ErrDuplicate, format, args...)
}

// IsDuplicate return true if the error is duplicate.
func IsDuplicate(err error) bool {
	return Is(err, ErrDuplicate)
}

// WrapNotImplement return the not implement error.
func WrapNotImplement(format string, args ...any) error {
	return Wrapf(ErrNotImplement, format, args...)
}

// IsNotImplement return true if the error is not implement.
func IsNotImplement(err error) bool {
	return Is(err, ErrNotImplement)
}

// WrapContinue return the Continue error.
func WrapContinue(format string, args ...any) error {
	return Wrapf(ErrContinue, format, args...)
}

// IsContinue return true if the error is continue.
func IsContinue(err error) bool {
	return Is(err, ErrContinue)
}

// WrapRetryable return the Retryable error.
func WrapRetryable(format string, args ...any) error {
	return Wrapf(ErrRetryable, format, args...)
}

// IsRetryable return true if the error is Retryable.
func IsRetryable(err error) bool {
	return Is(err, ErrRetryable)
}

// WrapNonRetryable return the NonRetryable error.
func WrapNonRetryable(format string, args ...any) error {
	return Wrapf(ErrNonRetryable, format, args...)
}

// IsNonRetryable return true if the error is non-retryable.
func IsNonRetryable(err error) bool {
	return Is(err, ErrNonRetryable)
}
