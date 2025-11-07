package xapp

import (
	"time"
)

// ApplicationEvent is to be implemented by user defined event.
type ApplicationEvent interface {
	// GetTimestamp return the time in milliseconds when the event occurred.
	GetTimestamp() uint64

	// GetSource return the object on which the Event initially occurred.
	GetSource() any
}

// DefaultApplicationEvent is the default implement of application event.
type DefaultApplicationEvent struct {
	Source    any    `json:"-"`
	Timestamp uint64 `json:"timestamp"`
}

// GetTimestamp impelement the ApplicationEvent.GetTimestamp.
func (e *DefaultApplicationEvent) GetTimestamp() uint64 {
	return e.Timestamp
}

// GetSource impelement the ApplicationEvent.GetSource.
func (e *DefaultApplicationEvent) GetSource() any {
	return e.Source
}

// NewDefaultApplicationEvent return the default impl.
func NewDefaultApplicationEvent(source any) *DefaultApplicationEvent {
	return &DefaultApplicationEvent{
		Source:    source,
		Timestamp: uint64(time.Now().UnixNano()) / 1000,
	}
}
