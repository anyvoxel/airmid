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
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/anyvoxel/airmid/anvil"
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
			anvil.SafeRun(context.TODO(), func(context.Context) {
				handler(msg)
			})
		}
	})
}
