package ioc

import (
	"reflect"
	"testing"

	. "github.com/onsi/gomega"
)

func TestDefaultBeanDefinitionOption(t *testing.T) {
	g := NewWithT(t)

	opt := defaultBeanDefinitionOption()
	g.Expect(opt).ToNot(BeNil())
	g.Expect(opt.scope).To(Equal(ScopeSingleton))
	g.Expect(opt.lazy).To(BeFalse())
	g.Expect(opt.primary).To(BeFalse())
	g.Expect(opt.name).To(BeEmpty())
	g.Expect(opt.construtorArguments).To(BeEmpty())
}

func TestWithBeanName(t *testing.T) {
	g := NewWithT(t)

	opt := &beanDefinitionOption{}
	WithBeanName("testBean").Apply(opt)
	g.Expect(opt.name).To(Equal("testBean"))
}

func TestWithBeanScope(t *testing.T) {
	g := NewWithT(t)

	t.Run("SingletonScope", func(t *testing.T) {
		opt := &beanDefinitionOption{}
		WithBeanScope(ScopeSingleton).Apply(opt)
		g.Expect(opt.scope).To(Equal(ScopeSingleton))
	})

	t.Run("PrototypeScope", func(t *testing.T) {
		opt := &beanDefinitionOption{}
		WithBeanScope(ScopePrototype).Apply(opt)
		g.Expect(opt.scope).To(Equal(ScopePrototype))
	})
}

func TestWithLazyMode(t *testing.T) {
	g := NewWithT(t)

	opt := &beanDefinitionOption{}
	WithLazyMode().Apply(opt)
	g.Expect(opt.lazy).To(BeTrue())
}

func TestWithPrimary(t *testing.T) {
	g := NewWithT(t)

	opt := &beanDefinitionOption{}
	WithPrimary().Apply(opt)
	g.Expect(opt.primary).To(BeTrue())
}

func TestWithConstructorArguments(t *testing.T) {
	g := NewWithT(t)

	args := []ConstructorArgument{
		{
			Type: reflect.TypeOf(""),
		},
	}
	opt := &beanDefinitionOption{}
	WithConstructorArguments(args).Apply(opt)
	g.Expect(opt.construtorArguments).To(Equal(args))
}

func TestBeanDefinitionOptionValidate(t *testing.T) {
	g := NewWithT(t)

	t.Run("ValidScopes", func(t *testing.T) {
		opt := &beanDefinitionOption{scope: ScopeSingleton}
		g.Expect(opt.Validate()).ToNot(HaveOccurred())

		opt = &beanDefinitionOption{scope: ScopePrototype}
		g.Expect(opt.Validate()).ToNot(HaveOccurred())
	})

	t.Run("InvalidScope", func(t *testing.T) {
		opt := &beanDefinitionOption{scope: "invalid"}
		err := opt.Validate()
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("Unsupport scope 'invalid'"))
	})
}
