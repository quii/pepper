package string

import (
	"fmt"
	. "github.com/quii/pepper/matching"
	. "github.com/quii/pepper/matching/matchers/comparable"
	. "github.com/quii/pepper/matching/matchers/spytb"
	"testing"
)

func Example() {
	t := &SpyTB{}

	Expect(t, "hello").To(
		HaveLength(5),
		EqualTo("hello"),
		HaveSubstring("ell"),
		Doesnt(HaveAllCaps),
	)

	fmt.Println(len(t.ErrorCalls)) // no error calls means it passed
	//Output: 0
}
func TestStringMatchers(t *testing.T) {
	t.Run("Have length", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			Expect(t, "hello").To(HaveLength(5))
		})

		t.Run("failing", func(t *testing.T) {
			spyTB := &SpyTB{}
			Expect(spyTB, "goodbye").To(HaveLength(5))
			Expect(t, spyTB).To(HaveError("expected goodbye to have length of 5, but it has a length of 7"))
		})
	})
}
