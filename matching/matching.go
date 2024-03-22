package matching

type (
	TB interface {
		Error(args ...any)
		Errorf(format string, args ...any)
		Helper()
	}
	MatchResult struct {
		Description string
		Matches     bool
	}

	Expect[T any] struct {
		t       TB
		Subject T
	}

	Matcher[T any] func(T) MatchResult
)

func ExpectThat[T any](t TB, subject T) Expect[T] {
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
