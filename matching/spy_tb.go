package matching

import "fmt"

type SpyTB struct {
	ErrorCalls []string
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
