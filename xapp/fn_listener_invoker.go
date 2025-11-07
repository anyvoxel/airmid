package xapp

import (
	"context"
	"reflect"
	"strings"

	"github.com/anyvoxel/airmid/utils"
	"github.com/anyvoxel/airmid/utils/xreflect"
	"github.com/anyvoxel/airmid/xerrors"
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

		utils.Must(i.app.Submit(func() {
			//nolint:contextcheck
			utils.SafeRun(func() {
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
