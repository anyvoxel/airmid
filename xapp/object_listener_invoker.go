package xapp

import (
	"context"

	"github.com/anyvoxel/airmid/beans"
	"github.com/anyvoxel/airmid/utils"
	"github.com/anyvoxel/airmid/xerrors"
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
		utils.Must(i.app.Submit(func() {
			//nolint:contextcheck
			utils.SafeRun(func() {
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

	if vobj := beans.IndirectTo[ApplicationListener](obj); vobj != nil {
		o.obj = vobj
	}
	if vobj := beans.IndirectTo[AsyncApplicationListener](obj); vobj != nil {
		o.objAsync = vobj
	}

	if o.obj != nil || o.objAsync != nil {
		return o, nil
	}
	return nil, xerrors.WrapContinue("type '%T' of object doesn't implement Listener", obj)
}
