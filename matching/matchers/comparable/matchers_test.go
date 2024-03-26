package comparable

import (
	. "github.com/quii/pepper/matching"
	. "github.com/quii/pepper/matching/matchers/spytb"
	"testing"
)

// todo: testing matchers feels boilerplatey, can probably be further extracted
func TestComparisonMatchers(t *testing.T) {
	t.Run("Less than", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			Expect(t, 5).To(LessThan(6))
		})

		t.Run("failing", func(t *testing.T) {
			spyTB := &SpyTB{}
			Expect(spyTB, 6).To(LessThan(6))
			Expect(t, spyTB).To(HaveError("expected 6 to be less than 6"))
		})
	})

	t.Run("Greater than", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			Expect(t, 5).To(GreaterThan(4))
		})

		t.Run("failing", func(t *testing.T) {
			spyTB := &SpyTB{}
			Expect(spyTB, 6).To(GreaterThan(6))
			Expect(t, spyTB).To(HaveError("expected 6 to be greater than 6"))
		})
	})
}
