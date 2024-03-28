package pepper

import "fmt"

type (
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

	Matcher[T any] func(T) MatchResult
)

func Expect[T any](t TB, subject T) Expecter[T] {
	return Expecter[T]{t, subject}
}

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

func Doesnt[T any](matcher Matcher[T]) Matcher[T] {
	return negate(matcher)
}

func Not[T any](matcher Matcher[T]) Matcher[T] {
	return negate(matcher)
}

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

func (m Matcher[T]) And(matchers ...Matcher[T]) Matcher[T] {
	return func(got T) MatchResult {
		result := m(got)

		for _, matcher := range matchers {
			result.Description += " and " + matcher(got).Description
		}

		if !result.Matches {
			return result
		}

		for _, matcher := range matchers {
			r := matcher(got)
			if !r.Matches {
				result.Matches = false
				result.But = r.But
				return result
			}
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

func (m MatchResult) Zero() bool {
	return m.Description == "" && m.But == "" && !m.Matches
}

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

func Compose[A, B any](getSubject func(A) B, subjectName string, matchers ...Matcher[B]) Matcher[A] {
	return func(parentSubject A) MatchResult {
		childSubject := getSubject(parentSubject)

		var combinedMatchResults = MatchResult{
			SubjectName: subjectName,
		}

		hasMatched := true

		for _, matcher := range matchers {
			result := matcher(childSubject)
			result.SubjectName = subjectName
			if !result.Matches {
				hasMatched = false
				combinedMatchResults = combinedMatchResults.Combine(result)
			}
		}

		combinedMatchResults.Matches = hasMatched

		return combinedMatchResults
	}
}
