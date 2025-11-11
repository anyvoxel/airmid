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
	"log/slog"
	"reflect"
	"unsafe"

	slogctx "github.com/veqryn/slog-context"
)

// PropertyValues is the values holder for bean property.
type PropertyValues interface {
	// AddValue will add the field value to container
	AddValue(fieldIndex int, value reflect.Value)

	// SetProperty will update the field value with property
	SetProperty(obj reflect.Value, fieldDescriptors []FieldDescriptor) error
}

type propertyValuesImpl struct {
	values map[int]reflect.Value
}

// NewPropertyValues return the PropertyValues impl.
func NewPropertyValues() PropertyValues {
	return &propertyValuesImpl{
		values: make(map[int]reflect.Value),
	}
}

func (p *propertyValuesImpl) AddValue(fieldIndex int, value reflect.Value) {
	p.values[fieldIndex] = value
}

func (p *propertyValuesImpl) SetProperty(obj reflect.Value, fieldDescriptors []FieldDescriptor) error {
	for _, fd := range fieldDescriptors {
		value, ok := p.values[fd.FieldIndex]
		if !ok {
			slogctx.FromCtx(context.TODO()).DebugContext(
				context.TODO(),
				"the value of field not found, skip it",
				slog.Int("FieldIndex", fd.FieldIndex),
				slog.String("FieldName", fd.Name),
			)
			continue
		}

		fv := obj.Field(fd.FieldIndex)
		if !fd.Unexported {
			fv.Set(value)
			continue
		}

		// The field is unexported (private field), we cannot directly use Set
		// It will panic with 'reflect.Value.Set using unaddressable value'
		reflect.NewAt(fv.Type(), unsafe.Pointer(fv.UnsafeAddr())).Elem().Set(value)
	}

	return nil
}
