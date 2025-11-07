package xerrors

import (
	"errors"
	"testing"

	"github.com/onsi/gomega"
)

func TestErrNotFound(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := gomega.NewWithT(t)
		err := WrapNotFound("normal")
		g.Expect(err).To(gomega.HaveOccurred())

		g.Expect(IsNotFound(err)).To(gomega.BeTrue())
	})

	t.Run("is_not_found_with_nil", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsNotFound(nil)).To(gomega.BeFalse())
	})

	t.Run("is_not_found_with_other", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsNotFound(errors.New("other"))).To(gomega.BeFalse())
	})
}

func TestErrDuplicate(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := gomega.NewWithT(t)
		err := WrapDuplicate("normal")
		g.Expect(err).To(gomega.HaveOccurred())

		g.Expect(IsDuplicate(err)).To(gomega.BeTrue())
	})

	t.Run("is_duplidate_with_nil", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsDuplicate(nil)).To(gomega.BeFalse())
	})

	t.Run("is_duplidate_with_other", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsDuplicate(errors.New("other"))).To(gomega.BeFalse())
	})
}

func TestErrNotImplement(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := gomega.NewWithT(t)
		err := WrapNotImplement("normal")
		g.Expect(err).To(gomega.HaveOccurred())

		g.Expect(IsNotImplement(err)).To(gomega.BeTrue())
	})

	t.Run("is_not_implement_with_nil", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsNotImplement(nil)).To(gomega.BeFalse())
	})

	t.Run("is_not_implement_with_other", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsNotImplement(errors.New("other"))).To(gomega.BeFalse())
	})
}

func TestErrContinue(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := gomega.NewWithT(t)
		err := WrapContinue("normal")
		g.Expect(err).To(gomega.HaveOccurred())

		g.Expect(IsContinue(err)).To(gomega.BeTrue())
	})

	t.Run("is_continue_with_nil", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsContinue(nil)).To(gomega.BeFalse())
	})

	t.Run("is_continue_with_other", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsContinue(errors.New("other"))).To(gomega.BeFalse())
	})
}

func TestErrRetryable(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := gomega.NewWithT(t)
		err := WrapRetryable("normal")
		g.Expect(err).To(gomega.HaveOccurred())

		g.Expect(IsRetryable(err)).To(gomega.BeTrue())
	})

	t.Run("is_retryable_with_nil", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsRetryable(nil)).To(gomega.BeFalse())
	})

	t.Run("is_retryable_with_other", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsRetryable(errors.New("other"))).To(gomega.BeFalse())
	})
}

func TestErrNonRetryable(t *testing.T) {
	t.Run("normal test", func(t *testing.T) {
		g := gomega.NewWithT(t)
		err := WrapNonRetryable("normal")
		g.Expect(err).To(gomega.HaveOccurred())

		g.Expect(IsNonRetryable(err)).To(gomega.BeTrue())
	})

	t.Run("is_non_retryable_with_nil", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsNonRetryable(nil)).To(gomega.BeFalse())
	})

	t.Run("is_non_retryable_with_other", func(t *testing.T) {
		g := gomega.NewWithT(t)

		g.Expect(IsNonRetryable(errors.New("other"))).To(gomega.BeFalse())
	})
}
