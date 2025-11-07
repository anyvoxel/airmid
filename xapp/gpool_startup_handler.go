package xapp

import (
	"context"
	"reflect"

	"github.com/anyvoxel/airmid/beans"
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
		beans.MustNewBeanDefinition(
			reflect.TypeOf((*GPoolFactory)(nil)),
			beans.WithLazyMode(),
		),
	)
}

func (*gpoolStartupHandler) AfterLoadProps(_ context.Context, app *airmidApplication, _ *option) error {
	gpoolFactoryObject, err := beans.GetBean[*GPoolFactory](app, "airmidGPoolFactory")
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
