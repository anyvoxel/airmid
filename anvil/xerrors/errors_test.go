// Copyright (c) 2025 The anyvoxel Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package xerrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/onsi/gomega"
)

type MyCustomError struct{ Msg string }

func (e *MyCustomError) Error() string { return e.Msg }

func TestTypedError(t *testing.T) {
	g := gomega.NewWithT(t)

	t.Run("CreationAndWrapping", func(t *testing.T) {
		cause := errors.New("underlying cause")
		err := NewTyped(NotFound{}).WithCause(cause).WithMessage("resource not found")

		g.Expect(err).To(gomega.HaveOccurred())
		g.Expect(err.Message).To(gomega.Equal("resource not found"))
		g.Expect(Is(err, cause)).To(gomega.BeTrue())
		g.Expect(err.Error()).To(gomega.Equal("NotFound: resource not found (cause: underlying cause)"))
	})

	t.Run("AsSuccess", func(t *testing.T) {
		err := NewTyped(NotFound{}).WithMessage("test")
		var notFoundErr *TypedError[NotFound]
		g.Expect(As(err, &notFoundErr)).To(gomega.BeTrue())
		g.Expect(notFoundErr).NotTo(gomega.BeNil())
	})

	t.Run("AsFailureForDifferentTypedError", func(t *testing.T) {
		err := NewTyped(NotFound{}).WithMessage("test")
		var duplicateErr *TypedError[Duplicate]
		g.Expect(As(err, &duplicateErr)).To(gomega.BeFalse())
		g.Expect(duplicateErr).To(gomega.BeNil())
	})

	t.Run("AsFailureForNil", func(t *testing.T) {
		var notFoundErr *TypedError[NotFound]
		g.Expect(As(nil, &notFoundErr)).To(gomega.BeFalse())
	})

	t.Run("AsFailureForOtherError", func(t *testing.T) {
		err := errors.New("some other error")
		var notFoundErr *TypedError[NotFound]
		g.Expect(As(err, &notFoundErr)).To(gomega.BeFalse())
	})

	t.Run("TableDrivenForSemanticTypes", func(t *testing.T) {
		testCases := []struct {
			name      string
			err       error
			checkFunc func(err error) bool
		}{
			{
				name: "NotImplement",
				err:  NewTyped(NotImplement{}),
				checkFunc: func(err error) bool {
					var target *TypedError[NotImplement]
					return As(err, &target)
				},
			},
			{
				name: "NotFound",
				err:  NewTyped(NotFound{}),
				checkFunc: func(err error) bool {
					var target *TypedError[NotFound]
					return As(err, &target)
				},
			},
			{
				name: "Duplicate",
				err:  NewTyped(Duplicate{}),
				checkFunc: func(err error) bool {
					var target *TypedError[Duplicate]
					return As(err, &target)
				},
			},
			{
				name: "Continue",
				err:  NewTyped(Continue{}),
				checkFunc: func(err error) bool {
					var target *TypedError[Continue]
					return As(err, &target)
				},
			},
			{
				name: "Retryable",
				err:  NewTyped(Retryable{}),
				checkFunc: func(err error) bool {
					var target *TypedError[Retryable]
					return As(err, &target)
				},
			},
			{
				name: "NonRetryable",
				err:  NewTyped(NonRetryable{}),
				checkFunc: func(err error) bool {
					var target *TypedError[NonRetryable]
					return As(err, &target)
				},
			},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("PositiveCheck_%s", tc.name), func(t *testing.T) {
				g := gomega.NewWithT(t)
				g.Expect(tc.checkFunc(tc.err)).To(gomega.BeTrue())
			})
			t.Run(fmt.Sprintf("NegativeCheck_%s", tc.name), func(t *testing.T) {
				g := gomega.NewWithT(t)
				g.Expect(tc.checkFunc(errors.New("different error"))).To(gomega.BeFalse())
			})
		}
	})

	t.Run("AsType", func(t *testing.T) {
		err := NewTyped(NotFound{}).WithMessage("test")

		// Successful match - getting *TypedError[NotFound]
		typedNotFoundErr, ok := AsType[*TypedError[NotFound]](err)
		g.Expect(ok).To(gomega.BeTrue())
		g.Expect(typedNotFoundErr).NotTo(gomega.BeNil())
		g.Expect(typedNotFoundErr.Context).To(gomega.Equal(NotFound{})) // Check the context
		g.Expect(typedNotFoundErr.Message).To(gomega.Equal("test"))     // Check the message

		// Failed match - getting *TypedError[Duplicate]
		typedDuplicateErr, ok := AsType[*TypedError[Duplicate]](err)
		g.Expect(ok).To(gomega.BeFalse())
		g.Expect(typedDuplicateErr).To(gomega.BeNil())

		// Test with an error chain involving a custom error
		customWrappedErr := NewTyped(InvalidArgument{}).WithCause(&MyCustomError{Msg: "specific custom error"})

		// Successfully extract MyCustomError
		myCustomErr, ok := AsType[*MyCustomError](customWrappedErr)
		g.Expect(ok).To(gomega.BeTrue())
		g.Expect(myCustomErr).NotTo(gomega.BeNil())
		g.Expect(myCustomErr.Msg).To(gomega.Equal("specific custom error"))

		// Failed to extract MyCustomError from a different error
		otherErr := errors.New("just another error")
		otherMyCustomErr, ok := AsType[*MyCustomError](otherErr)
		g.Expect(ok).To(gomega.BeFalse())
		g.Expect(otherMyCustomErr).To(gomega.BeNil())
	})
}
