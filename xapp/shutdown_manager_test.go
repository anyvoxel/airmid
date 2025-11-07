package xapp

import (
	"context"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/utils"
)

func TestShutdownWithSignal(t *testing.T) {
	g := NewWithT(t)
	ctx, cancel := context.WithCancel(context.Background())

	msg := ""
	guard := gomonkey.ApplyFunc(signal.Notify, func(c chan<- os.Signal, _ ...os.Signal) {
		time.Sleep(time.Second)
		c <- os.Interrupt
	})
	defer guard.Reset()

	_ = NewSignalShutdownManager([]ShutdownHandler{
		func(s string) {
			msg = s
			cancel()
		},
	})
	<-ctx.Done()

	g.Expect(msg).To(Equal("Receive signal: interrupt"))
}

func TestShutdownActive(t *testing.T) {
	g := NewWithT(t)
	ctx, cancel := context.WithCancel(context.Background())

	msg := ""

	m := NewSignalShutdownManager([]ShutdownHandler{
		func(s string) {
			msg = s
			cancel()
		},
	})
	m.Shutdown("123")
	<-ctx.Done()

	g.Expect(msg).To(Equal("123"))
}

func TestShutdownMixed(t *testing.T) {
	g := NewWithT(t)
	ctx, cancel := context.WithCancel(context.Background())

	msg := ""

	m := NewSignalShutdownManager([]ShutdownHandler{
		func(s string) {
			msg += s
			cancel()
		},
	})

	var wg utils.WaitGroupWrapper
	for i := 0; i < 100; i++ {
		wg.Wrap(func() {
			m.Shutdown("1")
		})
	}

	wg.Wait()
	<-ctx.Done()
	g.Expect(msg).To(Equal("1"))
}
