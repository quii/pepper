package spytb

import (
	"fmt"
	"github.com/quii/pepper/matching"
)

func HaveError(s string) matching.Matcher[*matching.SpyTB] {
	return func(tb *matching.SpyTB) matching.MatchResult {
		found := false
		for _, e := range tb.ErrorCalls {
			if e == s {
				found = true
			}
		}
		return matching.MatchResult{
			Description: fmt.Sprintf("have error %q", s),
			Matches:     found,
			But:         fmt.Sprintf("has %v", tb.ErrorCalls),
		}
	}
}

func HaveNoErrors(tb *matching.SpyTB) matching.MatchResult {
	return matching.MatchResult{
		Description: "have no errors",
		Matches:     len(tb.ErrorCalls) == 0,
		But:         fmt.Sprintf("it had errors %v", tb.ErrorCalls),
	}
}
