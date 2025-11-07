package xreflect

import (
	"context"
	"reflect"
)

var (
	// ContextType is the reflect.Type of context.Context.
	ContextType = reflect.TypeOf((*context.Context)(nil)).Elem()

	// ErrorType is the reflect.Type of error.
	ErrorType = reflect.TypeOf((*error)(nil)).Elem()
)
