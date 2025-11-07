package xapp

import (
	"context"
)

// ListenerInvoker will invoke listener function with given argument.
type ListenerInvoker interface {
	Invoke(ctx context.Context, event ApplicationEvent)
}
