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
	"reflect"
)

var (
	// applicationEventType is the reflect.Type of ApplicationEvent.
	applicationEventType = reflect.TypeOf((*ApplicationEvent)(nil)).Elem()

	// NOTE: when we compare with Method Type, we cannot distinct between different object.method
	//  on retrieve method from object.
	// If the Method Type is retrieve from reflect.Type, it will take the receiver as first argument.

	// onApplicationEventFuncType is the function reflect.Type of ApplicationListener.OnApplicationEvent.
	onApplicationEventFuncType = reflect.TypeOf((*ApplicationListener)(nil)).Elem().Method(0).Type

	// onApplicationEventFuncName is the function name of ApplicationListener.OnApplicationEvent.
	onApplicationEventFuncName = "OnApplicationEvent"

	// onApplicationEventAsyncFuncType is the function reflect.Type of AsyncApplicationListener.OnApplicationEventAsync.
	onApplicationEventAsyncFuncType = reflect.TypeOf((*AsyncApplicationListener)(nil)).Elem().Method(0).Type

	// onApplicationEventAsyncFuncName is the function name of AsyncApplicationListener.OnApplicationEventAsync.
	onApplicationEventAsyncFuncName = "OnApplicationEventAsync"
)

// ApplicationListener is to be implemented by user listener.
type ApplicationListener interface {
	OnApplicationEvent(ctx context.Context, event ApplicationEvent)
}

// AsyncApplicationListener is to be implemented by user async listener.
type AsyncApplicationListener interface {
	OnApplicationEventAsync(ctx context.Context, event ApplicationEvent)
}
