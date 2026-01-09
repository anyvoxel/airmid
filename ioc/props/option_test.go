package props

import (
	"reflect"
	"testing"

	"github.com/onsi/gomega"
)

const testString = "string"

func Test_defaultGetOption(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	opt := defaultGetOption()
	g.Expect(opt).ToNot(gomega.BeNil())
	g.Expect(opt.Target).To(gomega.Equal(reflect.Value{}))
	g.Expect(opt.Typ).To(gomega.Equal(reflect.TypeOf("")))
	g.Expect(opt.Default).To(gomega.BeNil())
}

func Test_getOption_Validate(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	t.Run("NilTyp", func(t *testing.T) {
		opt := &getOption{
			Typ: nil,
		}
		err := opt.Validate()
		g.Expect(err).To(gomega.HaveOccurred())
	})

	t.Run("TargetNotMatchTyp", func(t *testing.T) {
		s := testString
		opt := &getOption{
			Target: reflect.ValueOf(&s),
			Typ:    reflect.TypeOf(1),
		}
		err := opt.Validate()
		g.Expect(err).To(gomega.HaveOccurred())
	})

	t.Run("Valid", func(t *testing.T) {
		s := testString
		opt := &getOption{
			Target: reflect.ValueOf(&s),
			Typ:    reflect.TypeOf(&s),
		}
		err := opt.Validate()
		g.Expect(err).ToNot(gomega.HaveOccurred())
	})
}

func Test_getOption_Complete(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	t.Run("ValidateFail", func(t *testing.T) {
		opt := &getOption{
			Typ: nil,
		}
		err := opt.Complete()
		g.Expect(err).To(gomega.HaveOccurred())
	})

	t.Run("InvalidTargetNotPtr", func(t *testing.T) {
		opt := &getOption{
			Typ: reflect.TypeOf(""),
		}
		err := opt.Complete()
		g.Expect(err).ToNot(gomega.HaveOccurred())
		g.Expect(opt.IsTargetValid()).To(gomega.BeTrue())
		g.Expect(opt.Target.Type()).To(gomega.Equal(reflect.TypeOf("")))
	})

	t.Run("InvalidTargetPtr", func(t *testing.T) {
		opt := &getOption{
			Typ: reflect.TypeOf((*string)(nil)),
		}
		err := opt.Complete()
		g.Expect(err).ToNot(gomega.HaveOccurred())
		g.Expect(opt.IsTargetValid()).To(gomega.BeTrue())
		g.Expect(opt.Target.Type()).To(gomega.Equal(reflect.TypeOf((*string)(nil))))
	})

	t.Run("ValidTarget", func(t *testing.T) {
		s := testString
		opt := &getOption{
			Target: reflect.ValueOf(&s),
			Typ:    reflect.TypeOf(&s),
		}
		err := opt.Complete()
		g.Expect(err).ToNot(gomega.HaveOccurred())
		g.Expect(opt.Target).To(gomega.Equal(reflect.ValueOf(&s)))
	})
}

func TestOptions(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	t.Run("WithDefault", func(t *testing.T) {
		opt := &getOption{}
		WithDefault("default").Apply(opt)
		g.Expect(*opt.Default).To(gomega.Equal("default"))
	})

	t.Run("WithTarget", func(t *testing.T) {
		opt := &getOption{}
		s := testString
		WithTarget(&s).Apply(opt)
		g.Expect(opt.Target).To(gomega.Equal(reflect.ValueOf(&s)))
		g.Expect(opt.Typ).To(gomega.Equal(reflect.TypeOf(&s)))
	})

	t.Run("WithType", func(t *testing.T) {
		opt := &getOption{}
		typ := reflect.TypeOf(1)
		WithType(typ).Apply(opt)
		g.Expect(opt.Typ).To(gomega.Equal(typ))
	})
}
