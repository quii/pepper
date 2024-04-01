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

// Act is a helper function that will run the "act" function and expect no error. If there is an error it will call Errorf on the testing.TB. Otherwise, an inspector is returned. This is to cater for the fairly common case of writing tests around functions and methods that return T or an error.
func Act[T any](t TB, act func() (T, error)) Inspector[T] {
	t.Helper()
	result, err := act()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	return Inspector[T]{t, result}
}

// To is the method that actually runs the matchers. It will call Errorf on the testing.TB if any of the matchers fail.
func (e Inspector[T]) To(matchers ...Matcher[T]) {
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

func (e Inspector[T]) AndAssertSubject(matchers ...Matcher[T]) Inspector[T] {
	e.To(matchers...)
	return e
}
