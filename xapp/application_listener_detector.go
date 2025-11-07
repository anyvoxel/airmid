package xapp

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/anyvoxel/airmid/beans"
	"github.com/anyvoxel/airmid/logger"
	"github.com/anyvoxel/airmid/xerrors"
)

type applicationListenerDetector struct {
	// TODO: change to interface
	app *airmidApplication

	singletonNames map[string]bool
}

func (l *applicationListenerDetector) PostProcessBeanDefinition(beanName string, beanDefinition beans.BeanDefinition) {
	l.singletonNames[beanName] = beanDefinition.Scope() == beans.ScopeSingleton
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
