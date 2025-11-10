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

	"github.com/anyvoxel/airmid/anvil/xreflect"
)

// Constructor is the constructor for Bean.
type Constructor interface {
	// NewObject will return the object from constructor
	NewObject(resolver ConstructorArgumentResolver) (reflect.Value, error)
}

// ReflectConstructor will build the object from reflect.
type ReflectConstructor struct {
	typ reflect.Type
}

// NewObject return the object from Reflect.
func (r *ReflectConstructor) NewObject(_ ConstructorArgumentResolver) (reflect.Value, error) {
	return xreflect.NewValue(r.typ)
}

// MethodConstructor will build the object from typ's method.
type MethodConstructor struct {
	typ    reflect.Type
	method reflect.Method

	args []ConstructorArgument
}

// NewMethodConstructor return the constructor for method.
func NewMethodConstructor(
	typ reflect.Type, method reflect.Method, args []ConstructorArgument,
) (*MethodConstructor, error) {
	for i := 0; i < len(args); i++ {
		// We must use +1 to indicate the argument's type, because the method first argument it (*self)
		args[i].Type = method.Type.In(i + 1)
	}

	return &MethodConstructor{
		typ:    typ,
		method: method,
		args:   args,
	}, nil
}

// NewObject return the object from Constructor.
func (m *MethodConstructor) NewObject(resolver ConstructorArgumentResolver) (reflect.Value, error) {
	ins, err := resolver.Resolve(m.args)
	if err != nil {
		return reflect.Value{}, err
	}

	argin := []reflect.Value{reflect.Zero(m.typ)}
	argin = append(argin, ins...)

	ret := m.method.Func.Call(argin)
	if len(ret) == 1 {
		return ret[0], nil
	}

	ret1 := ret[1].Interface()
	if ret1 == nil {
		return ret[0], nil
	}

	return ret[0], ret1.(error) //nolint:revive
}
