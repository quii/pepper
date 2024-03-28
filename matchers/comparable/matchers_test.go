package comparable

import (
	. "github.com/quii/pepper"
	"github.com/quii/pepper/matchers/spytb"
	"testing"
)

func TestComparisonMatchers(t *testing.T) {
	t.Run("Less than", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			Expect(t, 5).To(LessThan(6))
		})

		t.Run("failing", func(t *testing.T) {
			spytb.VerifyFailingMatcher(t, 6, LessThan(6), "expected 6 to be less than 6, but it was 6")
			spytb.VerifyFailingMatcher(t, 6, LessThan(3), "expected 6 to be less than 3, but it was 6")
		})
	})

	t.Run("Greater than", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			Expect(t, 5).To(GreaterThan(4))
		})

		t.Run("failing", func(t *testing.T) {
			spytb.VerifyFailingMatcher(t, 6, GreaterThan(6), "expected 6 to be greater than 6, but it was 6")
			spytb.VerifyFailingMatcher(t, 2, GreaterThan(10), "expected 2 to be greater than 10, but it was 2")
		})
	})
}
