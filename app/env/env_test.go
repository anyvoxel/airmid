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

package env

import (
	"context"
	"os"
	"testing"

	"github.com/onsi/gomega"
)

func TestLoaderImpl_Load(t *testing.T) {
	type testCase struct {
		desc       string
		newLoader  func() Loader
		testBefore func()
		testAfter  func()
		expected   map[string]string
	}

	testCases := []testCase{
		{
			desc: "test prefix",
			newLoader: func() Loader {
				return NewEnvLoader(WithPrefixOption("AIRMID_"))
			},
			testBefore: func() {
				_ = os.Setenv("AIRMID_TEST", "test")
				_ = os.Setenv("AIRMIDA_TEST", "test")
			},
			testAfter: func() {
				_ = os.Unsetenv("AIRMID_TEST")
				_ = os.Unsetenv("AIRMIDA_TEST")
			},
			expected: map[string]string{
				"TEST": "test",
			},
		},
		{
			desc: "test include",
			newLoader: func() Loader {
				return NewEnvLoader(WithPrefixOption("AIRMID_"), WithEnvIncludePattern("^*_INCLUDE$"))
			},
			testBefore: func() {
				_ = os.Setenv("AIRMID_TEST_INCLUDE", "test")
				_ = os.Setenv("AIRMID_TEST_EXCLUDE", "test")
			},
			testAfter: func() {
				_ = os.Unsetenv("AIRMID_TEST_INCLUDE")
				_ = os.Unsetenv("AIRMID_TEST_EXCLUDE")
			},
			expected: map[string]string{
				"TEST_INCLUDE": "test",
			},
		},
		{
			desc: "test exclude",
			newLoader: func() Loader {
				return NewEnvLoader(WithPrefixOption("AIRMID_"), WithEnvExcludePattern("^*_EXCLUDE$"))
			},
			testBefore: func() {
				_ = os.Setenv("AIRMID_TEST_INCLUDE", "test")
				_ = os.Setenv("AIRMID_TEST_EXCLUDE", "test")
			},
			testAfter: func() {
				_ = os.Unsetenv("AIRMID_TEST_INCLUDE")
				_ = os.Unsetenv("AIRMID_TEST_EXCLUDE")
			},
			expected: map[string]string{
				"TEST_INCLUDE": "test",
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.desc, func(t *testing.T) {
			g := gomega.NewWithT(t)
			if c.testBefore != nil {
				c.testBefore()
			}
			loader := c.newLoader()
			actual := loader.Load(context.Background())
			g.Expect(actual).Should(gomega.Equal(c.expected))
			if c.testAfter != nil {
				c.testAfter()
			}
		})
	}
}
