package reader

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestExtReaderRead(t *testing.T) {
	g := NewWithT(t)

	vv := map[string]any{
		"1": 1,
	}
	r := &extReader{
		fn: func(data []byte) (map[string]any, error) {
			return vv, nil
		},
	}
	g.Expect(r.Read([]byte{})).To(Equal(vv))
}

func TestExtReaderMatch(t *testing.T) {
	type testCase struct {
		desp     string
		r        *extReader
		filename string
		err      string
	}
	testCases := []testCase{
		{
			desp: "match",
			r: &extReader{
				exts: []string{"1", "2"},
			},
			filename: "x2",
			err:      "",
		},
		{
			desp: "match",
			r: &extReader{
				exts: []string{"1", "2"},
			},
			filename: "x",
			err:      "'1,2' cannot support filename 'x'",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			err := tc.r.Match(tc.filename)
			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
		})
	}
}

func TestExtReaderName(t *testing.T) {
	g := NewWithT(t)
	r := &extReader{
		exts: []string{"1", "2"},
	}
	g.Expect(r.Name()).To(Equal("1,2"))
}
