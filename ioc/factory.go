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

// Package ioc provide the implement to build Bean object
package ioc

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sort"
	"sync"

	"github.com/anyvoxel/airmid/anvil/logger"
	"github.com/anyvoxel/airmid/anvil/parallel"
	"github.com/anyvoxel/airmid/anvil/xerrors"
	"github.com/anyvoxel/airmid/ioc/props"
)

// BeanFactory providing the full capabilities of SPI.
type BeanFactory interface {
	// GetBean return an instance, which may be shared or independent, of the specified bean.
	GetBean(name string) (any, error)

	// ResolveBeanNames return the bean names of the specified bean type.
	ResolveBeanNames(typ reflect.Type) ([]string, error)

	// RegisterScope will register scope with name to factory
	RegisterScope(name string, scope Scope) error

	// RegisterSingleton will register the given existing object as singleton in the bean registry
	// under the given bean name.
	// The object is supposed to be fully initialized.
	RegisterSingleton(name string, bean any) error

	// AddBeanPostProcessor will add beanPostProcessor to factory.
	AddBeanPostProcessor(beanPostProcessor BeanPostProcessor)

	// PreInstantiateSingletons will pre initializing the non-lazy mode singletons
	PreInstantiateSingletons() error

	// Destroy will destroy the beans before bean destruction.
	Destroy()

	BeanDefinitionRegistry
	props.Properties
}

// NewBeanFactory return the BeanFactory impl.
func NewBeanFactory() BeanFactory {
	beanFactory := &beanFactoryImpl{
		BeanDefinitionRegistry: NewBeanDefinitionRegistry(),
		Properties:             props.NewProperties(),
		singletonObjects:       make(map[string]reflect.Value),
		allBeans:               make(map[string]any),
		beansInCreating:        make(map[string]any),
		scopes:                 make(map[string]Scope),
	}
	beanFactory.beanPostProcessorCompositor = NewBeanPostProcessorCompositor(&BeanAwareProcessor{
		beanFactory: beanFactory,
	})

	return beanFactory
}

type beanFactoryImpl struct {
	BeanDefinitionRegistry
	props.Properties

	// singletonObjects is the cache for singleton scope instance
	singletonObjects map[string]reflect.Value
	// all created beans
	allBeans        map[string]any
	beansInCreating map[string]any
	scopes          map[string]Scope

	beanPostProcessorCompositor BeanPostProcessorCompositor

	mu sync.RWMutex
}

func (f *beanFactoryImpl) GetBean(name string) (any, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.getBeanLocked(name)
}

//nolint:revive,cyclop
func (f *beanFactoryImpl) getBeanLocked(name string) (any, error) {
	if obj, ok := f.singletonObjects[name]; ok {
		return obj.Interface(), nil
	}

	beanDefinition, err := f.GetBeanDefinition(name)
	if err != nil {
		return nil, err
	}

	switch beanDefinition.Scope() {
	case ScopeSingleton, ScopePrototype:
	default:
		return nil, xerrors.Errorf(
			"BeanDefinition %s has invalid scope '%v', it must be ['%v', '%v']",
			beanDefinition.Name(), ScopeSingleton, ScopePrototype)
	}

	v, err := beanDefinition.Constructor().NewObject(&factoryConstructorArgumentResolver{
		bf: f,
	})
	if err != nil {
		return nil, err
	}

	cachedBean, ok := f.beansInCreating[name]
	if ok {
		switch beanDefinition.Scope() {
		case ScopePrototype:
			return nil, xerrors.Errorf("cannot get bean '%s' circularly", name)
		case ScopeSingleton:
			return cachedBean, nil
		}
	}
	f.beansInCreating[name] = v.Interface()
	defer delete(f.beansInCreating, name)

	err = f.wireStruct(v, beanDefinition.FieldDescriptors())
	if err != nil {
		return nil, err
	}

	// TODO: optimize this. When A depends B and B depends A,
	// if B's post process and initialization depends on A,
	// A's properties is not set, something unexpected will happen.
	f.beanPostProcessorCompositor.PostProcessBeanDefinition(name, beanDefinition)

	obj := v.Interface()

	if obj, err = f.beanPostProcessorCompositor.PostProcessBeforeInitialization(obj, name); err != nil {
		return nil, err
	}

	if vobj := IndirectTo[InitializingBean](obj); vobj != nil {
		logger.FromContext(context.TODO()).DebugContext(
			context.TODO(),
			"bean impelemented the InitializingBean interface, will execute it",
			slog.String("BeanName", name),
		)

		err := vobj.AfterPropertiesSet()
		if err != nil {
			return nil, err
		}
	} else {
		logger.FromContext(context.TODO()).DebugContext(
			context.TODO(),
			"bean doesn't impelemented the InitializingBean interface, will skip it",
			slog.String("BeanName", name),
		)
	}

	//nolint
	if obj, err = f.beanPostProcessorCompositor.PostProcessAfterInitialization(obj, name); err != nil {
		return nil, err
	}

	if beanDefinition.Scope() == ScopeSingleton {
		f.singletonObjects[name] = v
	}

	f.allBeans[name] = v.Interface()

	return v.Interface(), nil
}

