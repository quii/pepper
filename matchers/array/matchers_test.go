package array

import (
	"fmt"
	. "github.com/quii/pepper"
	. "github.com/quii/pepper/matchers/comparable"
	. "github.com/quii/pepper/matchers/spytb"
	. "github.com/quii/pepper/matchers/string"
	"testing"
)

func ExampleContainItem() {
	t := &SpyTB{}

	anArray := []string{"hello", "world"}
	Expect(t, anArray).To(ContainItem(HaveAllCaps))

	fmt.Println(t.LastError())
	//Output: expected [hello world] to contain item be in all caps, but it did not
}

func TestArrayMatchers(t *testing.T) {
	t.Run("contain with other matcher to find matcher", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			Expect(t, []string{"hello", "WORLD"}).To(ContainItem(HaveAllCaps))
		})

		t.Run("failing", func(t *testing.T) {
			t.Run("equal to", func(t *testing.T) {
				VerifyFailingMatcher(
					t,
					[]string{"hello", "world"},
					ContainItem(EqualTo("goodbye")),
					`expected [hello world] to contain item be equal to goodbye, but it did not`,
				)
			})
			t.Run("all caps", func(t *testing.T) {
				VerifyFailingMatcher(
					t,
					[]string{"hello", "world"},
					ContainItem(HaveAllCaps),
					`expected [hello world] to contain item be in all caps, but it did not`,
				)
			})
		})
	})
}
