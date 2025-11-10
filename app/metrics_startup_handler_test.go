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
	"testing"

	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/attribute"

	"github.com/anyvoxel/airmid/ioc"
)

func TestMetricsStartupHandler(t *testing.T) {
	g := NewWithT(t)

	app := &airmidApplication{
		BeanFactory: ioc.NewBeanFactory(),
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
