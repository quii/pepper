package slice

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

	anArray := []string{"HELLO", "WORLD"}
	Expect(t, anArray).To(ContainItem(HaveAllCaps))

	fmt.Println(t.LastError())
	//Output:
}

func ExampleContainItem_fail() {
	t := &SpyTB{}

	anArray := []string{"hello", "world"}
	Expect(t, anArray).To(ContainItem(HaveAllCaps))

	fmt.Println(t.LastError())
	//Output: expected [hello world] to contain an item in all caps, but it did not
}

func ExampleHaveSize() {
	t := &SpyTB{}

	anArray := []string{"hello", "world"}
	Expect(t, anArray).To(HaveSize[string](EqualTo(2)))

	fmt.Println(t.LastError())
	//Output:
}

func ExampleHaveSize_fail() {
	t := &SpyTB{}

	anArray := []string{"hello", "world"}
	Expect(t, anArray).To(HaveSize[string](EqualTo(3)))

	fmt.Println(t.LastError())
	//Output: expected [hello world] to have a size be equal to 3, but it was 2
}

func ExampleEveryItem() {
	t := &SpyTB{}

	anArray := []string{"hello", "world"}
	Expect(t, anArray).To(EveryItem(HaveSubstring("o")))

	fmt.Println(t.LastError())
	//Output:
}

func ExampleEveryItem_fail() {
	t := &SpyTB{}

	anArray := []string{"hello", "world"}
	Expect(t, anArray).To(EveryItem(HaveSubstring("h")))

	fmt.Println(t.LastError())
	//Output: expected [hello world] to have every item contain "h"
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
					`expected [hello world] to contain an item be equal to goodbye, but it did not`,
				)
			})
			t.Run("all caps", func(t *testing.T) {
				VerifyFailingMatcher(
					t,
					[]string{"hello", "world"},
					ContainItem(HaveAllCaps),
					`expected [hello world] to contain an item in all caps, but it did not`,
				)
			})
		})
	})
}
