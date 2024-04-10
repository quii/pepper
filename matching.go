package pepper

import (
	"errors"
	"fmt"
)

type (
	// TB is a cut-down version of testing.TB.
	TB interface {
		Error(args ...any)
		Errorf(format string, args ...any)
		Fatalf(format string, args ...any)
		Helper()
	}
	Inspector[T any] struct {
		t       TB
		Subject T
	}
)

// Expect is the entry point for the matcher DSL. Pass in the testing.TB and the subject you want to test.
func Expect[T any](t TB, subject T) Inspector[T] {
	return Inspector[T]{t, subject}
}

// ExpectNoError is a helper function that will call t.Fatalf if the error is not nil.
func ExpectNoError(t TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ExpectError is a helper function that will call t.Fatalf if the error is nil.
func ExpectError(t TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected an error")
	}
}

// ExpectErrorOfType is a helper function that will call t.Fatalf if the error is not of the expected type.
func ExpectErrorOfType(t TB, err error, expectedType error) {
	t.Helper()
	if !errors.Is(err, expectedType) {
		t.Fatalf("expected error of type %T, but got %q", expectedType, err.Error())
	}
}

// To is the method that actually runs the matchers. It will call Errorf on the testing.TB if any of the matchers fail.
func (e Inspector[T]) To(matchers ...Matcher[T]) {
	e.t.Helper()
	for _, matcher := range matchers {
		result := matcher(e.Subject)

		if result.SubjectName == "" {
			result.SubjectName = calculateSubjectName(e)
		}

		if !result.Matches {
			e.t.Error(result.Error())
		}
	}
}

func calculateSubjectName[T any](e Inspector[T]) string {
	var subjectName = fmt.Sprintf("%v", e.Subject)

	if str, isStringer := any(e.Subject).(fmt.Stringer); isStringer {
		subjectName = str.String()
	}
	return subjectName
}
