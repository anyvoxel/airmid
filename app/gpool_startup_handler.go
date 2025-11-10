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

	"github.com/anyvoxel/airmid/ioc"
)

type gpoolStartupHandler struct{}

var (
	_ ApplicationStartupHandler = (*gpoolStartupHandler)(nil)
)

func (*gpoolStartupHandler) Name() string {
	return "GpoolStartupHander"
}

func (*gpoolStartupHandler) BeforeLoadProps(_ context.Context, app *airmidApplication, _ *option) error {
	return app.RegisterBeanDefinition(
		"airmidGPoolFactory",
		ioc.MustNewBeanDefinition(
			reflect.TypeOf((*GPoolFactory)(nil)),
			ioc.WithLazyMode(),
		),
	)
}

func (*gpoolStartupHandler) AfterLoadProps(_ context.Context, app *airmidApplication, _ *option) error {
	gpoolFactoryObject, err := ioc.GetBean[*GPoolFactory](app, "airmidGPoolFactory")
	if err != nil {
		return err
	}

	gpool, err := gpoolFactoryObject.CreatePool()
	if err != nil {
		return err
	}

	app.gpool = gpool
	return nil
}

func (*gpoolStartupHandler) BeforeStartRunner(_ context.Context, _ *airmidApplication, _ *option) error {
	return nil
}
