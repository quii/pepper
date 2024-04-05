package comparable

import (
	"cmp"
	"fmt"
	"github.com/quii/pepper"
)

// EqualTo checks if a value is equal to another value.
func EqualTo[T comparable](in T) pepper.Matcher[T] {
	return func(got T) pepper.MatchResult {
		return pepper.MatchResult{
			Description: fmt.Sprintf("be equal to %+v", in),
			Matches:     got == in,
			But:         fmt.Sprintf("it was %v", got),
		}
	}
}

// Equal is an alias for EqualTo.
func Equal[T comparable](in T) pepper.Matcher[T] {
	return EqualTo(in)
}

// LessThan checks if a value is less than another value.
func LessThan[T cmp.Ordered](in T) pepper.Matcher[T] {
	return func(got T) pepper.MatchResult {
		return pepper.MatchResult{
			Description: fmt.Sprintf("be less than %v", in),
			Matches:     got < in,
			But:         fmt.Sprintf("it was %v", got),
		}
	}
}

// GreaterThan checks if a value is greater than another value.
func GreaterThan[T cmp.Ordered](in T) pepper.Matcher[T] {
	return func(got T) pepper.MatchResult {
		return pepper.MatchResult{
			Description: fmt.Sprintf("be greater than %v", in),
			Matches:     got > in,
			But:         fmt.Sprintf("it was %v", got),
		}
	}
}
