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
)

// BeanPriority is the priority of bean, more higher the priority order is, the index of a bean in the slice/array
// smaller is.
type BeanPriority int32

// BeanPriorityOrder define the priority order interface.
type BeanPriorityOrder interface {
	// GetPriority return the priority of a bean, more higher the priority order is,
	// the index of a bean in the slice/array smaller is.
	GetPriority() BeanPriority
}

// CandidateBeans is a slice type of reflect.Value.
type CandidateBeans []reflect.Value

func (c CandidateBeans) Len() int {
	return len(c)
}

// Less return true in following conditions:
// 1. c[j] implements BeanPriorityOrder, but c[i] not,
// 2. c[i] and c[j] implements BeanPriorityOrder, c[j]'s priority > c[i]'s priority.
func (c CandidateBeans) Less(i, j int) bool {
	vi := IndirectTo[BeanPriorityOrder](c[i].Interface())
	vj := IndirectTo[BeanPriorityOrder](c[j].Interface())

	if vi != nil && vj == nil {
		return true
	}
	if vi != nil && vj != nil {
		return vi.GetPriority() > vj.GetPriority()
	}

	return false
}

func (c CandidateBeans) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
