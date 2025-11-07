package beans

import (
	"reflect"

	"github.com/anyvoxel/airmid/xerrors"
)

// IsTypeMatched return true if the beanTyp can convert to targetTyp (assigned).
func IsTypeMatched(targetTyp reflect.Type, beanTyp reflect.Type) bool {
	if targetTyp.Kind() == reflect.Interface {
		return beanTyp.Implements(targetTyp)
	}

	return beanTyp.AssignableTo(targetTyp)
}

// Proxyer is an wrapper for bean object.
type Proxyer interface {
	// OriginalObject return the wrapped bean object
	OriginalObject() any
}

// IndirectTo return the target object which implement specific interface.
func IndirectTo[T any](o any) T {
	var v T

	vv, ok := o.(T)
	if ok {
		return vv
	}

	if po, ok := o.(Proxyer); ok {
		return IndirectTo[T](po.OriginalObject())
	}

	return v
}

// GetBean return the target beanObject.
func GetBean[T any](f BeanFactory, beanName string) (T, error) {
	var v T

	beanObject, err := f.GetBean(beanName)
	if err != nil {
		return v, err
	}

	vv, ok := beanObject.(T)
	if !ok {
		return v, xerrors.Errorf("cannot convert %T to %T", beanObject, v)
	}

	return vv, nil
}
