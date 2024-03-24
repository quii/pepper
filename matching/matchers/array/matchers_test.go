package array

import (
	. "github.com/quii/pepper/matching"
	. "github.com/quii/pepper/matching/matchers/spytb"
	"testing"
)

func TestArrayMatchers(t *testing.T) {
	t.Run("Contain comparable", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			Expect(t, []int{1, 2, 3}).To(Contain(2))
		})

		t.Run("failing", func(t *testing.T) {
			spyTB := &SpyTB{}
			Expect(spyTB, []int{1, 2, 3}).To(Contain(4))
			Expect(t, spyTB).To(HaveError("expected [1 2 3] to contain 4, but it did not"))
		})
	})
}
