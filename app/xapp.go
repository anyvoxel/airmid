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
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"time"

	"github.com/anyvoxel/airmid/anvil/xerrors"
	"github.com/anyvoxel/airmid/ioc"
	"github.com/anyvoxel/airmid/ioc/props"
	slogctx "github.com/veqryn/slog-context"
)

// Application is the interface for app.
type Application interface {
	// Run will start the application, it will:
	// 1. prepare all of bean definitions used by application
	// 2. load the properties from env、flag、configfile
	// 3. initialize application's component, such as gopool、logger、metrics
	// 4. start all AppRunner
	// 5. Waiting for shutdown signals
	Run(ctx context.Context, opts ...Option) error

	// Shutdown will stop the application, the application will start graceful shutdown progress:
	// 1. invoke all AppRunner
	// 2. Waiting until shutdown.duration or all AppRunner exited
	// 3. exit the application
	Shutdown()

	// Submit will async run task in another goroutine
	Submit(task func()) error

	ApplicationEventPublisher
	ioc.BeanFactory
}

type airmidApplication struct {
	startupHandlers []ApplicationStartupHandler
	exitChan        chan struct{}

	ioc.BeanFactory
	listenerInvoker []ListenerInvoker

	// Use another struct to store the props, so we can autowire it
	props *airmidApplicationProps

	// TODO: change this to bootstrap config?
	appConfig       *config
	shutdownManager ShutdownManager

	gpool interface {
		Submit(func()) error
	}
}

type airmidApplicationProps struct {
	// NOTE: we must autowire it by concrete type
	runnerCompositor *RunnerCompositor `airmid:"autowire:?"`

	shutdownDuration time.Duration `airmid:"value:${airmid.shutdown.duration:=30s}"`
}

// NewApplication return the application.
func NewApplication() Application {
	v := &airmidApplication{
		startupHandlers: []ApplicationStartupHandler{
			&loggerStartupHandler{},
			&metricsStartupHandler{},
			&gpoolStartupHandler{},
			&shutdownStartupHandler{},
		},
		exitChan:        make(chan struct{}),
		BeanFactory:     ioc.NewBeanFactory(),
		listenerInvoker: make([]ListenerInvoker, 0),
		props:           &airmidApplicationProps{},
		appConfig:       &config{},
	}

	return v
}

