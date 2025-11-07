package beans

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/anyvoxel/airmid/utils/pointer"
	"github.com/anyvoxel/airmid/xerrors"
)

func TestNewBeanPostProcessorCompositor(t *testing.T) {
	g := NewWithT(t)

	p1 := &FuncBeanPostProcessor{}
	p2 := &FuncBeanDefinitionPostProcessor{}
	p3 := &FuncBeanDefinitionPostProcessor{}

	p := NewBeanPostProcessorCompositor(p1, p2, p2, p1, p3).(*beanPostProcessorCompositorImpl)
	g.Expect(p).ToNot(BeNil())

	g.Expect(p.beanPostProcessors).To(Equal([]BeanPostProcessor{p2, p1, p3}))
	g.Expect(p.beanDefinitionPostProcessors).To(Equal([]BeanDefinitionPostProcessor{p2, p3}))
}

func TestPostProcessBeforeInitialization(t *testing.T) {
	type testCase struct {
		desp string
		p    *beanPostProcessorCompositorImpl
		init func(tc *testCase)

		obj    any
		err    string
		expect any
	}
	testCases := []testCase{
		{
			desp: "all pass",
			p:    NewBeanPostProcessorCompositor().(*beanPostProcessorCompositorImpl),
			init: func(tc *testCase) {
				tc.p.AddBeanPostProcessor(
					&FuncBeanPostProcessor{
						BeforeInitializationFunc: func(obj any, beanName string) (v any, err error) {
							vv := obj.(*int)
							*vv++
							return obj, nil
						},
					},
				)
				tc.p.AddBeanPostProcessor(
					&FuncBeanPostProcessor{
						BeforeInitializationFunc: func(obj any, beanName string) (v any, err error) {
							vv := obj.(*int)
							*vv += 2
							return obj, nil
						},
					},
				)
			},
			obj:    pointer.IntPtr(0),
			err:    "",
			expect: pointer.IntPtr(3),
		},
		{
			desp: "failed at 1",
			p:    NewBeanPostProcessorCompositor().(*beanPostProcessorCompositorImpl),
			init: func(tc *testCase) {
				tc.p.AddBeanPostProcessor(
					&FuncBeanPostProcessor{
						BeforeInitializationFunc: func(obj any, beanName string) (v any, err error) {
							return nil, xerrors.Errorf(tc.desp)
						},
					},
				)
				tc.p.AddBeanPostProcessor(
					&FuncBeanPostProcessor{
						BeforeInitializationFunc: func(obj any, beanName string) (v any, err error) {
							vv := obj.(*int)
							*vv += 2
							return obj, nil
						},
					},
				)
			},
			obj:    pointer.IntPtr(0),
			err:    "failed at 1",
			expect: nil,
		},
		{
			desp: "failed at 2",
			p:    NewBeanPostProcessorCompositor().(*beanPostProcessorCompositorImpl),
			init: func(tc *testCase) {
				tc.p.AddBeanPostProcessor(
					&FuncBeanPostProcessor{
						BeforeInitializationFunc: func(obj any, beanName string) (v any, err error) {
							vv := obj.(*int)
							*vv++
							return obj, nil
						},
					},
				)
				tc.p.AddBeanPostProcessor(
					&FuncBeanPostProcessor{
						BeforeInitializationFunc: func(obj any, beanName string) (v any, err error) {
							return nil, xerrors.Errorf(tc.desp)
						},
					},
				)
			},
			obj:    pointer.IntPtr(0),
			err:    "failed at 2",
			expect: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			tc.init(&tc)

			actual, err := tc.p.PostProcessBeforeInitialization(tc.obj, tc.desp)

			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestPostProcessAfterInitialization(t *testing.T) {
	type testCase struct {
		desp string
		p    *beanPostProcessorCompositorImpl
		init func(tc *testCase)

		obj    any
		err    string
		expect any
	}
	testCases := []testCase{
		{
			desp: "all pass",
			p:    NewBeanPostProcessorCompositor().(*beanPostProcessorCompositorImpl),
			init: func(tc *testCase) {
				tc.p.AddBeanPostProcessor(
					&FuncBeanPostProcessor{
						AfterInitializationFunc: func(obj any, beanName string) (v any, err error) {
							vv := obj.(*int)
							*vv++
							return obj, nil
						},
					},
				)
				tc.p.AddBeanPostProcessor(
					&FuncBeanPostProcessor{
						AfterInitializationFunc: func(obj any, beanName string) (v any, err error) {
							vv := obj.(*int)
							*vv += 2
							return obj, nil
						},
					},
				)
			},
			obj:    pointer.IntPtr(0),
			err:    "",
			expect: pointer.IntPtr(3),
		},
		{
			desp: "failed at 1",
			p:    NewBeanPostProcessorCompositor().(*beanPostProcessorCompositorImpl),
			init: func(tc *testCase) {
				tc.p.AddBeanPostProcessor(
					&FuncBeanPostProcessor{
						AfterInitializationFunc: func(obj any, beanName string) (v any, err error) {
							return nil, xerrors.Errorf(tc.desp)
						},
					},
				)
				tc.p.AddBeanPostProcessor(
					&FuncBeanPostProcessor{
						AfterInitializationFunc: func(obj any, beanName string) (v any, err error) {
							vv := obj.(*int)
							*vv += 2
							return obj, nil
						},
					},
				)
			},
			obj:    pointer.IntPtr(0),
			err:    "failed at 1",
			expect: nil,
		},
		{
			desp: "failed at 2",
			p:    NewBeanPostProcessorCompositor().(*beanPostProcessorCompositorImpl),
			init: func(tc *testCase) {
				tc.p.AddBeanPostProcessor(
					&FuncBeanPostProcessor{
						AfterInitializationFunc: func(obj any, beanName string) (v any, err error) {
							vv := obj.(*int)
							*vv++
							return obj, nil
						},
					},
				)
				tc.p.AddBeanPostProcessor(
					&FuncBeanPostProcessor{
						AfterInitializationFunc: func(obj any, beanName string) (v any, err error) {
							return nil, xerrors.Errorf(tc.desp)
						},
					},
				)
			},
			obj:    pointer.IntPtr(0),
			err:    "failed at 2",
			expect: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			tc.init(&tc)

			actual, err := tc.p.PostProcessAfterInitialization(tc.obj, tc.desp)

			if tc.err != "" {
				g.Expect(err).To(HaveOccurred())
				g.Expect(err.Error()).To(MatchRegexp(tc.err))
				return
			}

			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(actual).To(Equal(tc.expect))
		})
	}
}

func TestPostProcessBeanDefinition(t *testing.T) {
	type testCase struct {
		desp string
		p    *beanPostProcessorCompositorImpl
		init func(tc *testCase)

		name   string
		expect string
	}
	testCases := []testCase{
		{
			desp: "all pass",
			p:    NewBeanPostProcessorCompositor().(*beanPostProcessorCompositorImpl),
			init: func(tc *testCase) {
				tc.p.AddBeanPostProcessor(
					&FuncBeanDefinitionPostProcessor{
						PostBeanDefinitionFunc: func(s string, bd BeanDefinition) {
							tc.name = s + "1"
						},
					},
				)
			},
			expect: "all pass1",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			tc.init(&tc)

			tc.p.PostProcessBeanDefinition(tc.desp, nil)
			g.Expect(tc.name).To(Equal(tc.expect))
		})
	}
}
