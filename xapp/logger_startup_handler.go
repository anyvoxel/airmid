package xapp

import (
	"context"
	"log/slog"
	"os"
	"reflect"

	"github.com/anyvoxel/airmid/beans"
	"github.com/anyvoxel/airmid/xerrors"
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
		beans.MustNewBeanDefinition(
			reflect.TypeOf((*loggerStartupHandlerConfiguration)(nil)),
			beans.WithLazyMode(),
		),
	)
}

// AfterLoadProps change the default logger to the bean which implement it.
func (*loggerStartupHandler) AfterLoadProps(ctx context.Context, app *airmidApplication, _ *option) error {
	loggerC, err := beans.GetBean[*loggerStartupHandlerConfiguration](app, "airmidLoggerStartupHandlerConfiguraion")
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
