package comparable

import (
	"fmt"
	. "github.com/quii/pepper"
	"github.com/quii/pepper/matchers/spytb"
	"testing"
)

func ExampleEqual() {
	t := &SpyTB{}
	Expect(t, 2).To(Equal(2))
	fmt.Println(t.LastError())
	// Output:
}

func ExampleEqual_fail() {
	t := &SpyTB{}
	Expect(t, 2).To(Equal(1))
	fmt.Println(t.LastError())
	// Output: expected 2 to be equal to 1, but it was 2
}

func ExampleEqualTo() {
	t := &SpyTB{}
	Expect(t, 5).To(EqualTo(5))
	fmt.Println(t.LastError())
	// Output:
}

func ExampleEqualTo_fail() {
	t := &SpyTB{}
	Expect(t, 5).To(EqualTo(4))
	fmt.Println(t.LastError())
	// Output: expected 5 to be equal to 4, but it was 5
}

func ExampleGreaterThan() {
	t := &SpyTB{}
	Expect(t, 5).To(GreaterThan(4))
	fmt.Println(t.LastError())
	// Output:
}

func ExampleGreaterThan_fail() {
	t := &SpyTB{}
	Expect(t, 5).To(GreaterThan(6))
	fmt.Println(t.LastError())
	// Output: expected 5 to be greater than 6, but it was 5
}

func ExampleLessThan() {
	t := &SpyTB{}
	Expect(t, 5).To(LessThan(6))
	fmt.Println(t.LastError())
	// Output:
}

func ExampleLessThan_fail() {
	t := &SpyTB{}
	Expect(t, 5).To(LessThan(4))
	fmt.Println(t.LastError())
	// Output: expected 5 to be less than 4, but it was 5
}

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

	t.Run("equal to with empty strings", func(t *testing.T) {
		t.Run("when it is an empty string, failing output should be quoted", func(t *testing.T) {
			spytb.VerifyFailingMatcher(
				t,
				"",
				EqualTo("Bob"),
				`expected "" to be equal to "Bob", but it was ""`,
			)
		})
	})
}
