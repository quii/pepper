package comparable

import (
	"cmp"
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

func LessThan[T cmp.Ordered](in T) matching.Matcher[T] {
	return func(got T) matching.MatchResult {
		return matching.MatchResult{
			Description: fmt.Sprintf("be less than %v", in),
			Matches:     got < in,
		}
	}
}

func GreaterThan[T cmp.Ordered](in T) matching.Matcher[T] {
	return func(got T) matching.MatchResult {
		return matching.MatchResult{
			Description: fmt.Sprintf("be greater than %v", in),
			Matches:     got > in,
		}
	}
}
