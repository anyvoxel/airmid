package runner1

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/anyvoxel/airmid/anvil"
	"github.com/anyvoxel/airmid/app"
	"github.com/anyvoxel/airmid/ioc"
	slogctx "github.com/veqryn/slog-context"
)

type runner1 struct {
	Attr        string `airmid:"value:${attr}"`
	publisher   app.ApplicationEventPublisher
	beanFactory ioc.BeanFactory
	application app.Application
}

type runner1RunEvent struct {
	desp string
	app.ApplicationEvent
}

var (
	k1 = &[]int{1}
)

func (r *runner1) SetApplication(application app.Application) {
	r.application = application
}

func (r *runner1) SetApplicationEventPublisher(publisher app.ApplicationEventPublisher) {
	r.publisher = publisher
}

func (r *runner1) SetBeanFactory(beanFactory ioc.BeanFactory) {
	r.beanFactory = beanFactory
}

func (r *runner1) Run(ctx context.Context) {
	slogctx.FromCtx(ctx).InfoContext(ctx, fmt.Sprintf("Runner1.Run called: attr='%v'\n", r.Attr))

	r.publisher.PublishEvent(context.WithValue(ctx, k1, "v1"), &runner1RunEvent{
		ApplicationEvent: app.NewDefaultApplicationEvent(r),
		desp:             "e1",
	})

	attr, _ := r.application.Get(ctx, "attr")
	slogctx.FromCtx(ctx).InfoContext(ctx, fmt.Sprintf("Runner1 get properties: %s=%v", "attr", attr))

	r.application.Submit(func() {
		slogctx.FromCtx(ctx).InfoContext(ctx, "submit running in gpool")
	})
}

// todo(xiangqilin): stop publish event.
func (r *runner1) Stop(ctx context.Context) {
	slogctx.FromCtx(ctx).InfoContext(ctx, "Runner1.Stop called: runner is stopping")
}

type eventListener struct {
}

// NOTE: we doesn't need to instantiate this object manual, the app will call factory.PreInstantiateSingletons to
//
//	pre initialize all singleton & non-lazy-mode instance.
func (l *eventListener) OnRunEvent(ctx context.Context, event *runner1RunEvent) {
	slogctx.FromCtx(ctx).InfoContext(ctx, "OnRunEvent", slog.Any("k1", ctx.Value(k1)), slog.Any("event", event.desp))
}

func init() {
	anvil.Must(
		app.RegisterBeanDefinition("runner1", ioc.MustNewBeanDefinition(reflect.TypeOf((*runner1)(nil)))),
	)
	anvil.Must(
		app.RegisterBeanDefinition("eventListener1", ioc.MustNewBeanDefinition(reflect.TypeOf((*eventListener)(nil)))),
	)
}
