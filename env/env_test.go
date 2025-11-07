package env

import (
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
			actual := loader.Load()
			g.Expect(actual).Should(gomega.Equal(c.expected))
			if c.testAfter != nil {
				c.testAfter()
			}
		})
	}
}
