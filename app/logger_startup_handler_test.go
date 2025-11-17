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
	"testing"

	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/ioc"
)

type testLoggerProvider struct {
}

func (t *testLoggerProvider) GetLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
}

func TestLoggerStartupHandler(t *testing.T) {
	t.Run("user provide logger", func(t *testing.T) {
		g := NewWithT(t)

		app := &airmidApplication{
			BeanFactory: ioc.NewBeanFactory(),
		}
		app.RegisterBeanDefinition("testLoggerProvider",
			ioc.MustNewBeanDefinition(
				reflect.TypeOf((*testLoggerProvider)(nil)),
			),
		)
		h := &loggerStartupHandler{}
		err := h.BeforeLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())

		err = h.AfterLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("default handler", func(t *testing.T) {
		g := NewWithT(t)
		app := &airmidApplication{
			BeanFactory: ioc.NewBeanFactory(),
		}
		h := &loggerStartupHandler{}
		err := h.BeforeLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())

		err = h.AfterLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("json handler", func(t *testing.T) {
		g := NewWithT(t)
		app := &airmidApplication{
			BeanFactory: ioc.NewBeanFactory(),
		}
		h := &loggerStartupHandler{}
		app.Set(context.Background(), "airmid.logger.handler.type", "json")
		err := h.BeforeLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())

		err = h.AfterLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("invalid level", func(t *testing.T) {
		g := NewWithT(t)
		app := &airmidApplication{
			BeanFactory: ioc.NewBeanFactory(),
		}
		h := &loggerStartupHandler{}
		app.Set(context.Background(), "airmid.logger.handler.opt.level", "iii")
		err := h.BeforeLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())

		err = h.AfterLoadProps(context.Background(), app, nil)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`level string "iii": unknown name`))
	})

	t.Run("invalid type", func(t *testing.T) {
		g := NewWithT(t)
		app := &airmidApplication{
			BeanFactory: ioc.NewBeanFactory(),
		}
		h := &loggerStartupHandler{}
		app.Set(context.Background(), "airmid.logger.handler.type", "iii")
		err := h.BeforeLoadProps(context.Background(), app, nil)
		g.Expect(err).ToNot(HaveOccurred())

		err = h.AfterLoadProps(context.Background(), app, nil)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(MatchRegexp(`unknown logger handler type iii`))
	})
}
