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
	"fmt"
	"reflect"
)

// ConstructorArgumentResolver to resolve arguments.
type ConstructorArgumentResolver interface {
	// Resolve will build the args value from definition
	Resolve(args []ConstructorArgument) ([]reflect.Value, error)
}

type factoryConstructorArgumentResolver struct {
	bf *beanFactoryImpl
}

func (r *factoryConstructorArgumentResolver) Resolve(args []ConstructorArgument) ([]reflect.Value, error) {
	stFields := make([]reflect.StructField, 0, len(args))
	fds := make([]FieldDescriptor, 0, len(args))
	for i, arg := range args {
		st := reflect.StructField{
			Name: fmt.Sprintf("F%d", i),
			Type: arg.Type,
		}

		stFields = append(stFields, st)

		switch {
		case arg.Property != nil:
			fds = append(fds, FieldDescriptor{
				FieldIndex: i,
				Name:       st.Name,
				Typ:        st.Type,
				Unexported: false,
				Property:   arg.Property,
			})
		case arg.Bean != nil:
			fds = append(fds, FieldDescriptor{
				FieldIndex: i,
				Name:       st.Name,
				Typ:        st.Type,
				Unexported: false,
				Bean:       arg.Bean,
			})
		default:
			// We do nothing for constant value
		}
	}
	beanObject := reflect.New(reflect.StructOf(stFields))
	err := r.bf.wireStruct(beanObject, fds)
	if err != nil {
		return nil, err
	}

	ret := make([]reflect.Value, 0, len(args))
	for i, arg := range args {
		switch {
		case arg.Bean != nil || arg.Property != nil:
			// We need get value from beanObject, it has been inject with wireStruct
			ret = append(ret, beanObject.Elem().Field(i))
		default:
			ret = append(ret, arg.Value)
		}
	}

	return ret, nil
}
