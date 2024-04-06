package json

import (
	"encoding/json"
	"fmt"
	"github.com/quii/pepper"
	"io"
)

func Parse[T any](matcher pepper.Matcher[T]) pepper.Matcher[io.Reader] {
	return func(rdr io.Reader) pepper.MatchResult {
		var thing T
		err := json.NewDecoder(rdr).Decode(&thing)
		if err != nil {
			return pepper.MatchResult{
				Description: fmt.Sprintf("be parseable into %T", thing),
				SubjectName: "JSON",
				Matches:     false,
				But:         fmt.Sprintf("it could not be parsed: %v", err),
			}
		}
		return matcher(thing)
	}
}
