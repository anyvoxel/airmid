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

package app

// Option applies a configuration option value to a application.
type Option interface {
	apply(*option)
}

// option contains configuration options for a application.
type option struct {
	attrs []Attribute
}

// optionFunc applies a set of options to a option.
type optionFunc func(*option)

// apply the function with a option.
func (f optionFunc) apply(o *option) {
	f(o)
}

// Attribute holds a key and value pair.
type Attribute struct {
	Key   string
	Value string
}

// WithAttributes sets the attrs.
func WithAttributes(attrs ...Attribute) Option {
	return optionFunc(func(o *option) {
		o.attrs = append(o.attrs, attrs...)
	})
}

func newOption(options []Option) *option {
	o := &option{}
	for _, opt := range options {
		opt.apply(o)
	}

	return o
}
