package io

import (
	"bytes"
	"fmt"
	"github.com/quii/pepper"
	"io"
)

// ContainingByte will check if the given byte slice is contained in the byte slice.
func ContainingByte(want []byte) pepper.Matcher[[]byte] {
	return func(have []byte) pepper.MatchResult {
		return pepper.MatchResult{
			Description: fmt.Sprintf("contain %q", want),
			Matches:     bytes.Contains(have, want),
			SubjectName: "the reader",
			But:         fmt.Sprintf("it didn't have %q", want),
		}
	}
}

// ContainingString will check if the given string is contained in the byte slice.
func ContainingString(want string) pepper.Matcher[[]byte] {
	return func(have []byte) pepper.MatchResult {
		return pepper.MatchResult{
			Description: fmt.Sprintf("contain %q", want),
			Matches:     bytes.Contains(have, []byte(want)),
			SubjectName: "the reader",
			But:         fmt.Sprintf("it didn't have %q", want),
		}
	}
}

// HaveData will read all the data from the io.Reader and run the given matcher on it.
func HaveData(matcher pepper.Matcher[[]byte]) pepper.Matcher[io.Reader] {
	return func(reader io.Reader) pepper.MatchResult {
		all, err := io.ReadAll(reader)
		if err != nil {
			return pepper.MatchResult{
				Description: "have data in io.Reader",
				Matches:     false,
				But:         "it could not be read",
			}
		}
		return matcher(all)
	}
}
