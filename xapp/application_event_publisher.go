package xapp

import (
	"context"
)

// ApplicationEventPublisher is the interface that encapsulates event publication functionality.
type ApplicationEventPublisher interface {
	// PublishEvent will notify all matching listeners registered with this application of an application event.
	PublishEvent(ctx context.Context, event ApplicationEvent)
}
