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
