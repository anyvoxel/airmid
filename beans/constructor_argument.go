package beans

import (
	"reflect"
)

// ConstructorArgument is the definition for constructor arguments.
type ConstructorArgument struct {
	// Type is the value's type of argument
	Type reflect.Type

	// Property is the field property descriptor
	Property *PropertyFieldDescriptor

	// Bean is the field bean descriptor
	Bean *BeanFieldDescriptor

	// Value is the fixed argument value
	Value reflect.Value
}
