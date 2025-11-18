package runner2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/anyvoxel/airmid/anvil"
	"github.com/anyvoxel/airmid/app"
	"github.com/anyvoxel/airmid/ioc"
)

// nolint
type Runner2 struct {
	Attr    int           `airmid:"value:${attr2}"`
	Port    int           `airmid:"value:${port}"`
	timeout time.Duration `airmid:"value:${timeout}"`

	Config       *config
	Config2      *config2     `airmid:"autowire:?,optional"`
	configReader configReader `airmid:"autowire:?"`

	serv *http.Server
	r    metric.Reader
}

func (*Runner2) NewRunner2(Config *config) *Runner2 {
	return &Runner2{
		Config: Config,
	}
}

type config struct {
	Address *string `airmid:"value:${address}"`
}

type config2 struct {
}

type configReader interface {
	Name() string
}

func (r *Runner2) GetOptions() []metric.Option {
	return []metric.Option{
		metric.WithReader(func() metric.Reader {
			exporter, err := stdoutmetric.New()
			if err != nil {
				panic(err)
			}

			rr := metric.NewPeriodicReader(exporter)
			r.r = rr
			return rr
		}()),
		metric.WithResource(func() *resource.Resource {
			res, err := resource.New(
				context.Background(),
				resource.WithContainerID(),
				resource.WithFromEnv(),
				resource.WithHost(),
				resource.WithOS(),
				resource.WithProcess(),
				resource.WithAttributes(attribute.Key("k1").String("v1")),
			)
			if err != nil {
				panic(err)
			}

			return res
		}()),
	}
}

type configReaderImpl struct {
	nam   string
	value string `airmid:"value:${not_exists:=not_found_value}"`
}

func (*configReaderImpl) NewConfigReaderImpl(nam string) *configReaderImpl {
	return &configReaderImpl{
		nam: nam,
	}
}

func (c *configReaderImpl) Name() string {
	return "configReaderImpl/" + c.nam + "  value:" + c.value
}

// nolint
func (r *Runner2) Run(ctx context.Context) {
	fmt.Printf("Runner2.Run called: attr='%v', port='%v', addr='%v', name='%v', timeout='%v', config2='%v'\n", r.Attr, r.Port, *r.Config.Address, r.configReader.Name(), r.timeout, r.Config2)

	go func() {
		err := r.serv.ListenAndServe()
		if err != nil {
			fmt.Printf("Runner2.Run ListenAndServe failed, %v\n", err)
			app.Shutdown()
		}
	}()
}

func (r *Runner2) AfterPropertiesSet(context.Context) error {
	handler := http.NewServeMux()
	handler.Handle("/", http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		fmt.Printf("Handler http request\n")
		response.Write([]byte("hello"))
	}))
	handler.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var mockData metricdata.ResourceMetrics
		err := r.r.Collect(req.Context(), &mockData)
		if err != nil {
			panic(err)
		}

		data, err := json.MarshalIndent(mockData, "", "  ")
		if err != nil {
			panic(err)
		}

		w.Write(data)
	}))
	serv := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", r.Port),
		Handler: handler,
	}
	r.serv = serv
	return nil
}

func (r *Runner2) Stop(ctx context.Context) {
	_ = r.serv.Shutdown(ctx)
}

func init() {
	anvil.Must(
		app.RegisterBeanDefinition(
			"runner2",
			ioc.MustNewBeanDefinition(
				reflect.TypeOf((*Runner2)(nil)),
				ioc.WithConstructorArguments([]ioc.ConstructorArgument{
					{
						Bean: &ioc.BeanFieldDescriptor{
							Name: "config",
						},
					},
				}),
			),
		),
	)
	anvil.Must(
		app.RegisterBeanDefinition("config", ioc.MustNewBeanDefinition(reflect.TypeOf((*config)(nil)))),
	)
	anvil.Must(
		app.RegisterBeanDefinition(
			"configReader",
			ioc.MustNewBeanDefinition(
				reflect.TypeOf((*configReaderImpl)(nil)),
				ioc.WithConstructorArguments([]ioc.ConstructorArgument{
					{
						Property: &ioc.PropertyFieldDescriptor{
							Name: "name",
						},
					},
				}),
			),
		),
	)
}
