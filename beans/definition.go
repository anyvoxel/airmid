package beans

import (
	"reflect"

	"github.com/anyvoxel/airmid/xerrors"
)

const (
	// ScopeSingleton is the scope identifier for the singleton scope.
	// The bean will been initialize exactly once and cached.
	ScopeSingleton string = "singleton"

	// ScopePrototype is the scope identifier for the prototype scope.
	// The bean will been initialize when every get is called.
	ScopePrototype string = "prototype"
)

// BeanDefinition is the definition of bean.
type BeanDefinition interface {
	// Type return the bean type
	Type() reflect.Type

	// Name return the bean name
	Name() string

	// FieldDescriptors return the struct field descriptors
	FieldDescriptors() []FieldDescriptor

	// Scope return the scope name of current bean definition
	Scope() string

	// IsLazyMode return true if the bean is lazy init
	IsLazyMode() bool

	// IsPrimary return true if the bean is primary
	IsPrimary() bool

	// Constructor return the constructor for bean
	Constructor() Constructor
}

// MustNewBeanDefinition return the BeanDefinition impl, and it panics when some error happened.
func MustNewBeanDefinition(typ reflect.Type, opts ...BeanDefinitionOption) BeanDefinition {
	b, err := NewBeanDefinition(typ, opts...)
	if err != nil {
		panic(err)
	}

	return b
}

// NewBeanDefinition return the BeanDefinition impl.
func NewBeanDefinition(typ reflect.Type, opts ...BeanDefinitionOption) (BeanDefinition, error) {
	opt := defaultBeanDefinitionOption()
	for _, o := range opts {
		o.Apply(opt)
	}

	if err := opt.Validate(); err != nil {
		return nil, err
	}

	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, xerrors.Errorf("Cannot build bean definition from '%s', it must be *struct", typ.String())
	}

	beanName := opt.name
	if beanName == "" {
		// TODO: change to name generator
		beanName = typ.Elem().Name()
	}

	fieldDescriptors := []FieldDescriptor{}
	for idx := 0; idx < typ.Elem().NumField(); idx++ {
		field := typ.Elem().Field(idx)

		fd, err := NewFieldDescriptor(field, idx)
		if err != nil {
			return nil, err
		}
		if fd == nil {
			continue
		}

		fieldDescriptors = append(fieldDescriptors, *fd)
	}

	constructor, err := constructorBuilder.Build(typ, opt.construtorArguments)
	if err != nil {
		return nil, err
	}

	b := &beanDefinitionHolder{
		Typ:              typ,
		name:             beanName,
		fieldDescriptors: fieldDescriptors,
		scope:            opt.scope,
		lazy:             opt.lazy,
		primary:          opt.primary,
		constructor:      constructor,
	}
	return b, nil
}

type beanDefinitionHolder struct {
	Typ     reflect.Type
	name    string
	scope   string
	lazy    bool
	primary bool

	fieldDescriptors []FieldDescriptor
	constructor      Constructor
}

func (b *beanDefinitionHolder) Type() reflect.Type {
	return b.Typ
}

func (b *beanDefinitionHolder) Name() string {
	return b.name
}

func (b *beanDefinitionHolder) FieldDescriptors() []FieldDescriptor {
	return b.fieldDescriptors
}

func (b *beanDefinitionHolder) Scope() string {
	return b.scope
}

func (b *beanDefinitionHolder) IsLazyMode() bool {
	return b.lazy
}

func (b *beanDefinitionHolder) Constructor() Constructor {
	return b.constructor
}

func (b *beanDefinitionHolder) IsPrimary() bool {
	return b.primary
}
