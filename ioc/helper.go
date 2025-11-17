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

package ioc

import (
	"context"
	"reflect"

	"github.com/anyvoxel/airmid/anvil/xerrors"
)

// IsTypeMatched return true if the beanTyp can convert to targetTyp (assigned).
func IsTypeMatched(targetTyp reflect.Type, beanTyp reflect.Type) bool {
	if targetTyp.Kind() == reflect.Interface {
		return beanTyp.Implements(targetTyp)
	}

	return beanTyp.AssignableTo(targetTyp)
}

// Proxyer is an wrapper for bean object.
type Proxyer interface {
	// OriginalObject return the wrapped bean object
	OriginalObject() any
}

// IndirectTo return the target object which implement specific interface.
func IndirectTo[T any](o any) T {
	var v T

	vv, ok := o.(T)
	if ok {
		return vv
	}

	if po, ok := o.(Proxyer); ok {
		return IndirectTo[T](po.OriginalObject())
	}

	return v
}

// GetBean return the target beanObject.
func GetBean[T any](ctx context.Context, f BeanFactory, beanName string) (T, error) {
	var v T

	beanObject, err := f.GetBean(ctx, beanName)
	if err != nil {
		return v, err
	}

	vv, ok := beanObject.(T)
	if !ok {
		return v, xerrors.Errorf("cannot convert %T to %T", beanObject, v)
	}

	return vv, nil
}
