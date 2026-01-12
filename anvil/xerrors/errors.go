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

// Package xerrors defines a generic, type-safe error framework.
package xerrors

import (
	"errors"
	"fmt"
	"reflect"
)

// Is alias the errors.Is.
var Is = errors.Is

// As alias the errors.As.
var As = errors.As

// AsType is a generic wrapper around As that simplifies typed error checking.
// It checks if the error `err` can be assigned to a value of type `E` and returns
// that value and a boolean indicating success. `E` must be an error type (e.g., *MyError).
func AsType[E error](err error) (E, bool) {
	var target E
	if As(err, &target) {
		return target, true
	}
	return target, false // target will be the zero value of E
}

// New alias the errors.New.
var New = errors.New

// Errorf formats according to a format specifier and returns the string as a
// value that satisfies error.
var Errorf = fmt.Errorf

// TypedError is a generic, structured error type.
// Its type parameter T carries the semantic meaning and structured data of the error.
type TypedError[T any] struct {
	// Message is an optional, human-friendly description.
	Message string
	// Context is the strongly-typed semantic payload.
	Context T
	// cause is the underlying wrapped error.
	cause error
}

// Messager is the interface that contains the Message method.
type Messager interface {
	// Message return the message of the Messager.
	Message() string
}

// NewTyped creates a new TypedError with the given context.
func NewTyped[T any](context T) *TypedError[T] {
	e := &TypedError[T]{
		Context: context,
	}

	if msger, ok := (any)(context).(Messager); ok {
		e.Message = msger.Message()
	}

	return e
}

// WithCause adds an underlying cause to the error.
func (e *TypedError[T]) WithCause(cause error) *TypedError[T] {
	e.cause = cause
	return e
}

// WithMessage adds a human-friendly message to the error.
func (e *TypedError[T]) WithMessage(message string) *TypedError[T] {
	e.Message = message
	return e
}

// Error implements the error interface, providing a human-readable representation.
func (e *TypedError[T]) Error() string {
	typeName := reflect.TypeOf(e.Context).Name()

	baseMsg := typeName
	if e.Message != "" {
		baseMsg = fmt.Sprintf("%s: %s", typeName, e.Message)
	}

	if e.cause != nil {
		return fmt.Sprintf("%s (cause: %s)", baseMsg, e.cause.Error())
	}

	return baseMsg
}

// Unwrap provides access to the underlying error for error chains.
func (e *TypedError[T]) Unwrap() error {
	return e.cause
}

// --- Semantic Context Structs ---

// NotImplement indicates a feature or function is not implemented.
type NotImplement struct{}

// NotFound indicates a requested resource could not be found.
type NotFound struct {
	// ResourceType is the type of the resource that was not found.
	ResourceType string
	// ResourceID is the ID of the resource that was not found.
	ResourceID any
}

// Message returns the structured error message.
func (e NotFound) Message() string {
	return fmt.Sprintf("%s with ID '%v' not found", e.ResourceType, e.ResourceID)
}

// Duplicate indicates an item already exists.
type Duplicate struct {
	// ResourceType is the type of the resource that already exists.
	ResourceType string
	// ResourceID is the ID of the resource that already exists.
	ResourceID any
}

// Message returns the structured error message.
func (e Duplicate) Message() string {
	return fmt.Sprintf("%s with ID '%v' already exists", e.ResourceType, e.ResourceID)
}

// Continue indicates a transient state that can be continued.
type Continue struct{}

// Retryable indicates an operation failed but can be retried.
type Retryable struct{}

// NonRetryable indicates an operation failed and should not be retried.
type NonRetryable struct{}

// ConversionError indicates a failure during data type conversion.
type ConversionError struct {
	// SourceType is the type of the source value.
	SourceType string
	// TargetType is the type of the target value.
	TargetType string
	// Value is the value that failed to convert.
	Value any
}

// Message returns the structured error message.
func (e ConversionError) Message() string {
	return fmt.Sprintf("cannot convert from %s to %s", e.SourceType, e.TargetType)
}

// UnsupportedTypeError indicates that a given type is not supported.
type UnsupportedTypeError struct {
	// TargetType is the unsupported type.
	TargetType string
}

// Message returns the structured error message.
func (e UnsupportedTypeError) Message() string {
	return fmt.Sprintf("type %s is not supported", e.TargetType)
}

// InvalidArgument indicates that a function or method received an invalid argument.
type InvalidArgument struct {
	// ArgumentName is the name of the invalid argument.
	ArgumentName string
	// Reason is the reason why the argument is invalid.
	Reason string
}

// Message returns the structured error message.
func (e InvalidArgument) Message() string {
	return fmt.Sprintf("invalid argument '%s': %s", e.ArgumentName, e.Reason)
}

// InitializationError indicates an error occurred during an initialization phase.
type InitializationError struct {
	// Reason is the reason why the initialization failed.
	Reason string
}

