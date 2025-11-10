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
	"strings"

	"github.com/anyvoxel/airmid/anvil"
	"github.com/anyvoxel/airmid/anvil/xerrors"
	"github.com/anyvoxel/airmid/anvil/xreflect"
)

type fnListenerInvoker struct {
	async bool
	name  string
	fn    reflect.Value
	arg1  reflect.Type

	app Application
}

// Invoke implement ListenerInvoker.Invoke.
func (i *fnListenerInvoker) Invoke(ctx context.Context, event ApplicationEvent) {
	in1 := reflect.ValueOf(event)
	if event == nil {
		// If the event is nil, we must construct it's zero value.
		// otherwise:
		//   1. in1.Type will panic
		//   2. reflect.Call will panic `reflect: Call using zero Value argument`
		// TODO: should we just panic?
		in1 = reflect.New(i.arg1).Elem()
	}

	if in1.Type().AssignableTo(i.arg1) {
		if !i.async {
			i.fn.Call([]reflect.Value{
				reflect.ValueOf(ctx),
				in1,
			})
			return
		}

		anvil.Must(i.app.Submit(func() {
			//nolint:contextcheck
			anvil.SafeRun(func() {
				i.fn.Call([]reflect.Value{
					reflect.ValueOf(ctx),
					in1,
				})
			})
		}))
	}
}

// NewFnListenerInvoker will return fn ListenerInvoker impl.
func NewFnListenerInvoker(fn reflect.Value, name string, app Application) (ListenerInvoker, error) {
	typ := fn.Type()

	if typ.Kind() != reflect.Func {
		return nil, xerrors.WrapContinue("type '%s' of object '%s' doesn't match kind func", typ.String(), name)
	}

	if typ.NumIn() != 2 || typ.NumOut() != 0 {
		return nil, xerrors.WrapContinue(
			"func '%s' expect in '%d' and out '%d', got in '%d' and out '%d'", name, 2, 0, typ.NumIn(), typ.NumOut(),
		)
	}

	arg0 := typ.In(0)
	arg1 := typ.In(1)
	if arg0 != xreflect.ContextType {
		return nil, xerrors.WrapContinue(
			"func '%s' arg0 type '%s', expect 'context.Context'", name, arg0.String())
	}

	if !arg1.Implements(applicationEventType) {
		return nil, xerrors.WrapContinue(
			"func '%s' arg1 type '%s', expect implement ApplicationEvent", name, arg1.String())
	}

	if arg1 == applicationEventType && (name == onApplicationEventFuncName || name == onApplicationEventAsyncFuncName) {
		return nil, xerrors.WrapContinue("func '%s' is object listener interface", name)
	}

	return &fnListenerInvoker{
		async: strings.HasSuffix(name, "Async"),
		name:  name,
		fn:    fn,
		arg1:  arg1,
		app:   app,
	}, nil
}
