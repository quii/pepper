package string

import (
	"fmt"
	"github.com/quii/pepper/matching"
	"strings"
)

// todo: this should be a higher order matcher so you can do HaveLength(LessThan(2)) - which then fits neatly into re-using comparable matchers
func HaveLength(length int) matching.Matcher[string] {
	return func(in string) matching.MatchResult {
		return matching.MatchResult{
			Description: fmt.Sprintf("have length of %d", length),
			Matches:     len(in) == length,
			But:         fmt.Sprintf("it has a length of %d", len(in)),
		}
	}
}

func HaveAllCaps(in string) matching.MatchResult {
	return matching.MatchResult{
		Description: "be in all caps",
		Matches:     strings.ToUpper(in) == in,
		But:         "it was not in all caps",
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