// Message returns the structured error message.
func (e InitializationError) Message() string {
	return e.Reason
}

// SettingError indicates an error occurred while setting a value.
type SettingError struct {
	// Key is the key of the setting.
	Key string
	// Value is the value of the setting.
	Value any
	// Reason is the reason why the setting failed.
	Reason string
}

// Message returns the structured error message.
func (e SettingError) Message() string {
	return fmt.Sprintf("cannot set key '%s' with value '%v': %s", e.Key, e.Value, e.Reason)
}

// CircularDependency indicates a circular dependency was detected.
type CircularDependency struct {
	// BeanName is the name of the bean that caused the circular dependency.
	BeanName string
	// Path is the dependency path that formed a cycle.
	Path []string
}

// Message returns the structured error message.
func (e CircularDependency) Message() string {
	return fmt.Sprintf("cannot get bean '%s' circularly", e.BeanName)
}

// InvalidScopeError indicates an invalid scope name was used.
type InvalidScopeError struct {
	// ScopeName is the invalid scope name.
	ScopeName string
	// Reason is the reason why the scope is invalid.
	Reason string
}

// Message returns the structured error message.
func (e InvalidScopeError) Message() string {
	return fmt.Sprintf("invalid scope name '%s': %s", e.ScopeName, e.Reason)
}

// TooManyCandidatesError indicates too many candidates were found for a given type.
type TooManyCandidatesError struct {
	// Count is the number of candidates found.
	Count int
	// Type is the type of the candidates.
	Type string
	// Field is the field where the candidates were found.
	Field string
	// IsPrimary indicates whether the candidates are primary.
	IsPrimary bool
}

// Message returns the structured error message.
func (e TooManyCandidatesError) Message() string {
	if e.IsPrimary {
		return fmt.Sprintf("'%v' primary candidates found for field '%v' with type %s", e.Count, e.Field, e.Type)
	}
	return fmt.Sprintf("'%v' candidates found for field '%v' with type %s", e.Count, e.Field, e.Type)
}

// NoCandidateError indicates no candidate was found for a given type.
type NoCandidateError struct {
	// Type is the type of the candidate.
	Type string
	// Field is the field where the candidate was not found.
	Field string
}

// Message returns the structured error message.
func (e NoCandidateError) Message() string {
	return fmt.Sprintf("no candidate found for field '%v' with type %s", e.Field, e.Type)
}

// InvalidBeanDefinitionError indicates an invalid bean definition.
type InvalidBeanDefinitionError struct {
	// BeanName is the name of the invalid bean.
	BeanName string
	// Reason is the reason why the bean definition is invalid.
	Reason string
}

// Message returns the structured error message.
func (e InvalidBeanDefinitionError) Message() string {
	return fmt.Sprintf("invalid bean definition for bean '%s': %s", e.BeanName, e.Reason)
}

// StartupHandlerError indicates an error occurred within an application startup handler.
type StartupHandlerError struct {
	// HandlerName is the name of the startup handler.
	HandlerName string
	// Phase is the phase of the startup handler.
	Phase string // e.g., "BeforeLoadProps", "AfterLoadProps"
}

// Message returns the structured error message.
func (e StartupHandlerError) Message() string {
	return fmt.Sprintf("startup handler '%s' failed at phase '%s'", e.HandlerName, e.Phase)
}

// UnsetableError indicates that a value cannot be set.
type UnsetableError struct {
	// TargetType is the type of the value that cannot be set.
	TargetType string
}

// Message returns the structured error message.
func (e UnsetableError) Message() string {
	return fmt.Sprintf("value of type %s cannot be set", e.TargetType)
}

// InvalidTagError indicates that a struct tag is invalid.
type InvalidTagError struct {
	// Tag is the invalid tag.
	Tag string
	// Reason is the reason why the tag is invalid.
	Reason string
}

// Message returns the structured error message.
func (e InvalidTagError) Message() string {
	return fmt.Sprintf("invalid tag '%s': %s", e.Tag, e.Reason)
}

// TooManyConstructorsError indicates too many constructors were found for a given type.
type TooManyConstructorsError struct {
	// Count is the number of constructors found.
	Count int
	// Type is the type of the constructors.
	Type string
}

// Message returns the structured error message.
func (e TooManyConstructorsError) Message() string {
	return fmt.Sprintf("too many constructors found for type %s: %d", e.Type, e.Count)
}

// TooManyInstancesError indicates too many instances were found for a singleton scope object.
type TooManyInstancesError struct{}

// Message returns the structured error message.
func (e TooManyInstancesError) Message() string {
	return "too many instances found for a singleton scope object"
}

// PanicError indicates that a panic occurred.
type PanicError struct {
	// Recovered is the recovered value from the panic.
	Recovered any
	// Stack is the stack trace of the panic.
	Stack string
}

// Message returns the structured error message.
func (e PanicError) Message() string {
	return fmt.Sprintf("recovered from: '%v'", e.Recovered)
}
