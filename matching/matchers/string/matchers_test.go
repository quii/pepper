package string

import (
	. "github.com/quii/pepper/matching"
	. "github.com/quii/pepper/matching/matchers/spy_tb"
	"testing"
)

func TestStringMatchers(t *testing.T) {

	t.Run("Have length", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			spyTB := &SpyTB{}
			ExpectThat(spyTB, "hello").To(HaveLength(5))
			ExpectThat(t, spyTB).To(HaveNoErrors)
		})

		t.Run("failing", func(t *testing.T) {
			spyTB := &SpyTB{}
			ExpectThat(spyTB, "goodbye").To(HaveLength(5))
			ExpectThat(t, spyTB).To(HaveError("expected goodbye to have length of 5, got 7"))
		})
	})
}
