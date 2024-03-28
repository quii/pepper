package string

import (
	"fmt"
	. "github.com/quii/pepper"
	. "github.com/quii/pepper/matchers/comparable"
	. "github.com/quii/pepper/matchers/spytb"
	"testing"
)

func Example() {
	t := &SpyTB{}

	Expect(t, "hello").To(
		HaveLength(EqualTo(5)),
		EqualTo("hello"),
		HaveSubstring("ell"),
		Doesnt(HaveAllCaps),
	)

	fmt.Println(t.LastError()) // no error calls means it passed
	//Output:
}
func TestStringMatchers(t *testing.T) {
	t.Run("Have length", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			Expect(t, "hello").To(HaveLength(EqualTo(5)))
		})

		t.Run("failing", func(t *testing.T) {
			spyTB := &SpyTB{}
			Expect(spyTB, "goodbye").To(HaveLength(EqualTo(5)))
			Expect(t, spyTB).To(HaveError("expected goodbye to have length be equal to 5, but it was 7"))
		})
	})
}
