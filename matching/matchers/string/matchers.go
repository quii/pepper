package string

import (
	"fmt"
	"github.com/quii/pepper/matching"
	"strings"
)

func HaveLength(length int) matching.Matcher[string] {
	return func(in string) matching.MatchResult {
		return matching.MatchResult{
			Description: fmt.Sprintf("have length of %d, got %d", length, len(in)),
			Matches:     len(in) == length,
		}
	}
}

func HaveAllCaps(in string) matching.MatchResult {
	return matching.MatchResult{
		Description: "be in all caps",
		Matches:     strings.ToUpper(in) == in,
	}
}

func Contain(substring string) matching.Matcher[string] {
	return func(in string) matching.MatchResult {
		return matching.MatchResult{
			Description: fmt.Sprintf("contain %q", substring),
			Matches:     strings.Contains(in, substring),
		}
	}
}
