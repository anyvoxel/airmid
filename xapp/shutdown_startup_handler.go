package xapp

import (
	"context"
	"os"
	"reflect"
	"syscall"

	"github.com/anyvoxel/airmid/beans"
)

type shutdownStartupHandlerConfigration struct {
	// shutdownSignals is the os.Signal we wait for graceful shutdown
	// Defaults: SIGHUP、SIGINT、SIGTERM
	shutdownSignals []int `airmid:"value:${airmid.shutdown.signals:=1,2,15}"`
}

type shutdownStartupHandler struct{}

var (
	_ ApplicationStartupHandler = (*shutdownStartupHandler)(nil)
)

func (*shutdownStartupHandler) Name() string {
	return "ShutdownStartupHandler"
}

func (*shutdownStartupHandler) BeforeLoadProps(_ context.Context, app *airmidApplication, _ *option) error {
	beanDefinitions := map[string]beans.BeanDefinition{
		"airmidShutdownStartupHandlerConfiguration": beans.MustNewBeanDefinition(
			reflect.TypeOf((*shutdownStartupHandlerConfigration)(nil)),
			beans.WithLazyMode(),
		),
	}

	for name, beanDefinition := range beanDefinitions {
		err := app.RegisterBeanDefinition(name, beanDefinition)
		if err != nil {
			return err
		}
	}

	return nil
}

func (*shutdownStartupHandler) AfterLoadProps(_ context.Context, _ *airmidApplication, _ *option) error {
	return nil
}

func (*shutdownStartupHandler) BeforeStartRunner(_ context.Context, app *airmidApplication, _ *option) error {
	shutdownC, err := beans.GetBean[*shutdownStartupHandlerConfigration](app, "airmidShutdownStartupHandlerConfiguration")
	if err != nil {
		return err
	}

	sigs := make([]os.Signal, len(shutdownC.shutdownSignals))

	for _, signal := range shutdownC.shutdownSignals {
		sigs = append(sigs, syscall.Signal(signal))
	}

	//nolint:contextcheck
	app.shutdownManager = NewSignalShutdownManager([]ShutdownHandler{
		app.shutdownWithMessage,
	}, sigs...)
	return nil
}
