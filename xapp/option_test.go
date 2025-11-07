package xapp

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestWithAttributes(t *testing.T) {
	g := NewWithT(t)

	o := newOption([]Option{
		WithAttributes(Attribute{
			Key:   "k1",
			Value: "v1",
		}),
		WithAttributes(
			Attribute{
				Key:   "k2",
				Value: "v2",
			},
			Attribute{
				Key:   "k1",
				Value: "v3",
			},
		),
	})
	g.Expect(o).To(Equal(&option{
		attrs: []Attribute{
			{
				Key:   "k1",
				Value: "v1",
			},
			{
				Key:   "k2",
				Value: "v2",
			},
			{
				Key:   "k1",
				Value: "v3",
			},
		},
	}))
}
