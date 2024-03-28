package pepper

import "fmt"

type SpyTB struct {
	ErrorCalls []string
}

func (s *SpyTB) LastError() string {
	if len(s.ErrorCalls) == 0 {
		return ""
	}

	return s.ErrorCalls[len(s.ErrorCalls)-1]
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

func (s *SpyTB) Reset() {
	s.ErrorCalls = nil
}
