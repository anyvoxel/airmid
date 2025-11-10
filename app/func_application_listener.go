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

import (
	"context"
)

// FuncApplicationListener is the handler for compose function.
// NOTE: this is used for unit test, be careful to use it.
type FuncApplicationListener struct {
	OnEventFunc      func(ctx context.Context, event ApplicationEvent)
	OnEventAsyncFunc func(ctx context.Context, event ApplicationEvent)
}

// OnApplicationEvent implement ApplicationListener.OnApplicationEvent.
func (l *FuncApplicationListener) OnApplicationEvent(ctx context.Context, event ApplicationEvent) {
	if l.OnEventFunc != nil {
		l.OnEventFunc(ctx, event)
	}
}

// OnApplicationEventAsync implement AsyncApplicationListener.OnApplicationEventAsync.
func (l *FuncApplicationListener) OnApplicationEventAsync(ctx context.Context, event ApplicationEvent) {
	if l.OnEventAsyncFunc != nil {
		l.OnEventAsyncFunc(ctx, event)
	}
}

var (
	_ ApplicationListener      = &FuncApplicationListener{}
	_ AsyncApplicationListener = &FuncApplicationListener{}
)
