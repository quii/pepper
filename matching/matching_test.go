package matching_test

import (
	. "github.com/quii/pepper/matching"
	. "github.com/quii/pepper/matching/matchers/comparable"
	. "github.com/quii/pepper/matching/matchers/string"
	"testing"
)

func TestStringMatchers(t *testing.T) {
	t.Run("passing example", func(t *testing.T) {
		Expect(t, "hello").To(
			HaveLength(5),
			EqualTo("hello"),
			Contain("ell"),
			Doesnt(HaveAllCaps),
		)
	})
}
