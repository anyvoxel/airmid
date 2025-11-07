package xapp

import (
	"context"
	"log/slog"

	"github.com/anyvoxel/airmid/beans"
	"github.com/anyvoxel/airmid/logger"
	"github.com/anyvoxel/airmid/utils/parallel"
)

// AppRunner is the interface for runner.
type AppRunner interface {
	// Run will been called when the application is ready.
	Run(ctx context.Context)

	// Stop will be called when the application is shutting down.
	Stop(ctx context.Context)
}

// AppRunnerCompoistorProcessor is the postprocessor for AppRunner.
type AppRunnerCompoistorProcessor struct {
	appRunnerNames map[AppRunner]string
}

// PostProcessBeforeInitialization implement the BeanPostProcessor.PostProcessBeforeInitialization.
func (p *AppRunnerCompoistorProcessor) PostProcessBeforeInitialization(obj any, beanName string) (v any, err error) {
	if vobj := beans.IndirectTo[AppRunner](obj); vobj != nil {
		// NOTE: we cannot return a WrapObject for AppRunner, because
		// when we use embedded interface, the WrapObject didn't implement
		// those interface which implemented by concrete type
		p.appRunnerNames[vobj] = beanName
	}

	return obj, nil
}

// PostProcessAfterInitialization implement the BeanPostProcessor.PostProcessAfterInitialization.
func (*AppRunnerCompoistorProcessor) PostProcessAfterInitialization(obj any, _ string) (v any, err error) {
	return obj, nil
}

// AppRunnerCompositor is the compositor for all AppRunner.
type AppRunnerCompositor struct {
	runners []AppRunner `airmid:"autowire:?"`

	appRunnerNames map[AppRunner]string
}

// GetAppRunnerBeanName will return the AppRunner's BeanName for logger.
func (c *AppRunnerCompositor) GetAppRunnerBeanName(_ context.Context, runner AppRunner) string {
	if c.appRunnerNames == nil {
		return ""
	}

	return c.appRunnerNames[runner]
}

// Run implement AppRunner.Run, it will start all the AppRunner one by one.
func (c *AppRunnerCompositor) Run(ctx context.Context) {
	logger.FromContext(ctx).InfoContext(
		ctx,
		"start to run all AppRunner",
	)
	for _, r := range c.runners {
		logger.FromContext(ctx).InfoContext(
			ctx,
			"AppRunning is starting",
			slog.String("AppRunnerName", c.GetAppRunnerBeanName(ctx, r)),
		)
		r.Run(ctx)
	}
}

// Stop implement AppRunner.Stop, it will stop all the AppRunner
// parallel.
func (c *AppRunnerCompositor) Stop(ctx context.Context) error {
	logger.FromContext(ctx).InfoContext(
		ctx,
		"start to stop all AppRunner",
	)

	err := parallel.Run(ctx, len(c.runners), func(i int) error {
		r := c.runners[i]
		logger.FromContext(ctx).InfoContext(
			ctx,
			"AppRunning is stopping",
			slog.String("AppRunnerName", c.GetAppRunnerBeanName(ctx, r)),
		)

		r.Stop(ctx)
		return nil
	}, parallel.WithConcurrent(len(c.runners)))
	return err
}
