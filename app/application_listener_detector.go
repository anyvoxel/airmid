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
	"log/slog"
	"reflect"

	"github.com/anyvoxel/airmid/anvil/logger"
	"github.com/anyvoxel/airmid/anvil/xerrors"
	"github.com/anyvoxel/airmid/ioc"
)

type applicationListenerDetector struct {
	// TODO: change to interface
	app *airmidApplication

	singletonNames map[string]bool
}

func (l *applicationListenerDetector) PostProcessBeanDefinition(beanName string, beanDefinition ioc.BeanDefinition) {
	l.singletonNames[beanName] = beanDefinition.Scope() == ioc.ScopeSingleton
}

func (*applicationListenerDetector) PostProcessBeforeInitialization(obj any, _ string) (v any, err error) {
	return obj, nil
}

func (l *applicationListenerDetector) PostProcessAfterInitialization(obj any, beanName string) (v any, err error) {
	if !l.singletonNames[beanName] {
		return obj, nil
	}

	// TODO: only invoke on Singleton scope
	vv := reflect.ValueOf(obj)
	n := vv.Type().NumMethod()

	for i := 0; i < n; i++ {
		fn := vv.Method(i)
		m := vv.Type().Method(i)
		invoker, err := NewFnListenerInvoker(fn, m.Name, l.app)
		if err != nil {
			if xerrors.IsContinue(err) {
				logger.FromContext(context.TODO()).DebugContext(
					context.TODO(),
					"bean method doesn't implement the listener, skip it",
					slog.String("BeanName", beanName),
					slog.String("MethodName", m.Name),
					slog.Any("Error", err),
				)
				continue
			}

			return nil, err
		}

		l.app.listenerInvoker = append(l.app.listenerInvoker, invoker)
	}

	invoker, err := NewObjectListenerInvoker(obj, l.app)
	if err != nil {
		if !xerrors.IsContinue(err) {
			return nil, err
		}

		logger.FromContext(context.TODO()).DebugContext(
			context.TODO(),
			"bean doesn't implement the listener, skip it",
			slog.String("BeanName", beanName),
			slog.Any("Error", err),
		)
	} else {
		l.app.listenerInvoker = append(l.app.listenerInvoker, invoker)
	}

	return obj, nil
}
