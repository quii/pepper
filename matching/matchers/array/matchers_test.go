package array

import (
	. "github.com/quii/pepper/matching"
	. "github.com/quii/pepper/matching/matchers/comparable"
	. "github.com/quii/pepper/matching/matchers/spytb"
	. "github.com/quii/pepper/matching/matchers/string"
	"testing"
)

func TestArrayMatchers(t *testing.T) {
	t.Run("contain with other matcher to find matcher", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			Expect(t, []string{"hello", "WORLD"}).To(ContainItem(HaveAllCaps))
		})

		t.Run("failing", func(t *testing.T) {
			t.Run("equal to", func(t *testing.T) {
				spyTB := &SpyTB{}
				Expect(spyTB, []string{"hello", "world"}).To(ContainItem(EqualTo("goodbye")))
				Expect(t, spyTB).To(HaveError(`expected [hello world] to contain item be equal to goodbye, but it did not`))
			})
			t.Run("all caps", func(t *testing.T) {
				spyTB := &SpyTB{}
				Expect(spyTB, []string{"hello", "world"}).To(ContainItem(HaveAllCaps))
				Expect(t, spyTB).To(HaveError(`expected [hello world] to contain item be in all caps, but it did not`))
			})
		})
	})
}
