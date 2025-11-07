package utils

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/xerrors"
)

func TestMust(t *testing.T) {
	type testCase struct {
		desp   string
		err    error
		expect error
	}
	testCases := []testCase{
		{
			desp:   "panic with err",
			err:    xerrors.ErrNotFound,
			expect: xerrors.ErrNotFound,
		},
		{
			desp:   "not panic",
			err:    nil,
			expect: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			if tc.expect != nil {
				g.Ω(func() {
					Must(tc.err)
				}).To(PanicWith(tc.expect))
			} else {
				g.Ω(func() {
					Must(tc.err)
				}).NotTo(Panic())
			}
		})
	}
}
