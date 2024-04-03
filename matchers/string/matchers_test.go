package string

import (
	"fmt"
	. "github.com/quii/pepper"
	. "github.com/quii/pepper/matchers/comparable"
	. "github.com/quii/pepper/matchers/spytb"
	"testing"
)

func ExampleHaveAllCaps() {
	t := &SpyTB{}

	Expect(t, "HELLO").To(HaveAllCaps)

	fmt.Println(t.LastError())
	//Output:
}

func ExampleHaveAllCaps_fail() {
	t := &SpyTB{}

	Expect(t, "hello").To(HaveAllCaps)

	fmt.Println(t.LastError())
	//Output: expected hello to in all caps, but it was not in all caps
}

func ExampleHaveLength() {
	t := &SpyTB{}

	Expect(t, "hello").To(HaveLength(EqualTo(5)))

	fmt.Println(t.LastError())
	//Output:
}

func ExampleHaveLength_fail() {
	t := &SpyTB{}

	Expect(t, "hello").To(HaveLength(EqualTo(4)))

	fmt.Println(t.LastError())
	//Output: expected hello to have length be equal to 4, but it was 5
}

func ExampleHaveSubstring() {
	t := &SpyTB{}

	Expect(t, "hello").To(HaveSubstring("ell"))

	fmt.Println(t.LastError())
	//Output:
}

func ExampleHaveSubstring_fail() {
	t := &SpyTB{}

	Expect(t, "hello").To(HaveSubstring("goodbye"))

	fmt.Println(t.LastError())
	//Output: expected hello to contain "goodbye"
}

func Example() {
	t := &SpyTB{}

	Expect(t, "hello").To(
		HaveLength(EqualTo(5)),
		EqualTo("hello"),
		HaveSubstring("ell"),
		Doesnt(HaveAllCaps),
	)

	fmt.Println(t.LastError())
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
