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

	"github.com/anyvoxel/airmid/anvil"
	"github.com/anyvoxel/airmid/anvil/xerrors"
	"github.com/anyvoxel/airmid/ioc"
)

type objectListenerInvoker struct {
	obj      ApplicationListener
	objAsync AsyncApplicationListener

	app Application
}

func (i *objectListenerInvoker) Invoke(ctx context.Context, event ApplicationEvent) {
	if i.obj != nil {
		i.obj.OnApplicationEvent(ctx, event)
	}

	if i.objAsync != nil {
		anvil.Must(i.app.Submit(func() {
			//nolint:contextcheck
			anvil.SafeRun(ctx, func(ctx context.Context) {
				i.objAsync.OnApplicationEventAsync(ctx, event)
			})
		}))
	}
}

// NewObjectListenerInvoker will return the object ListenInvoker impl.
func NewObjectListenerInvoker(obj any, app Application) (ListenerInvoker, error) {
	o := &objectListenerInvoker{
		app: app,
	}

	if vobj := ioc.IndirectTo[ApplicationListener](obj); vobj != nil {
		o.obj = vobj
	}
	if vobj := ioc.IndirectTo[AsyncApplicationListener](obj); vobj != nil {
		o.objAsync = vobj
	}

	if o.obj != nil || o.objAsync != nil {
		return o, nil
	}
	return nil, xerrors.WrapContinue("type '%T' of object doesn't implement Listener", obj)
}
