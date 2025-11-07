package beans

import (
	"context"
	"log/slog"
	"reflect"
	"unsafe"

	"github.com/anyvoxel/airmid/logger"
)

// PropertyValues is the values holder for bean property.
type PropertyValues interface {
	// AddValue will add the field value to container
	AddValue(fieldIndex int, value reflect.Value)

	// SetProperty will update the field value with property
	SetProperty(obj reflect.Value, fieldDescriptors []FieldDescriptor) error
}

type propertyValuesImpl struct {
	values map[int]reflect.Value
}

// NewPropertyValues return the PropertyValues impl.
func NewPropertyValues() PropertyValues {
	return &propertyValuesImpl{
		values: make(map[int]reflect.Value),
	}
}

func (p *propertyValuesImpl) AddValue(fieldIndex int, value reflect.Value) {
	p.values[fieldIndex] = value
}

func (p *propertyValuesImpl) SetProperty(obj reflect.Value, fieldDescriptors []FieldDescriptor) error {
	for _, fd := range fieldDescriptors {
		value, ok := p.values[fd.FieldIndex]
		if !ok {
			logger.FromContext(context.TODO()).DebugContext(
				context.TODO(),
				"the value of field not found, skip it",
				slog.Int("FieldIndex", fd.FieldIndex),
				slog.String("FieldName", fd.Name),
			)
			continue
		}

		fv := obj.Field(fd.FieldIndex)
		if !fd.Unexported {
			fv.Set(value)
			continue
		}

		// The field is unexported (private field), we cannot directly use Set
		// It will panic with 'reflect.Value.Set using unaddressable value'
		reflect.NewAt(fv.Type(), unsafe.Pointer(fv.UnsafeAddr())).Elem().Set(value)
	}

	return nil
}
