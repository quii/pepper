package comparable

import (
	"fmt"
	"github.com/quii/pepper/matching"
)

func EqualTo[T comparable](in T) matching.Matcher[T] {
	return func(got T) matching.MatchResult {
		return matching.MatchResult{
			Description: fmt.Sprintf("be equal to %v", in),
			Matches:     got == in,
			But:         fmt.Sprintf("it was %v", got),
		}
	}
}
