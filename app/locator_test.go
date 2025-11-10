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

package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/anvil/xerrors"
)

func TestLocalResourceLocatorLocate(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := NewWithT(t)

		l := &localResourceLocator{
			configDir: []string{"1", "2"},
		}

		guard := gomonkey.ApplyFunc(os.Open, func(string) (*os.File, error) {
			return &os.File{}, nil
		})
		defer guard.Reset()

		rr, err := l.Locate("f1")

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(rr).To(HaveLen(2))
	})

	t.Run("abs failed", func(t *testing.T) {
		g := NewWithT(t)
		l := &localResourceLocator{
			configDir: []string{"1", "2"},
		}

		guard := gomonkey.ApplyFunc(filepath.Abs, func(string) (string, error) {
			return "", xerrors.ErrContinue
		})
		defer guard.Reset()

		rr, err := l.Locate("f1")
		g.Expect(err).To(HaveOccurred())
		g.Expect(rr).To(BeNil())
	})

	t.Run("with not found", func(t *testing.T) {
		g := NewWithT(t)

		l := &localResourceLocator{
			configDir: []string{"1", "2"},
		}

		count := 0
		guard := gomonkey.ApplyFunc(os.Open, func(string) (*os.File, error) {
			if count == 0 {
				count++
				return &os.File{}, nil
			}
			return nil, os.ErrNotExist
		})
		defer guard.Reset()

		rr, err := l.Locate("f1")

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(rr).To(HaveLen(1))
	})

	t.Run("open failed", func(t *testing.T) {
		g := NewWithT(t)

		l := &localResourceLocator{
			configDir: []string{"1", "2"},
		}

		guard := gomonkey.ApplyFunc(os.Open, func(string) (*os.File, error) {
			return nil, xerrors.ErrContinue
		})
		defer guard.Reset()

		rr, err := l.Locate("f1")

		g.Expect(err).To(HaveOccurred())
		g.Expect(rr).To(BeNil())
	})
}
