package beans

import (
	"fmt"
	"reflect"
)

// ConstructorArgumentResolver to resolve arguments.
type ConstructorArgumentResolver interface {
	// Resolve will build the args value from definition
	Resolve(args []ConstructorArgument) ([]reflect.Value, error)
}

type factoryConstructorArgumentResolver struct {
	bf *beanFactoryImpl
}

func (r *factoryConstructorArgumentResolver) Resolve(args []ConstructorArgument) ([]reflect.Value, error) {
	stFields := make([]reflect.StructField, 0, len(args))
	fds := make([]FieldDescriptor, 0, len(args))
	for i, arg := range args {
		st := reflect.StructField{
			Name: fmt.Sprintf("F%d", i),
			Type: arg.Type,
		}

		stFields = append(stFields, st)

		switch {
		case arg.Property != nil:
			fds = append(fds, FieldDescriptor{
				FieldIndex: i,
				Name:       st.Name,
				Typ:        st.Type,
				Unexported: false,
				Property:   arg.Property,
			})
		case arg.Bean != nil:
			fds = append(fds, FieldDescriptor{
				FieldIndex: i,
				Name:       st.Name,
				Typ:        st.Type,
				Unexported: false,
				Bean:       arg.Bean,
			})
		default:
			// We do nothing for constant value
		}
	}
	beanObject := reflect.New(reflect.StructOf(stFields))
	err := r.bf.wireStruct(beanObject, fds)
	if err != nil {
		return nil, err
	}

	ret := make([]reflect.Value, 0, len(args))
	for i, arg := range args {
		switch {
		case arg.Bean != nil || arg.Property != nil:
			// We need get value from beanObject, it has been inject with wireStruct
			ret = append(ret, beanObject.Elem().Field(i))
		default:
			ret = append(ret, arg.Value)
		}
	}

	return ret, nil
}
