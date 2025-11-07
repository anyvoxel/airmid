package beans

import (
	"reflect"

	"github.com/anyvoxel/airmid/utils/xreflect"
)

// Constructor is the constructor for Bean.
type Constructor interface {
	// NewObject will return the object from constructor
	NewObject(resolver ConstructorArgumentResolver) (reflect.Value, error)
}

// ReflectConstructor will build the object from reflect.
type ReflectConstructor struct {
	typ reflect.Type
}

// NewObject return the object from Reflect.
func (r *ReflectConstructor) NewObject(_ ConstructorArgumentResolver) (reflect.Value, error) {
	return xreflect.NewValue(r.typ)
}

// MethodConstructor will build the object from typ's method.
type MethodConstructor struct {
	typ    reflect.Type
	method reflect.Method

	args []ConstructorArgument
}

// NewMethodConstructor return the constructor for method.
func NewMethodConstructor(
	typ reflect.Type, method reflect.Method, args []ConstructorArgument,
) (*MethodConstructor, error) {
	for i := 0; i < len(args); i++ {
		// We must use +1 to indicate the argument's type, because the method first argument it (*self)
		args[i].Type = method.Type.In(i + 1)
	}

	return &MethodConstructor{
		typ:    typ,
		method: method,
		args:   args,
	}, nil
}

// NewObject return the object from Constructor.
func (m *MethodConstructor) NewObject(resolver ConstructorArgumentResolver) (reflect.Value, error) {
	ins, err := resolver.Resolve(m.args)
	if err != nil {
		return reflect.Value{}, err
	}

	argin := []reflect.Value{reflect.Zero(m.typ)}
	argin = append(argin, ins...)

	ret := m.method.Func.Call(argin)
	if len(ret) == 1 {
		return ret[0], nil
	}

	ret1 := ret[1].Interface()
	if ret1 == nil {
		return ret[0], nil
	}

	return ret[0], ret1.(error) //nolint:revive
}
