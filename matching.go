package pepper

import "fmt"

type (
	// TB is a cut-down version of testing.TB.
	TB interface {
		Error(args ...any)
		Errorf(format string, args ...any)
		Helper()
	}
	MatchResult struct {
		Description string
		Matches     bool
		But         string
		SubjectName string
	}

	Expecter[T any] struct {
		t       TB
		Subject T
	}

	// Matcher is a function that takes a subject T and returns a MatchResult.
	Matcher[T any] func(T) MatchResult
)

// Expect is the entry point for the matcher DSL. Pass in the testing.TB and the subject you want to test.
func Expect[T any](t TB, subject T) Expecter[T] {
	return Expecter[T]{t, subject}
}

// To is the method that actually runs the matchers. It will call Errorf on the testing.TB if any of the matchers fail.
func (e Expecter[T]) To(matchers ...Matcher[T]) {
	e.t.Helper()
	for _, matcher := range matchers {
		result := matcher(e.Subject)
		if result.SubjectName == "" {
			if str, isStringer := any(e.Subject).(fmt.Stringer); isStringer {
				result.SubjectName = str.String()
			} else {
				result.SubjectName = fmt.Sprintf("%v", e.Subject)
			}
		}
		if !result.Matches {
			if result.But != "" {
				e.t.Errorf("expected %+v to %+v, but %s", result.SubjectName, result.Description, result.But)
			} else {
				e.t.Errorf("expected %+v to %+v", result.SubjectName, result.Description)
			}
		}
	}
}

// Doesnt is a helper function to negate a matcher.
func Doesnt[T any](matcher Matcher[T]) Matcher[T] {
	return negate(matcher)
}

// Not is a helper function to negate a matcher.
func Not[T any](matcher Matcher[T]) Matcher[T] {
	return negate(matcher)
}

// Or combines matchers with a boolean OR.
func (m Matcher[T]) Or(matchers ...Matcher[T]) Matcher[T] {
	return func(got T) MatchResult {
		result := m(got)

		for _, matcher := range matchers {
			result.Description += " or " + matcher(got).Description
		}

		if result.Matches {
			return result
		}

		for _, matcher := range matchers {
			r := matcher(got)
			if r.Matches {
				result.Matches = true
				return result
			}
		}

		return result
	}
}

// And combines matchers with a boolean AND.
func (m Matcher[T]) And(matchers ...Matcher[T]) Matcher[T] {
	return func(got T) MatchResult {
		result := m(got)

		for _, matcher := range matchers {
			r := matcher(got)
			result = result.Combine(r)
		}

		return result
	}
}

func negate[T any](matcher Matcher[T]) Matcher[T] {
	return func(got T) MatchResult {
		result := matcher(got)
		return MatchResult{
			Description: "not " + result.Description,
			Matches:     !result.Matches,
		}
	}
}

// Zero returns true if the MatchResult is the zero value.
func (m MatchResult) Zero() bool {
	return m.Description == "" && m.But == "" && !m.Matches
}

// Combine merges two MatchResults into one.
func (m MatchResult) Combine(other MatchResult) MatchResult {
	if m.Zero() {
		return other
	}

	but := m.But + " and " + other.But

	if m.Matches && other.Matches {
		but = ""
	}

	if m.Matches && !other.Matches {
		but = other.But
	}

	if !m.Matches && other.Matches {
		but = m.But
	}

	return MatchResult{
		Description: m.Description + " and " + other.Description,
		Matches:     m.Matches && other.Matches,
		But:         but,
		SubjectName: m.SubjectName,
	}
}
