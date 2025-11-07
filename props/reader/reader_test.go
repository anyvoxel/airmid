package reader

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestRegisterReader(t *testing.T) {
	t.Run("register success", func(t *testing.T) {
		readers = []Reader{}

		g := NewWithT(t)
		r1 := &extReader{}
		r2 := &extReader{}
		err := RegisterReader(r1)
		g.Expect(err).ToNot(HaveOccurred())

		err = RegisterReader(r2)
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("register duplicate", func(t *testing.T) {
		r1 := &extReader{}
		readers = []Reader{r1}

		g := NewWithT(t)

		err := RegisterReader(r1)
		g.Expect(err).To(HaveOccurred())
	})
}

func TestRegisterExtFileReader(t *testing.T) {
	t.Run("register success", func(t *testing.T) {
		readers = []Reader{}

		g := NewWithT(t)
		err := RegisterExtFileReader(nil, "1", "2")
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(readers).To(HaveLen(1))
		g.Expect(readers[0]).To(Equal(&extReader{
			exts: []string{"1", "2"},
		}))
	})

	t.Run("register empty ext", func(t *testing.T) {
		readers = []Reader{}

		g := NewWithT(t)
		err := RegisterExtFileReader(nil)
		g.Expect(err).To(HaveOccurred())
	})
}

func TestRead(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := NewWithT(t)

		v1 := map[string]any{
			"1": 2,
		}
		v2 := map[string]any{
			"3": 4,
		}
		readers = []Reader{
			&extReader{
				fn: func(data []byte) (map[string]any, error) {
					return v1, nil
				},
				exts: []string{"1", "2"},
			},
			&extReader{
				fn: func(data []byte) (map[string]any, error) {
					return v2, nil
				},
				exts: []string{"3", "4"},
			},
		}
		vv, err := Read("x4", []byte{})
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(vv).To(Equal(v2))
	})

	t.Run("no reader", func(t *testing.T) {
		g := NewWithT(t)

		v1 := map[string]any{
			"1": 2,
		}
		v2 := map[string]any{
			"3": 4,
		}
		readers = []Reader{
			&extReader{
				fn: func(data []byte) (map[string]any, error) {
					return v1, nil
				},
				exts: []string{"1", "2"},
			},
			&extReader{
				fn: func(data []byte) (map[string]any, error) {
					return v2, nil
				},
				exts: []string{"3", "4"},
			},
		}
		_, err := Read("x5", []byte{})
		g.Expect(err).To(HaveOccurred())
	})
}
