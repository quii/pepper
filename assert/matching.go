package assert

import "testing"

type (
	MatchResult struct {
		Description string
		Matches     bool
	}

	Expect[T any] struct {
		t       *testing.T
		Subject T
	}

	Matcher[T any] func(T) MatchResult
)

func ExpectThat[T any](t *testing.T, subject T) Expect[T] {
	return Expect[T]{t, subject}
}

func (e Expect[T]) To(matchers ...Matcher[T]) {
	e.t.Helper()
	for _, matcher := range matchers {
		result := matcher(e.Subject)
		if !result.Matches {
			e.t.Errorf("expected %v to %v", e.Subject, result.Description)
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

func That[T any](t *testing.T, got T, matchers ...func(T) MatchResult) {
	t.Helper()
	for _, matcher := range matchers {
		result := matcher(got)
		if !result.Matches {
			t.Fatalf("expected %v to %v", got, result.Description)
		}
	}
}
