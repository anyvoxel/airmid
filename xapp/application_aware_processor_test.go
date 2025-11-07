package xapp

import (
	"testing"

	. "github.com/onsi/gomega"
)

type testApplicationEventPublisherAware struct {
	publisher ApplicationEventPublisher
}

func (t *testApplicationEventPublisherAware) SetApplicationEventPublisher(publisher ApplicationEventPublisher) {
	t.publisher = publisher
}

type testApplicationAware struct {
	app Application
}

func (t *testApplicationAware) SetApplication(app Application) {
	t.app = app
}

func TestPostProcessBeforeInitialization(t *testing.T) {
	app := &airmidApplication{}

	type testCase struct {
		desp   string
		obj    any
		name   string
		expect any
	}
	testCases := []testCase{
		{
			desp:   "not aware",
			obj:    "1",
			name:   "n1",
			expect: "1",
		},
		{
			desp: "application event publisher aware",
			obj: &testApplicationEventPublisherAware{
				publisher: nil,
			},
			name: "n1",
			expect: &testApplicationEventPublisherAware{
				publisher: app,
			},
		},
		{
			desp: "application aware",
			obj: &testApplicationAware{
				app: nil,
			},
			name: "n1",
			expect: &testApplicationAware{
				app: app,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			g := NewWithT(t)
			p := &ApplicationAwareProcessor{
				app: app,
			}
			_, err := p.PostProcessBeforeInitialization(tc.obj, tc.name)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(tc.obj).To(Equal(tc.expect))
		})
	}
}
