package xapp

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/anyvoxel/airmid/utils"
)

// ShutdownHandler is the function type which will be called at shutdown.
type ShutdownHandler func(string)

// ShutdownManager is the interface for manage shutdown for application.
type ShutdownManager interface {
	// Shutdown will exit the manager
	Shutdown(msg string)
}

// signalShutdownManager will shutdown when the signal received.
type signalShutdownManager struct {
	ctx    context.Context
	cancel func()

	handlers     []ShutdownHandler
	shutdownOnce sync.Once
}

// NewSignalShutdownManager will return the ShutdownManager implement.
func NewSignalShutdownManager(handlers []ShutdownHandler, sigs ...os.Signal) ShutdownManager {
	m := &signalShutdownManager{
		handlers: handlers,
	}
	m.ctx, m.cancel = context.WithCancel(context.Background())
	m.Start(sigs...)

	return m
}

func (s *signalShutdownManager) Start(sigs ...os.Signal) {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, sigs...)
		select {
		case <-s.ctx.Done():
			return
		case sig := <-ch:
			s.Shutdown(fmt.Sprintf("Receive signal: %v", sig))
		}
	}()
}

func (s *signalShutdownManager) Shutdown(msg string) {
	s.cancel()

	s.shutdownOnce.Do(func() {
		for _, handler := range s.handlers {
			utils.SafeRun(func() {
				handler(msg)
			})
		}
	})
}
