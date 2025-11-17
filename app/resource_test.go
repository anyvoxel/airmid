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
	"io"
	"testing"

	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	propsmocks "github.com/anyvoxel/airmid/ioc/props/mocks"
)

func TestConfigLoadProperty(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := NewWithT(t)

		mockCtrl := gomock.NewController(t)
		mockResourceLocator := NewMockResourceLocator(mockCtrl)
		mockProps := propsmocks.NewMockProperties(mockCtrl)

		c := &config{
			resourceLocator:  mockResourceLocator,
			ConfigExtensions: []string{".yaml", ".yml"},
			ActiveProfiles:   []string{"test"},
		}
		data := []byte(`k1: v1
k2: v2`)
		mockResource := NewMockResource(mockCtrl)
		mockResource.EXPECT().Read(gomock.Any()).DoAndReturn(
			func(p []byte) (n int, err error) {
				copy(p, data)
				return len(data), io.EOF
			},
		)
		mockResource.EXPECT().Name().Return("application.yaml").Times(2)

		mockResourceLocator.EXPECT().Locate(gomock.Eq("application.yaml")).Times(1).Return(
			[]Resource{mockResource},
			nil,
		)
		mockResourceLocator.EXPECT().Locate(gomock.Eq("application.yml")).Times(1).Return(
			nil,
			nil,
		)
		mockResourceLocator.EXPECT().Locate(gomock.Eq("application-test.yaml")).Times(1).Return(
			nil,
			nil,
		)
		mockResourceLocator.EXPECT().Locate(gomock.Eq("application-test.yml")).Times(1).Return(
			nil,
			nil,
		)
		mockProps.EXPECT().Set(context.Background(), gomock.Eq("k1"), gomock.Eq("v1")).Times(1).Return(nil)
		mockProps.EXPECT().Set(context.Background(), gomock.Eq("k2"), gomock.Eq("v2")).Times(1).Return(nil)

		err := c.loadProperty(context.Background(), mockProps)
		g.Expect(err).ToNot(HaveOccurred())
	})
}
