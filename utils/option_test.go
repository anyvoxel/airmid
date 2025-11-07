package utils

import (
	"testing"

	. "github.com/onsi/gomega"
)

type testOption struct {
	v int
}

func TestFnOptionApply(t *testing.T) {
	g := NewWithT(t)

	o := NewFnOption(func(opt *testOption) {
		opt.v = 1
	})
	opt := &testOption{}
	o.Apply(opt)

	g.Expect(opt.v).To(Equal(1))
}
