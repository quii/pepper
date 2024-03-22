package matching

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
			result.SubjectName = fmt.Sprintf("%v", e.Subject)
		}
		if !result.Matches {
			if result.But != "" {
				e.t.Errorf("expected %v to %v, but %s", result.SubjectName, result.Description, result.But)
			} else {
				e.t.Errorf("expected %v to %v", result.SubjectName, result.Description)
			}
		}
	}
}

func Doesnt[T any](matcher Matcher[T]) Matcher[T] {
	return func(got T) MatchResult {
		result := matcher(got)
		return MatchResult{
			Description: "not " + result.Description,
			Matches:     !result.Matches,
		}
	}
}

func Not[T any](matcher Matcher[T]) Matcher[T] {
	return func(got T) MatchResult {
		result := matcher(got)
		return MatchResult{
			Description: "not " + result.Description,
			Matches:     !result.Matches,
		}
	}
}
