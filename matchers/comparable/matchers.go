package comparable

import (
	"cmp"
	"fmt"
	"github.com/quii/pepper"
)

func EqualTo[T comparable](in T) pepper.Matcher[T] {
	return func(got T) pepper.MatchResult {
		return pepper.MatchResult{
			Description: fmt.Sprintf("be equal to %+v", in),
			Matches:     got == in,
			But:         fmt.Sprintf("it was %v", got),
		}
	}
}

func LessThan[T cmp.Ordered](in T) pepper.Matcher[T] {
	return func(got T) pepper.MatchResult {
		return pepper.MatchResult{
			Description: fmt.Sprintf("be less than %v", in),
			Matches:     got < in,
			But:         fmt.Sprintf("it was %v", got),
		}
	}
}

func GreaterThan[T cmp.Ordered](in T) pepper.Matcher[T] {
	return func(got T) pepper.MatchResult {
		return pepper.MatchResult{
			Description: fmt.Sprintf("be greater than %v", in),
			Matches:     got > in,
			But:         fmt.Sprintf("it was %v", got),
		}
	}
}
