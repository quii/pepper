package assert_test

import (
	"fmt"
	. "github.com/quii/pepper/assert"
	"strings"
	"testing"
)

func HaveLength(length int) Matcher[string] {
	return func(in string) MatchResult {
		return MatchResult{
			Description: fmt.Sprintf("have length of %d, got %d", length, len(in)),
			Matches:     len(in) == length,
		}
	}
}

func HaveAllCaps(in string) MatchResult {
	return MatchResult{
		Description: "be in all caps",
		Matches:     strings.ToUpper(in) == in,
	}
}

func EqualTo[T comparable](in T) Matcher[T] {
	return func(got T) MatchResult {
		return MatchResult{
			Description: fmt.Sprintf("be equal to %v", in),
			Matches:     got == in,
		}
	}
}

func Contain(substring string) Matcher[string] {
	return func(in string) MatchResult {
		return MatchResult{
			Description: fmt.Sprintf("contain %q", substring),
			Matches:     strings.Contains(in, substring),
		}
	}
}

func TestStringMatchers(t *testing.T) {

	t.Run("failing", func(t *testing.T) {
		ExpectThat(t, "DDT").To(
			HaveLength(5),
			EqualTo("hello"),
			Contain("ell"),
			Doesnt(HaveAllCaps),
		)
	})
}
