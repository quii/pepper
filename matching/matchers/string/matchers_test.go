package string

import (
	. "github.com/quii/pepper/matching"
	. "github.com/quii/pepper/matching/matchers/spytb"
	"testing"
)

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
