package spytb

import (
	"fmt"
	"github.com/quii/pepper"
)

func HaveError(s string) pepper.Matcher[*pepper.SpyTB] {
	return func(tb *pepper.SpyTB) pepper.MatchResult {
		found := false
		for _, e := range tb.ErrorCalls {
			if e == s {
				found = true
			}
		}
		return pepper.MatchResult{
			Description: fmt.Sprintf("have error %q", s),
			Matches:     found,
			But:         fmt.Sprintf("has %v", tb.ErrorCalls),
		}
	}
}

func HaveNoErrors(tb *pepper.SpyTB) pepper.MatchResult {
	return pepper.MatchResult{
		Description: "have no errors",
		Matches:     len(tb.ErrorCalls) == 0,
		But:         fmt.Sprintf("it had errors %v", tb.ErrorCalls),
	}
}
