package utils

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/xerrors"
)

func TestSafeRun(t *testing.T) {
	type testCase struct {
		desp string
		cmd  func()
	}
	testCases := []testCase{
		{
			desp: "normal runnable",
			cmd:  func() {},
		},
		{
			desp: "panic on error",
			cmd: func() {
				panic(xerrors.ErrNotFound)
			},
		},
		{
			desp: "panic on int",
			cmd: func() {
				panic(1)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)

			g.Ω(func() {
				SafeRun(tc.cmd)
			}).ToNot(Panic())
		})
	}
}
