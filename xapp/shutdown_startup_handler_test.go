package xapp

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/beans"
)

func TestShutdownStartupHandler(t *testing.T) {
	g := NewWithT(t)

	app := &airmidApplication{
		BeanFactory: beans.NewBeanFactory(),
	}
	h := &shutdownStartupHandler{}
	err := h.BeforeLoadProps(context.Background(), app, nil)
	g.Expect(err).ToNot(HaveOccurred())

	err = h.AfterLoadProps(context.Background(), app, nil)
	g.Expect(err).ToNot(HaveOccurred())
}
