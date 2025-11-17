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

	"github.com/anyvoxel/airmid/anvil/parallel"
	"github.com/anyvoxel/airmid/ioc"
	slogctx "github.com/veqryn/slog-context"
)

// Runner is the interface for runner.
type Runner interface {
	// Run will been called when the application is ready.
	Run(ctx context.Context)

	// Stop will be called when the application is shutting down.
	Stop(ctx context.Context)
}

// RunnerCompoistorProcessor is the postprocessor for AppRunner.
type RunnerCompoistorProcessor struct {
	appRunnerNames map[Runner]string
}

// PostProcessBeforeInitialization implement the BeanPostProcessor.PostProcessBeforeInitialization.
func (p *RunnerCompoistorProcessor) PostProcessBeforeInitialization(
	_ context.Context, obj any, beanName string) (v any, err error) {
	if vobj := ioc.IndirectTo[Runner](obj); vobj != nil {
		// NOTE: we cannot return a WrapObject for AppRunner, because
		// when we use embedded interface, the WrapObject didn't implement
		// those interface which implemented by concrete type
		p.appRunnerNames[vobj] = beanName
	}

	return obj, nil
}

// PostProcessAfterInitialization implement the BeanPostProcessor.PostProcessAfterInitialization.
func (*RunnerCompoistorProcessor) PostProcessAfterInitialization(
	_ context.Context, obj any, _ string) (v any, err error) {
	return obj, nil
}

// RunnerCompositor is the compositor for all AppRunner.
type RunnerCompositor struct {
	runners []Runner `airmid:"autowire:?"`

	appRunnerNames map[Runner]string
}

// GetAppRunnerBeanName will return the AppRunner's BeanName for logger.
func (c *RunnerCompositor) GetAppRunnerBeanName(_ context.Context, runner Runner) string {
	if c.appRunnerNames == nil {
		return ""
	}

	return c.appRunnerNames[runner]
}

// Run implement AppRunner.Run, it will start all the AppRunner one by one.
func (c *RunnerCompositor) Run(ctx context.Context) {
	slogctx.FromCtx(ctx).InfoContext(
		ctx,
		"start to run all AppRunner",
	)
	for _, r := range c.runners {
		slogctx.FromCtx(ctx).InfoContext(
			ctx,
			"AppRunning is starting",
			slog.String("AppRunnerName", c.GetAppRunnerBeanName(ctx, r)),
		)
		r.Run(ctx)
	}
}

// Stop implement AppRunner.Stop, it will stop all the AppRunner
// parallel.
func (c *RunnerCompositor) Stop(ctx context.Context) error {
	slogctx.FromCtx(ctx).InfoContext(
		ctx,
		"start to stop all AppRunner",
	)

	err := parallel.Run(ctx, len(c.runners), func(i int) error {
		r := c.runners[i]
		slogctx.FromCtx(ctx).InfoContext(
			ctx,
			"AppRunning is stopping",
			slog.String("AppRunnerName", c.GetAppRunnerBeanName(ctx, r)),
		)

		r.Stop(ctx)
		return nil
	}, parallel.WithConcurrent(len(c.runners)))
	return err
}
