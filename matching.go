package pepper

import "fmt"

type (
	// TB is a cut-down version of testing.TB.
	TB interface {
		Error(args ...any)
		Errorf(format string, args ...any)
		Helper()
	}
	Expecter[T any] struct {
		t       TB
		Subject T
	}
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
