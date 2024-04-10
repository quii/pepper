package pepper

import "fmt"

// SpyTB is a test helper that records calls to Error. This lets us create tests and examples to demonstrate what happens when a test fails.
type SpyTB struct {
	ErrorCalls []string
}

// Result will print the test result.
func (s *SpyTB) Result() string {
	if len(s.ErrorCalls) == 0 {
		return "Test passed"
	}
	return fmt.Sprintf("Test failed: %v", s.ErrorCalls)
}

func (s *SpyTB) String() string {
	return "Spy TB"
}

func (s *SpyTB) Helper() {
}

func (s *SpyTB) Error(args ...any) {
	s.ErrorCalls = append(s.ErrorCalls, fmt.Sprint(args...))
}

func (s *SpyTB) Errorf(format string, args ...any) {
	s.ErrorCalls = append(s.ErrorCalls, fmt.Sprintf(format, args...))
}

func (s *SpyTB) Fatalf(format string, args ...any) {
	s.ErrorCalls = append(s.ErrorCalls, fmt.Sprintf(format, args...))
}

func (s *SpyTB) Reset() {
	s.ErrorCalls = nil
}
