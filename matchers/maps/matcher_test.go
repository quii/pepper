package maps

import (
	"fmt"
	. "github.com/quii/pepper"
	. "github.com/quii/pepper/matchers/comparable"
	"github.com/quii/pepper/matchers/spytb"
	"testing"
)

func ExampleWithAnyValue() {
	t := &SpyTB{}

	Expect(t, map[string]string{"hello": "world"}).To(HaveKey("goodbye", WithAnyValue[string]()))

	fmt.Println(t.LastError())
	//Output: expected map[hello:world] to have key goodbye, but it did not
}

func ExampleHaveKey_fail() {
	t := &SpyTB{}

	Expect(t, map[string]int{"score": 4}).To(HaveKey("score", GreaterThan(5).And(LessThan(10))))

	fmt.Println(t.LastError())
	//Output: expected map[score:4] to have key score with value be greater than 5 and be less than 10, but it was 4
}

func ExampleHaveKey() {
	t := &SpyTB{}

	Expect(t, map[string]string{"hello": "world"}).To(HaveKey("hello", EqualTo("world")))

	fmt.Println(t.LastError())
	//Output:
}

func TestMapMatching(t *testing.T) {
	t.Run("HasKey WithValue", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			Expect(t, map[string]string{"hello": "world"}).To(HaveKey("hello", EqualTo("world")))
			Expect(t, map[string]int{"score": 7}).To(HaveKey("score", GreaterThan(5).And(LessThan(10))))
		})
	})

	t.Run("failures", func(t *testing.T) {
		t.Run("missing key", func(t *testing.T) {
			spytb.VerifyFailingMatcher(
				t,
				map[string]string{"hello": "world"},
				HaveKey("goodbye", WithAnyValue[string]()),
				`expected map[hello:world] to have key goodbye, but it did not`,
			)
		})

		t.Run("key exists but value does not match", func(t *testing.T) {
			spytb.VerifyFailingMatcher(
				t,
				map[string]string{"hello": "world"},
				HaveKey("hello", EqualTo("goodbye")),
				`expected map[hello:world] to have key hello with value be equal to goodbye, but it was world`,
			)
		})
	})
}
