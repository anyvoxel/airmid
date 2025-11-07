package beans

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/anyvoxel/airmid/utils/pointer"
	"github.com/anyvoxel/airmid/xerrors"
)

const (
	// AirmidTagName is the tag name for airmid struct tag field.
	AirmidTagName string = "airmid"

	// ValueContentPrefix is the prefix for property field.
	ValueContentPrefix string = "value:"

	// AutowireContentPrefix is the prefix for bean field.
	AutowireContentPrefix string = "autowire:"

	// OptionalAutowireField is the constant value for optional bean field.
	OptionalAutowireField string = "optional"
)

// FieldDescriptor is the descriptor for struct field
// The struct tag must format as:
//  1. `airmid:"value:${name,default}"` for property field
//  2. `airmid:"autowire:name,optional"` for bean field
type FieldDescriptor struct {
	FieldIndex int
	Name       string
	Typ        reflect.Type
	Unexported bool

	// Property is the field property descriptor.
	// The field should be marked as value=${name:=default}
	Property *PropertyFieldDescriptor

	// Bean is the field bean descriptor.
	// The field should be marked as autowire=name,optional
	Bean *BeanFieldDescriptor
}

// PropertyFieldDescriptor is the descriptor for property autowired value.
type PropertyFieldDescriptor struct {
	Name    string
	Default *string
}

// BeanFieldDescriptor is the descriptor for bean autowired value.
type BeanFieldDescriptor struct {
	Name     string
	Optional bool
}

// NewFieldDescriptor will return the descriptor from struct field
// 1. If the field doesn't have airmid tag, return nil, nil
// 2. If the field have airmid tag, return the field descriptor, otherwise return err.
func NewFieldDescriptor(field reflect.StructField, idx int) (*FieldDescriptor, error) {
	tag, ok := field.Tag.Lookup(AirmidTagName)
	if !ok {
		return nil, nil //nolint:nilnil
	}

	fd := &FieldDescriptor{
		FieldIndex: idx,
		Name:       field.Name,
		Typ:        field.Type,
		Unexported: false,
	}
	if field.PkgPath != "" {
		fd.Unexported = true
	}

	switch {
	case strings.HasPrefix(tag, ValueContentPrefix):
		v, err := NewPropertyFieldDescriptor(tag[len(ValueContentPrefix):])
		if err != nil {
			return nil, err
		}
		fd.Property = v
	case strings.HasPrefix(tag, AutowireContentPrefix):
		v, err := NewBeanFieldDescriptor(tag[len(AutowireContentPrefix):])
		if err != nil {
			return nil, err
		}
		fd.Bean = v
	default:
		return nil, xerrors.Errorf("Invalid tag '%v', it must start with 'value:' or 'autowire:'", tag)
	}

	return fd, nil
}

var (
	valueRegex = regexp.MustCompile(`^\$\{(.*)\}$`)
)

// NewPropertyFieldDescriptor will return the descriptor from tag value.
func NewPropertyFieldDescriptor(value string) (*PropertyFieldDescriptor, error) {
	res := valueRegex.FindAllStringSubmatch(value, -1)
	if len(res) != 1 || len(res[0]) != 2 {
		return nil, xerrors.Errorf("Invalid value '%v', it must format as ${x.y}", value)
	}

	value = res[0][1]
	if len(value) == 0 {
		return nil, xerrors.Errorf("Required value content, it cann't be empty")
	}

	vv := strings.SplitN(value, ":=", 2)
	fd := &PropertyFieldDescriptor{
		Name: vv[0],
	}
	if len(vv) > 1 {
		fd.Default = pointer.StringPtr(vv[1])
	}
	return fd, nil
}

// NewBeanFieldDescriptor will return the descriptor from tag autowire.
func NewBeanFieldDescriptor(value string) (*BeanFieldDescriptor, error) {
	if len(value) == 0 {
		return nil, xerrors.Errorf("Required autowire content, it cann't be empty")
	}

	vv := strings.SplitN(value, ",", 2)
	fd := &BeanFieldDescriptor{
		Name: vv[0],
	}
	if len(vv) > 1 {
		if vv[1] != OptionalAutowireField {
			return nil, xerrors.Errorf("Invalid autowire '%v', it must be '%v'", vv[1], OptionalAutowireField)
		}

		fd.Optional = true
	}

	return fd, nil
}
