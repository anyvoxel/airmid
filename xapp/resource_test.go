package xapp

import (
	"io"
	"testing"

	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	propsmocks "github.com/anyvoxel/airmid/props/mocks"
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
		mockProps.EXPECT().Set(gomock.Eq("k1"), gomock.Eq("v1")).Times(1).Return(nil)
		mockProps.EXPECT().Set(gomock.Eq("k2"), gomock.Eq("v2")).Times(1).Return(nil)

		err := c.loadProperty(mockProps)
		g.Expect(err).ToNot(HaveOccurred())
	})
}