// wireStruct will inject the field value into bean.
// 1. the bean must be pointer to struct.
// 2. it will return err when any inject failed.
func (f *beanFactoryImpl) wireStruct(bean reflect.Value, fds []FieldDescriptor) error {
	if len(fds) == 0 {
		return nil
	}

	propertyValues := NewPropertyValues()
	for _, fd := range fds {
		fn := func(_ FieldDescriptor, _ PropertyValues) error {
			return xerrors.ErrNotImplement
		}

		switch {
		case fd.Property != nil:
			fn = f.getPropertyValue
		case fd.Bean != nil:
			fn = f.getBeanValue
		}

		if err := fn(fd, propertyValues); err != nil {
			return err
		}
	}

	return propertyValues.SetProperty(bean.Elem(), fds)
}

func (f *beanFactoryImpl) getPropertyValue(fd FieldDescriptor, propertyValues PropertyValues) error {
	if fd.Property == nil {
		return nil
	}

	opts := []props.GetOption{}
	if fd.Property.Default != nil {
		opts = append(opts, props.WithDefault(*fd.Property.Default))
	}
	typ := fd.Typ
	if fd.Typ.Kind() != reflect.Ptr {
		// If the typ is not pointer, we change th target type to pointer,
		// so the value will be **singleton**
		typ = reflect.PointerTo(typ)
	}

	opts = append(opts, props.WithType(typ))
	value, err := f.Get(fd.Property.Name, opts...)
	if err != nil {
		return err
	}

	propertyValue := reflect.ValueOf(value)
	if fd.Typ.Kind() != reflect.Ptr {
		propertyValue = propertyValue.Elem()
	}
	propertyValues.AddValue(fd.FieldIndex, propertyValue)
	return nil
}

func (f *beanFactoryImpl) getBeanValue(fd FieldDescriptor, propertyValues PropertyValues) error {
	if fd.Bean == nil {
		return nil
	}

	if fd.Bean.Name != "?" {
		return f.getNamedBeanValue(fd, propertyValues)
	}

	return f.getTypedBeanValue(fd, propertyValues)
}

func (f *beanFactoryImpl) getNamedBeanValue(fd FieldDescriptor, propertyValues PropertyValues) error {
	obj, err := f.getBeanLocked(fd.Bean.Name)
	if err != nil {
		if xerrors.IsNotFound(err) && fd.Bean.Optional {
			return nil
		}
		return err
	}

	propertyValues.AddValue(fd.FieldIndex, reflect.ValueOf(obj))
	return nil
}

func (f *beanFactoryImpl) getTypedBeanValue(fd FieldDescriptor, propertyValues PropertyValues) error {
	primaryBeans, beans, err := f.getTypedCandidatesBeanNames(fd)
	if err != nil {
		return err
	}

	if fd.Typ.Kind() != reflect.Slice {
		return f.getTypedBeanNoneSliceValue(primaryBeans, beans, fd, propertyValues)
	}
	return f.getTypedBeanSliceValue(beans, fd, propertyValues)
}

func (f *beanFactoryImpl) getTypedBeanNoneSliceValue(
	primaryBeans []string, beans []string, fd FieldDescriptor, propertyValues PropertyValues) error {
	if len(primaryBeans) == 1 {
		beanValue, err := f.getBeanLockedAsValue(primaryBeans[0])
		if err != nil {
			return err
		}
		propertyValues.AddValue(fd.FieldIndex, beanValue)
		return nil
	}

	if len(primaryBeans) > 1 {
		return xerrors.Errorf(
			"'%v' primary candidates found for field '%v' with type %s",
			len(primaryBeans), fd.Name, fd.Typ.String())
	}

	if len(beans) == 0 {
		if !fd.Bean.Optional {
			// if no candidate beans and the field not optional, return error
			return xerrors.Errorf("No candidate found for field '%v' with type %v", fd.Name, fd.Typ.String())
		}
		return nil
	}

	if len(beans) > 1 {
		return xerrors.Errorf("'%v' candidates found for field '%v' with type %s", len(beans), fd.Name, fd.Typ.String())
	}

	beanValue, err := f.getBeanLockedAsValue(beans[0])
	if err != nil {
		return err
	}
	propertyValues.AddValue(fd.FieldIndex, beanValue)
	return nil
}

func (f *beanFactoryImpl) getTypedBeanSliceValue(
	beanNames []string, fd FieldDescriptor, propertyValues PropertyValues) error {
	candidates := make(CandidateBeans, 0, len(beanNames))
	for _, beanName := range beanNames {
		bean, err := f.getBeanLockedAsValue(beanName)
		if err != nil {
			return err
		}
		candidates = append(candidates, bean)
	}

	sort.Sort(candidates)
	ret := reflect.MakeSlice(fd.Typ, 0, len(candidates))
	for _, bean := range candidates {
		ret = reflect.Append(ret, bean)
	}

	propertyValues.AddValue(fd.FieldIndex, ret)
	return nil
}

