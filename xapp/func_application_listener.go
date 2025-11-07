package xapp

import "context"

// FuncApplicationListener is the handler for compose function.
// NOTE: this is used for unit test, be careful to use it.
type FuncApplicationListener struct {
	OnEventFunc      func(ctx context.Context, event ApplicationEvent)
	OnEventAsyncFunc func(ctx context.Context, event ApplicationEvent)
}

// OnApplicationEvent implement ApplicationListener.OnApplicationEvent.
func (l *FuncApplicationListener) OnApplicationEvent(ctx context.Context, event ApplicationEvent) {
	if l.OnEventFunc != nil {
		l.OnEventFunc(ctx, event)
	}
}

// OnApplicationEventAsync implement AsyncApplicationListener.OnApplicationEventAsync.
func (l *FuncApplicationListener) OnApplicationEventAsync(ctx context.Context, event ApplicationEvent) {
	if l.OnEventAsyncFunc != nil {
		l.OnEventAsyncFunc(ctx, event)
	}
}

var (
	_ ApplicationListener      = &FuncApplicationListener{}
	_ AsyncApplicationListener = &FuncApplicationListener{}
)
