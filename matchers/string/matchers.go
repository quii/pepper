package string

import (
	"fmt"
	"github.com/quii/pepper"
	"strings"
)

// HaveLength will check a string's length meets the given matcher's criteria.
func HaveLength(matcher pepper.Matcher[int]) pepper.Matcher[string] {
	return func(in string) pepper.MatchResult {
		result := matcher(len(in))
		result.Description = fmt.Sprintf("have length %v", result.Description)
		return result
	}
}

// HaveAllCaps will check if a string is in all caps.
func HaveAllCaps(in string) pepper.MatchResult {
	return pepper.MatchResult{
		Description: "in all caps",
		Matches:     strings.ToUpper(in) == in,
		But:         "it was not in all caps",
	}
}

// Containing will check if a string contains a given substring.
func Containing(substring string) pepper.Matcher[string] {
	return func(in string) pepper.MatchResult {
		return pepper.MatchResult{
			Description: fmt.Sprintf("contain %q", substring),
			Matches:     strings.Contains(in, substring),
		}
	}
}
