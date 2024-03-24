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

	t.Run("combining failures", func(t *testing.T) {
		t.Run("when it has a but and both failed", func(t *testing.T) {
			someString := "goodbye"
			result1 := HaveLength(5)(someString)
			result2 := HaveAllCaps(someString)

			expected := MatchResult{
				Description: `have length of 5 and be in all caps`,
				Matches:     false,
				But:         "it has a length of 7 and it was not in all caps",
			}

			actual := result1.Combine(result2)
			Expect(t, actual).To(EqualTo(expected))
		})

		t.Run("when nothing fails", func(t *testing.T) {
			someString := "HELLO"
			result1 := HaveLength(5)(someString)
			result2 := HaveAllCaps(someString)

			expected := MatchResult{
				Description: `have length of 5 and be in all caps`,
				Matches:     true,
			}

			actual := result1.Combine(result2)
			Expect(t, actual).To(EqualTo(expected))
		})

		t.Run("when first match is passing but second is failing", func(t *testing.T) {
			someString := "hello"
			result1 := HaveLength(5)(someString)
			result2 := HaveAllCaps(someString)

			expected := MatchResult{
				Description: `have length of 5 and be in all caps`,
				Matches:     false,
				But:         "it was not in all caps",
			}

			actual := result1.Combine(result2)
			Expect(t, actual).To(EqualTo(expected))
		})

		t.Run("when first match is failing but second is passing", func(t *testing.T) {
			someString := "GOODBYE"
			result1 := HaveLength(5)(someString)
			result2 := HaveAllCaps(someString)

			expected := MatchResult{
				Description: `have length of 5 and be in all caps`,
				Matches:     false,
				But:         "it has a length of 7",
			}

			actual := result1.Combine(result2)
			Expect(t, actual).To(EqualTo(expected))
		})
	})
}
