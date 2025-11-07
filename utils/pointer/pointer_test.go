package pointer

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestStringPtr(t *testing.T) {
	g := NewWithT(t)

	v := StringPtr("1")
	expect := "1"
	g.Expect(v).To(Equal(&expect))
}

func TestIntPtr(t *testing.T) {
	g := NewWithT(t)

	v := IntPtr(1)
	expect := 1
	g.Expect(v).To(Equal(&expect))
}
