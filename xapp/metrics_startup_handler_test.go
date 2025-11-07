package xapp

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/attribute"

	"github.com/anyvoxel/airmid/beans"
)

func TestMetricsStartupHandler(t *testing.T) {
	g := NewWithT(t)

	app := &airmidApplication{
		BeanFactory: beans.NewBeanFactory(),
	}
	h := &metricsStartupHandler{}
	err := h.BeforeLoadProps(context.Background(), app, nil)
	g.Expect(err).ToNot(HaveOccurred())

	err = h.AfterLoadProps(context.Background(), app, nil)
	g.Expect(err).ToNot(HaveOccurred())
}

func TestConvertOptionToAttributes(t *testing.T) {
	t.Run("nil option", func(t *testing.T) {
		g := NewWithT(t)
		g.Expect(convertOptionToAttributes(nil)).To(BeEmpty())
	})

	t.Run("normal test", func(t *testing.T) {
		g := NewWithT(t)
		opt := &option{
			attrs: []Attribute{
				{
					Key:   "k1",
					Value: "v1",
				},
				{
					Key:   "k2",
					Value: "v2",
				},
			},
		}
		g.Expect(convertOptionToAttributes(opt)).To(Equal([]attribute.KeyValue{
			attribute.String("k1", "v1"),
			attribute.String("k2", "v2"),
		}))
	})
}

func TestMetricsStartupHandlerBeforeStartRunner(t *testing.T) {
	t.Run("should return nil", func(t *testing.T) {
		g := NewWithT(t)
		m := &metricsStartupHandler{}
		g.Expect(m.BeforeStartRunner(context.Background(), nil, &option{
			attrs: []Attribute{
				{
					Key:   "k1",
					Value: "v1",
				},
				{
					Key:   "k2",
					Value: "v2",
				},
			},
		})).To(Succeed())
	})
}
