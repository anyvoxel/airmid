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
	"context"
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
			_, err := p.PostProcessBeforeInitialization(context.Background(), tc.obj, tc.name)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(tc.obj).To(Equal(tc.expect))
		})
	}
}