func (a *airmidApplication) registerAppBeanDefinitions() error {
	beanDefinitions := map[string]ioc.BeanDefinition{
		"airmid.app.config": ioc.MustNewBeanDefinition(
			reflect.TypeOf((*config)(nil)),
			ioc.WithLazyMode(),
		),
		"airmid.app.resource.locator": ioc.MustNewBeanDefinition(
			reflect.TypeOf((*localResourceLocator)(nil)),
			ioc.WithLazyMode(),
		),
		"airmid.app.props": ioc.MustNewBeanDefinition(
			reflect.TypeOf((*airmidApplicationProps)(nil)),
			ioc.WithLazyMode(),
		),
		"airmid.app.runner.compositor": ioc.MustNewBeanDefinition(
			reflect.TypeOf((*RunnerCompositor)(nil)),
			ioc.WithLazyMode(),
		),
	}

	for name, beanDefinition := range beanDefinitions {
		err := a.RegisterBeanDefinition(name, beanDefinition)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *airmidApplication) runBeforeLoadProps(ctx context.Context, opt *option) error {
	for _, h := range a.startupHandlers {
		err := h.BeforeLoadProps(ctx, a, opt)
		if err != nil {
			return xerrors.Wrapf(err, "Handler [%s] BeforeLoadProps failed", h.Name())
		}
	}

	return nil
}

func (a *airmidApplication) runAfterLoadProps(ctx context.Context, opt *option) error {
	for _, h := range a.startupHandlers {
		err := h.AfterLoadProps(ctx, a, opt)
		if err != nil {
			return xerrors.Wrapf(err, "Handler [%s] AfterLoadProps failed", h.Name())
		}
	}

	return nil
}

func (a *airmidApplication) runBeforeStartRunner(ctx context.Context, opt *option) error {
	for _, h := range a.startupHandlers {
		err := h.BeforeStartRunner(ctx, a, opt)
		if err != nil {
			return xerrors.Wrapf(err, "Handler [%s] BeforeStartRunner failed", h.Name())
		}
	}

	return nil
}

func (a *airmidApplication) loadProperties(ctx context.Context) (err error) {
	// First, we load the env & flags property, so user can
	// configuration some application's setting.
	if err := a.loadPropsFromEnvAndFlags(ctx, a); err != nil {
		return err
	}

	a.appConfig, err = ioc.GetBean[*config](ctx, a, "airmid.app.config")
	if err != nil {
		return err
	}

	// Second, we load the config file's property, base on
	// env & flags property.
	err = a.appConfig.loadProperty(ctx, a)
	if err != nil {
		return err
	}

	// At last, we reload env & flags property, so env & flags
	// property can overwrite config file's setting, and take
	// high priority.
	return a.loadPropsFromEnvAndFlags(ctx, a)
}

func (a *airmidApplication) Run(ctx context.Context, opts ...Option) (err error) {
	opt := newOption(opts)

	err = a.runBeforeLoadProps(ctx, opt)
	if err != nil {
		return err
	}

	// TODO: optimize this
	if err := a.registerAppBeanDefinitions(); err != nil {
		return err
	}

	appRunnerCompoistorProcessor := &RunnerCompoistorProcessor{
		appRunnerNames: map[Runner]string{},
	}
	a.AddBeanPostProcessor(&applicationListenerDetector{app: a, singletonNames: map[string]bool{}})
	a.AddBeanPostProcessor(&ApplicationAwareProcessor{app: a})
	a.AddBeanPostProcessor(appRunnerCompoistorProcessor)

	err = a.loadProperties(ctx)
	if err != nil {
		return err
	}

	err = a.runAfterLoadProps(ctx, opt)
	if err != nil {
		return err
	}

	a.props, err = ioc.GetBean[*airmidApplicationProps](ctx, a, "airmid.app.props")
	if err != nil {
		return err
	}

	err = a.runBeforeStartRunner(ctx, opt)
	if err != nil {
		return err
	}

	err = a.PreInstantiateSingletons(ctx)
	if err != nil {
		return err
	}

	a.props.runnerCompositor.appRunnerNames = appRunnerCompoistorProcessor.appRunnerNames
	a.props.runnerCompositor.Run(ctx)

	slogctx.FromCtx(ctx).InfoContext(
		ctx,
		"all AppRunner started, wait for exit signal",
	)
	<-a.exitChan
	return nil
}

func (*airmidApplication) loadPropsFromEnvAndFlags(ctx context.Context, p props.Properties) error {
	epl := NewEnvPropertiesLoader("AIRMID_", DefaultEnvKeyConvertFunc)
	apl := NewOptionArgsPropertiesLoader()
	loaders := NewPropertiesLoaders()
	return loaders.Add(epl).Add(apl).LoadProperties(ctx, p)
}

func (a *airmidApplication) Shutdown() {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	fnName := "<unknown>"
	if fn != nil {
		fnName = fn.Name()
	}
	caller := fmt.Sprintf("%s:%v %s", file, line, fnName)

	a.shutdownManager.Shutdown(fmt.Sprintf("Shutdown from %v", caller))
}

func (a *airmidApplication) shutdownWithMessage(msg string) {
	a.Destroy()
	ctx := context.Background()

	slogctx.FromCtx(ctx).InfoContext(
		ctx,
		fmt.Sprintf("Application will shutdown within: %s", msg),
		slog.Any("ShutdownDuration", a.props.shutdownDuration),
	)

	ctx, cancel := context.WithTimeout(context.Background(), a.props.shutdownDuration)
	defer cancel()

	err := a.props.runnerCompositor.Stop(ctx)
	if xerrors.Is(err, context.DeadlineExceeded) {
		slogctx.FromCtx(ctx).InfoContext(
			ctx,
			"Application shutdown duration was too long, force close...",
			slog.Any("ShutdownDuration", a.props.shutdownDuration),
		)
	} else {
		slogctx.FromCtx(ctx).InfoContext(
			ctx,
			"All runners were stopped, application gracefully shutdown",
		)
	}

	close(a.exitChan)
}

func (a *airmidApplication) PublishEvent(ctx context.Context, event ApplicationEvent) {
	for _, invoker := range a.listenerInvoker {
		invoker.Invoke(ctx, event)
	}
}

func (a *airmidApplication) Submit(task func()) error {
	return a.gpool.Submit(task)
}
