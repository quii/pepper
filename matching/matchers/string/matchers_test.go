package string

import (
	"github.com/quii/pepper/assert"
	. "github.com/quii/pepper/matching"
	"testing"
)

func TestStringMatchers(t *testing.T) {
	spyTB := &SpyTB{}

	t.Run("Have length", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			ExpectThat(spyTB, "hello").To(HaveLength(5))
			assert.True(t, len(spyTB.ErrorCalls) == 0)
		})

		t.Run("failing", func(t *testing.T) {
			ExpectThat(spyTB, "goodbye").To(HaveLength(5))
			assert.True(t, len(spyTB.ErrorCalls) == 1)
			assert.Equal(t, spyTB.ErrorCalls[0], "expected goodbye to have length of 5, got 7")
		})
	})
}
