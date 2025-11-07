package xapp

import (
	"context"
	"reflect"
)

var (
	// applicationEventType is the reflect.Type of ApplicationEvent.
	applicationEventType = reflect.TypeOf((*ApplicationEvent)(nil)).Elem()

	// NOTE: when we compare with Method Type, we cannot distinct between different object.method
	//  on retrieve method from object.
	// If the Method Type is retrieve from reflect.Type, it will take the receiver as first argument.

	// onApplicationEventFuncType is the function reflect.Type of ApplicationListener.OnApplicationEvent.
	onApplicationEventFuncType = reflect.TypeOf((*ApplicationListener)(nil)).Elem().Method(0).Type

	// onApplicationEventFuncName is the function name of ApplicationListener.OnApplicationEvent.
	onApplicationEventFuncName = "OnApplicationEvent"

	// onApplicationEventAsyncFuncType is the function reflect.Type of AsyncApplicationListener.OnApplicationEventAsync.
	onApplicationEventAsyncFuncType = reflect.TypeOf((*AsyncApplicationListener)(nil)).Elem().Method(0).Type

	// onApplicationEventAsyncFuncName is the function name of AsyncApplicationListener.OnApplicationEventAsync.
	onApplicationEventAsyncFuncName = "OnApplicationEventAsync"
)

// ApplicationListener is to be implemented by user listener.
type ApplicationListener interface {
	OnApplicationEvent(ctx context.Context, event ApplicationEvent)
}

// AsyncApplicationListener is to be implemented by user async listener.
type AsyncApplicationListener interface {
	OnApplicationEventAsync(ctx context.Context, event ApplicationEvent)
}
