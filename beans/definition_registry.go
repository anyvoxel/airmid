package beans

import (
	"sync"

	"github.com/anyvoxel/airmid/xerrors"
)

// BeanDefinitionVisitor is the interface that visitor all bean definitions.
type BeanDefinitionVisitor interface {
	// Visit will be called when walk through all bean definitions
	Visit(beanName string, beanDefinition BeanDefinition)
}

// FuncVisitor will wraps func to visit all bean definitions.
type FuncVisitor struct {
	VisitFunc func(string, BeanDefinition)
}

// Visit will forwart the argument to user defined VisitFunc.
func (v FuncVisitor) Visit(beanName string, beanDefinition BeanDefinition) {
	v.VisitFunc(beanName, beanDefinition)
}

// BeanDefinitionRegistry is the interface that hold bean definitions.
type BeanDefinitionRegistry interface {
	// RegisterBeanDefinition register a new bean definition with this registry.
	// It will return err if BeanDefinition is invalid or beanName is already exists.
	RegisterBeanDefinition(beanName string, beanDefinition BeanDefinition) error

	// RemoveBeanDefinition remove the bean definition for the given name.
	// It will return err if there is no such bean definition.
	RemoveBeanDefinition(beanName string) error

	// GetBeanDefinition return the bean definition for the given bean name.
	// It will return err if there is no such bean definition.
	GetBeanDefinition(beanName string) (BeanDefinition, error)

	// VisitBeanDefinition will walk through all bean definitions
	VisitBeanDefinition(visitor BeanDefinitionVisitor)
}

// NewBeanDefinitionRegistry return the BeanDefinitionRegistry impl.
func NewBeanDefinitionRegistry() BeanDefinitionRegistry {
	return &beanDefinitionRegistryImpl{
		beanDefinitionMap: make(map[string]BeanDefinition),
	}
}

type beanDefinitionRegistryImpl struct {
	beanDefinitionMap map[string]BeanDefinition
	lock              sync.RWMutex
}

func (r *beanDefinitionRegistryImpl) RegisterBeanDefinition(
	beanName string,
	beanDefinition BeanDefinition,
) (err error) {
	r.withLock(func() {
		_, ok := r.beanDefinitionMap[beanName]
		if ok {
			err = xerrors.WrapDuplicate("Cannot register bean '%v': It is already registered", beanName)
			return
		}

		r.beanDefinitionMap[beanName] = beanDefinition
	})

	return err
}

func (r *beanDefinitionRegistryImpl) RemoveBeanDefinition(beanName string) (err error) {
	r.withLock(func() {
		_, ok := r.beanDefinitionMap[beanName]
		if !ok {
			err = xerrors.WrapNotFound("No bean '%v' registered", beanName)
			return
		}

		delete(r.beanDefinitionMap, beanName)
	})

	return err
}

func (r *beanDefinitionRegistryImpl) GetBeanDefinition(beanName string) (beanDefinition BeanDefinition, err error) {
	r.withRLock(func() {
		var ok bool
		beanDefinition, ok = r.beanDefinitionMap[beanName]
		if !ok {
			beanDefinition = nil
			err = xerrors.WrapNotFound("No bean '%v' registered", beanName)
			return
		}
	})

	return beanDefinition, err
}

func (r *beanDefinitionRegistryImpl) VisitBeanDefinition(visitor BeanDefinitionVisitor) {
	r.withRLock(func() {
		for k, v := range r.beanDefinitionMap {
			visitor.Visit(k, v)
		}
	})
}

func (r *beanDefinitionRegistryImpl) withRLock(fn func()) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	fn()
}

func (r *beanDefinitionRegistryImpl) withLock(fn func()) {
	r.lock.Lock()
	defer r.lock.Unlock()

	fn()
}
