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
	"os"
	"reflect"

	"github.com/anyvoxel/airmid/anvil/xerrors"
	"github.com/anyvoxel/airmid/ioc"
)

// LoggerProvider is used to provide a logger for application usage.
type LoggerProvider interface {
	GetLogger() *slog.Logger
}

type loggerStartupHandlerConfiguration struct {
	loggerProviders []LoggerProvider `airmid:"autowire:?"`

	handlerType         string `airmid:"value:${airmid.logger.handler.type:=json}"`
	handlerOptAddSource bool   `airmid:"value:${airmid.logger.handler.opt.source:=true}"`
	handlerOptLevel     string `airmid:"value:${airmid.logger.handler.opt.level:=INFO}"`
}

type loggerStartupHandler struct{}

var (
	_ ApplicationStartupHandler = (*loggerStartupHandler)(nil)
)

func (*loggerStartupHandler) Name() string {
	return "LoggerStartupHandler"
}

func (*loggerStartupHandler) BeforeLoadProps(_ context.Context, app *airmidApplication, _ *option) error {
	return app.RegisterBeanDefinition(
		"airmidLoggerStartupHandlerConfiguraion",
		ioc.MustNewBeanDefinition(
			reflect.TypeOf((*loggerStartupHandlerConfiguration)(nil)),
			ioc.WithLazyMode(),
		),
	)
}

// AfterLoadProps change the default logger to the bean which implement it.
func (*loggerStartupHandler) AfterLoadProps(ctx context.Context, app *airmidApplication, _ *option) error {
	loggerC, err := ioc.GetBean[*loggerStartupHandlerConfiguration](app, "airmidLoggerStartupHandlerConfiguraion")
	if err != nil {
		return err
	}

	if len(loggerC.loggerProviders) > 0 {
		slog.SetDefault(loggerC.loggerProviders[0].GetLogger())
		slog.InfoContext(
			ctx, "change default logger to user specify bean",
			slog.String("LoggerBeanType", reflect.TypeOf(loggerC.loggerProviders[0]).Elem().String()),
		)

		return nil
	}

	level := slog.LevelVar{}
	err = level.UnmarshalText([]byte(loggerC.handlerOptLevel))
	if err != nil {
		return err
	}
	opt := &slog.HandlerOptions{
		AddSource:   loggerC.handlerOptAddSource,
		Level:       &level,
		ReplaceAttr: nil,
	}

	switch loggerC.handlerType {
	case "text":
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, opt)))
	case "json":
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opt)))
	default:
		return xerrors.Errorf("unknown logger handler type %s", loggerC.handlerType)
	}
	return nil
}

func (*loggerStartupHandler) BeforeStartRunner(_ context.Context, _ *airmidApplication, _ *option) error {
	return nil
}
