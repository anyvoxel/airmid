package xapp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/xerrors"
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
