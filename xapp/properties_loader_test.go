package xapp

import (
	"os"
	"testing"

	"github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/props"
)

func TestEnvPropertiesLoader_LoadProperties(t *testing.T) {
	type testCase struct {
		desc       string
		expected   map[string]string
		testBefore func()
		testAfter  func()
		isErr      bool
		errMsg     string
	}
	testCases := []testCase{
		{
			desc: "no include and exclude",
			testBefore: func() {
				_ = os.Unsetenv("AIRMID_INCLUDE_ENV_PATTERNS")
				_ = os.Unsetenv("AIRMID_EXCLUDE_ENV_PATTERNS")

				_ = os.Setenv("AIRMID_TEST", "test")
			},
			testAfter: func() {
				_ = os.Unsetenv("AIRMID_TEST")
			},
			isErr: false,
			expected: map[string]string{
				"test": "test",
			},
		},
		{
			desc: "include, no exclude",
			testBefore: func() {
				_ = os.Unsetenv("AIRMID_EXCLUDE_ENV_PATTERNS")

				_ = os.Setenv("AIRMID_INCLUDE_ENV_PATTERNS", "^TEST*")
				_ = os.Setenv("AIRMID_TEST_ENV", "test.env")
			},
			testAfter: func() {
				_ = os.Unsetenv("AIRMID_TEST_ENV")
				_ = os.Unsetenv("AIRMID_INCLUDE_ENV_PATTERNS")
			},
			isErr: false,
			expected: map[string]string{
				"test.env": "test.env",
			},
		},
		{
			desc: "include and exclude",
			testBefore: func() {
				_ = os.Setenv("AIRMID_INCLUDE_ENV_PATTERNS", "^TEST*")
				_ = os.Setenv("AIRMID_EXCLUDE_ENV_PATTERNS", "^*_EXCLUDE")
				_ = os.Setenv("AIRMID_TEST_ENV_EXCLUDE", "test.env.exclude")
			},
			testAfter: func() {
				_ = os.Unsetenv("AIRMID_TEST_ENV")
				_ = os.Unsetenv("AIRMID_INCLUDE_ENV_PATTERNS")
				_ = os.Unsetenv("AIRMID_EXCLUDE_ENV_PATTERNS")
			},
			isErr:  true,
			errMsg: "property with key='test.env.exclude' not found: ObjectNotFound",
			expected: map[string]string{
				"test.env.exclude": "test.env.exclude",
			},
		},
		{
			desc: "test array value",
			testBefore: func() {
				_ = os.Setenv("AIRMID_TEST_ENV", "test0,test1,test2")
			},
			testAfter: func() {
				_ = os.Unsetenv("AIRMID_TEST_ENV")
			},
			isErr: false,
			expected: map[string]string{
				"test.env[0]": "test0",
				"test.env[1]": "test1",
				"test.env[2]": "test2",
			},
		},
	}

	for _, tc := range testCases {
		g := gomega.NewWithT(t)
		t.Run(tc.desc, func(t *testing.T) {
			if tc.testBefore != nil {
				tc.testBefore()
			}
			p := props.NewProperties()
			loader := NewEnvPropertiesLoader("AIRMID_", DefaultEnvKeyConvertFunc)
			_ = loader.LoadProperties(p)
			for k, v := range tc.expected {
				actual, err := p.Get(k)
				if tc.isErr {
					g.Expect(err).Should(gomega.HaveOccurred())
					g.Expect(err.Error()).Should(gomega.Equal(tc.errMsg))
					return
				}
				g.Expect(err).ShouldNot(gomega.HaveOccurred())
				g.Expect(actual).To(gomega.Equal(v))
			}
			if tc.testAfter != nil {
				tc.testAfter()
			}
		})
	}
}

func TestOptionArgsPropertiesLoader_LoadProperties(t *testing.T) {
	type testCase struct {
		desc       string
		testBefore func(p props.Properties)
		testAfter  func()
		expect     map[string]string
		isErr      bool
	}

	testCases := []testCase{
		{
			desc: "Read from args",
			testBefore: func(_ props.Properties) {
				os.Args = []string{"", "--test=test", "-ok"}
			},
			expect: map[string]string{
				"test": "test",
				"ok":   "true",
			},
			isErr: false,
		},
		{
			desc: "Read from config, env and args  with same properties",
			testBefore: func(properties props.Properties) {
				_ = os.Setenv("AIRMID_TEST", "test_value")
				envLoader := NewEnvPropertiesLoader("AIRMID_", DefaultEnvKeyConvertFunc)
				_ = envLoader.LoadProperties(properties)
				os.Args = []string{"", "--test=test"}
			},
			testAfter: func() {
				_ = os.Unsetenv("AIRMID_TEST")
			},
			expect: map[string]string{
				"test": "test",
			},
			isErr: false,
		},
		{
			desc: "Read from env and args with different properties",
			testBefore: func(properties props.Properties) {
				_ = os.Setenv("AIRMID_TEST1", "test1")
				envLoader := NewEnvPropertiesLoader("AIRMID_", DefaultEnvKeyConvertFunc)
				_ = envLoader.LoadProperties(properties)

				os.Args = []string{"", "--test2=test2"}
			},
			testAfter: func() {
				_ = os.Unsetenv("AIRMID_TEST1")
			},
			expect: map[string]string{
				"test1": "test1",
				"test2": "test2",
			},
			isErr: false,
		},
		{
			desc: "Read from array flag",
			testBefore: func(_ props.Properties) {
				os.Args = []string{"", "--test.flag=test0,test1,test2"}
			},
			expect: map[string]string{
				"test.flag[0]": "test0",
				"test.flag[1]": "test1",
				"test.flag[2]": "test2",
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.desc, func(t *testing.T) {
			g := gomega.NewWithT(t)
			p := props.NewProperties()
			if c.testBefore != nil {
				c.testBefore(p)
			}
			loader := NewOptionArgsPropertiesLoader()
			err := loader.LoadProperties(p)
			if c.isErr {
				g.Expect(err).Should(gomega.HaveOccurred())
				return
			}
			for k, v := range c.expect {
				g.Expect(p.Get(k)).Should(gomega.Equal(v))
			}
			if c.testAfter != nil {
				c.testAfter()
			}
		})
	}
}

func TestPropertiesLoadersImpl_Add(t *testing.T) {
	g := gomega.NewWithT(t)

	envLoader := NewEnvPropertiesLoader("", nil)
	loaders := &propertiesLoaders{}
	loaders.Add(envLoader)

	g.Expect(loaders.loaders).Should(gomega.HaveLen(1))
}

func TestPropertiesLoadersImpl_LoadProperties(t *testing.T) {
	g := gomega.NewWithT(t)
	t.Setenv("AIRMID_TEST1", "test1")

	p := props.NewProperties()
	envLoader := NewEnvPropertiesLoader("AIRMID_", DefaultEnvKeyConvertFunc)
	loaders := NewPropertiesLoaders()
	err := loaders.Add(envLoader).LoadProperties(p)

	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	v, _ := p.Get("test1")
	g.Expect(v).Should(gomega.Equal("test1"))
}