func (f *beanFactoryImpl) getTypedCandidatesBeanNames(fd FieldDescriptor) ([]string, []string, error) {
	elemType := fd.Typ
	if fd.Typ.Kind() == reflect.Slice {
		elemType = elemType.Elem()
	}

	beanNames, err := f.ResolveBeanNames(elemType)
	if err != nil {
		return nil, nil, err
	}

	primaryBeans := make([]string, 0, len(beanNames))
	beans := make([]string, 0, len(beanNames))
	for _, beanName := range beanNames {
		def, err := f.GetBeanDefinition(beanName)
		if err != nil {
			return nil, nil, err
		}

		if def.IsPrimary() {
			primaryBeans = append(primaryBeans, beanName)
		}
		beans = append(beans, beanName)
	}
	return primaryBeans, beans, nil
}

func (f *beanFactoryImpl) getBeanLockedAsValue(beanName string) (reflect.Value, error) {
	obj, err := f.getBeanLocked(beanName)
	if err != nil {
		return reflect.Value{}, err
	}

	return reflect.ValueOf(obj), nil
}

func (f *beanFactoryImpl) ResolveBeanNames(typ reflect.Type) ([]string, error) {
	beanNames := []string{}
	f.VisitBeanDefinition(FuncVisitor{
		VisitFunc: func(s string, bd BeanDefinition) {
			if IsTypeMatched(typ, bd.Type()) {
				beanNames = append(beanNames, s)
			}
		},
	})

	return beanNames, nil
}

func (f *beanFactoryImpl) RegisterScope(name string, scope Scope) error {
	if name == ScopeSingleton || name == ScopePrototype {
		return xerrors.Errorf(
			"Invalid scope name '%s', cannot replace existing scopes '%s' and '%s'",
			name, ScopePrototype, ScopeSingleton)
	}

	pre, ok := f.scopes[name]
	if ok && pre != scope {
		logger.FromContext(context.TODO()).DebugContext(
			context.TODO(),
			"Old scope is exists, replicing it when new scope is register",
			slog.String("ScopeName", name),
			slog.String("OldScopeType", fmt.Sprintf("%v", pre)),
			slog.String("CurrentScopeType", fmt.Sprintf("%v", scope)),
		)
	}
	f.scopes[name] = scope
	return nil
}

func (f *beanFactoryImpl) RegisterSingleton(name string, bean any) error {
	oldObject, ok := f.singletonObjects[name]
	if ok {
		return xerrors.Errorf(
			"Invalid object '%v' under bean name '%v', it already exists with object '%v'",
			bean, name, oldObject.Interface())
	}

	// TODO: add more validate for bean
	f.singletonObjects[name] = reflect.ValueOf(bean)
	return nil
}

func (f *beanFactoryImpl) AddBeanPostProcessor(beanPostProcessor BeanPostProcessor) {
	f.beanPostProcessorCompositor.AddBeanPostProcessor(beanPostProcessor)
}

func (f *beanFactoryImpl) PreInstantiateSingletons() error {
	beanNames := []string{}
	f.VisitBeanDefinition(FuncVisitor{
		VisitFunc: func(s string, bd BeanDefinition) {
			if bd.Scope() == ScopeSingleton && !bd.IsLazyMode() {
				beanNames = append(beanNames, s)
			}
		},
	})

	if err := f.doConcurrentInstantiateSingleton(beanNames); err != nil {
		return err
	}

	for name, obj := range f.singletonObjects {
		if smartSingleton := IndirectTo[SmartInitializingSingleton](obj.Interface()); smartSingleton != nil {
			logger.FromContext(context.TODO()).DebugContext(
				context.TODO(),
				"singleton bean has implements SmartInitializingSingleton, will execute it",
				slog.String("BeanName", name),
			)

			if err := smartSingleton.AfterSingletonsInstantiated(); err != nil {
				logger.FromContext(context.TODO()).ErrorContext(
					context.TODO(),
					"Smart instantiate singleton bean failed",
					slog.String("BeanName", name),
					slog.Any("Error", err),
				)
				return err
			}
		} else {
			logger.FromContext(context.TODO()).DebugContext(
				context.TODO(),
				"singleton bean has not implements SmartInitializingSingleton, will skip it",
				slog.String("BeanName", name),
			)
		}
	}
	return nil
}

func (f *beanFactoryImpl) doConcurrentInstantiateSingleton(beanNames []string) error {
	if len(beanNames) == 0 {
		return nil
	}

	return parallel.Run(context.TODO(), len(beanNames), func(idx int) error {
		_, err := f.GetBean(beanNames[idx])
		if err != nil {
			logger.FromContext(context.TODO()).ErrorContext(
				context.TODO(),
				"instantiate singleton bean failed",
				slog.String("BeanName", beanNames[idx]),
				slog.Any("Error", err),
			)
			return err
		}

		logger.FromContext(context.TODO()).InfoContext(
			context.TODO(),
			"instantiate singleton bean",
			slog.String("BeanName", beanNames[idx]),
		)
		return nil
	})
}

func (f *beanFactoryImpl) Destroy() {
	for name, bean := range f.allBeans {
		if obj := IndirectTo[DestructionAwareBeanPostProcessor](bean); obj != nil {
			obj.PostProcessBeforeDestruction(name, bean)
		}
	}
}
