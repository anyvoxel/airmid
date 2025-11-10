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

// Package xreflect provide the helper function with reflect
package xreflect

import (
	"fmt"
	"reflect"

	"github.com/anyvoxel/airmid/anvil/xerrors"
)

// IndirectToInterface returns the interface's value
// If the interface is reflect.Value, return the reflect.Value.Interface().
func IndirectToInterface(v any) any {
	//nolint:gocritic,revive
	switch e := v.(type) {
	case reflect.Value:
		return e.Interface()
	}

	return v
}

// IndirectToValue returns the reflect.Value of v
// If the interface is non reflect.Value, return the reflect.ValueOf(v).
func IndirectToValue(v any) reflect.Value {
	switch e := v.(type) { //nolint:gocritic,revive
	case reflect.Value:
		return e
	}

	return reflect.ValueOf(v)
}

// IndirectToSetableValue returns the setable reflect.Value of i.
func IndirectToSetableValue(i any) (reflect.Value, error) {
	v := IndirectToValue(i)
	if v.Kind() == reflect.Ptr {
		// If the value is pointer and CanSet & IsNil, it maybe the pointer-type field in struct,
		// we must set it with the nil-non-pointer value.
		if v.CanSet() && v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		v = v.Elem()
	}

	if !v.CanSet() {
		return reflect.Value{}, xerrors.Errorf("The '%T' cannot been set, it must setable", i)
	}
	return v, nil
}

// NewValue return the value of type:
// 1. If type is Primitive type, it will return zero value
// 2. If type is Ptr type, it will return ptr to zero value
// 3. Other wise, it will return an error.
func NewValue(typ reflect.Type) (reflect.Value, error) {
	if typ.Kind() == reflect.Ptr {
		if typ.Elem().Kind() == reflect.Ptr {
			return reflect.Value{}, fmt.Errorf("Unsupport type '%v'", typ.String()) //nolint
		}

		ret := reflect.New(typ.Elem())
		// ret.Elem().Set(reflect.Zero(typ.Elem()))
		return ret, nil
	}
	return reflect.Zero(typ), nil
}

// IsPrimitiveType return nil error if type is primitive.
func IsPrimitiveType(typ reflect.Type) error {
	if v, ok := primitiveTypeMaps[typ.Kind()]; ok {
		if v {
			return nil
		}

		return fmt.Errorf("Kind '%v' is not primitive type", typ.Kind()) //nolint
	}

	return fmt.Errorf("Kind '%v' not found in primitive types", typ.Kind()) //nolint
}

var primitiveTypeMaps = map[reflect.Kind]bool{
	reflect.Invalid:       false,
	reflect.Bool:          true,
	reflect.Int:           true,
	reflect.Int8:          true,
	reflect.Int16:         true,
	reflect.Int32:         true,
	reflect.Int64:         true,
	reflect.Uint:          true,
	reflect.Uint8:         true,
	reflect.Uint16:        true,
	reflect.Uint32:        true,
	reflect.Uint64:        true,
	reflect.Uintptr:       true,
	reflect.Float32:       true,
	reflect.Float64:       true,
	reflect.Complex64:     true,
	reflect.Complex128:    true,
	reflect.Array:         false,
	reflect.Chan:          false,
	reflect.Func:          false,
	reflect.Interface:     false,
	reflect.Map:           false,
	reflect.Ptr:           true,
	reflect.Slice:         false,
	reflect.String:        true,
	reflect.Struct:        true,
	reflect.UnsafePointer: false,
}
