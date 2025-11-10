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
	"reflect"
	"strings"

	"github.com/anyvoxel/airmid/anvil/xerrors"
)

// ConstructorBuilder will build constructor.
type ConstructorBuilder interface {
	// Build will return the Constructor for type
	Build(typ reflect.Type, args []ConstructorArgument) (Constructor, error)
}

// DefaultConstructorBuilder will build constructor with default behavior.
//  1. If the type have method:
//     1.1 The receiver is pointer to struct
//     1.2 The method name is NewXXX
//     1.3 The method must have no input argument
//     1.3 The method have return value.
//     If the out num is 1, it must be *XXX
//     If the out num is 2, first must be *XXX, second must be error
type DefaultConstructorBuilder struct {
	filters MethodFilters
}

// NewDefaultConstructorBuilder return the ConstructorBuilder impl.
func NewDefaultConstructorBuilder() ConstructorBuilder {
	return &DefaultConstructorBuilder{
		filters: []MethodFilter{
			&FnMethodFilter{
				FnFilter: MethodNameFilter,
			},
			&FnMethodFilter{
				FnFilter: MethodInputNumFilter,
			},
			&FnMethodFilter{
				FnFilter: MethodOutputNumFilter,
			},
			&FnMethodFilter{
				FnFilter: MethodOutput0Filter,
			},
			&FnMethodFilter{
				FnFilter: MethodOutput1Filter,
			},
		},
	}
}

// Build return the constructor from typ.
func (b *DefaultConstructorBuilder) Build(typ reflect.Type, args []ConstructorArgument) (Constructor, error) {
	methods := []reflect.Method{}
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		if b.filters.Filter(typ, method, args) {
			methods = append(methods, method)
		}
	}

	switch len(methods) {
	case 0:
		return &ReflectConstructor{
			typ: typ,
		}, nil
	case 1:
		return NewMethodConstructor(typ, methods[0], args)
	default:
	}

	return nil, xerrors.Errorf(
		"Too many function '%d' match the constructor spec with type '%s'",
		len(methods), typ.String())
}

// MethodNameFilter will check the method name.
func MethodNameFilter(typ reflect.Type, method reflect.Method, _ []ConstructorArgument) bool {
	if !strings.HasPrefix(method.Name, "New") {
		return false
	}

	expectMethodName := "New" + typ.Elem().Name()
	return strings.EqualFold(method.Name, expectMethodName)
}

// MethodInputNumFilter will check the method input num.
func MethodInputNumFilter(_ reflect.Type, method reflect.Method, args []ConstructorArgument) bool {
	return method.Func.Type().NumIn() == len(args)+1
}

// MethodOutputNumFilter will check the method output num.
func MethodOutputNumFilter(_ reflect.Type, method reflect.Method, _ []ConstructorArgument) bool {
	n := method.Func.Type().NumOut()
	return n == 1 || n == 2
}

// MethodOutput0Filter will check the method output0 type.
func MethodOutput0Filter(typ reflect.Type, method reflect.Method, _ []ConstructorArgument) bool {
	out0 := method.Func.Type().Out(0)
	return typ == out0
}

// MethodOutput1Filter will check the method output1 type.
func MethodOutput1Filter(_ reflect.Type, method reflect.Method, _ []ConstructorArgument) bool {
	n := method.Func.Type().NumOut()
	if n != 2 {
		return true
	}

	out1 := method.Func.Type().Out(1)
	return out1 == reflect.TypeOf((*error)(nil)).Elem()
}

// MethodFilter is the filter to match method with constructor spec.
type MethodFilter interface {
	// Filter will check the method on typ, it return true if the method can be constructor
	Filter(typ reflect.Type, method reflect.Method, args []ConstructorArgument) bool
}

// MethodFilters is the list of method filter.
type MethodFilters []MethodFilter

// Filter will call all filters in order, it return true if all filters passed.
func (m MethodFilters) Filter(typ reflect.Type, method reflect.Method, args []ConstructorArgument) bool {
	for _, fl := range m {
		if !fl.Filter(typ, method, args) {
			return false
		}
	}

	return true
}

// FnMethodFilter will proxy the Filter call to FnFilter.
type FnMethodFilter struct {
	FnFilter func(reflect.Type, reflect.Method, []ConstructorArgument) bool
}

// Filter implement the MethodFilter.Filter.
func (f *FnMethodFilter) Filter(typ reflect.Type, method reflect.Method, args []ConstructorArgument) bool {
	return f.FnFilter(typ, method, args)
}

var (
	_ ConstructorBuilder = (*DefaultConstructorBuilder)(nil)
)

var (
	constructorBuilder = NewDefaultConstructorBuilder()
)
