package http

import (
	"fmt"
	"github.com/quii/pepper"
	"io"
	"net/http"
)

// HaveStatus returns a matcher that checks if the response status code is equal to the given status code.
func HaveStatus(status int) pepper.Matcher[*http.Response] {
	return func(res *http.Response) pepper.MatchResult {
		return pepper.MatchResult{
			Description: fmt.Sprintf("have status of %d", status),
			But:         fmt.Sprintf("it was %d", res.StatusCode),
			Matches:     res.StatusCode == status,
			SubjectName: "the response",
		}
	}
}

// BeOK is a convenience matcher for HaveStatus(http.StatusOK).
func BeOK(res *http.Response) pepper.MatchResult {
	return HaveStatus(http.StatusOK)(res)
}

// HaveHeader returns a matcher that checks if the response has a header with the given name and value.
func HaveHeader(header, value string) pepper.Matcher[*http.Response] {
	return func(res *http.Response) pepper.MatchResult {
		return pepper.MatchResult{
			Description: fmt.Sprintf("have header %q of %q", header, value),
			Matches:     res.Header.Get(header) == value,
			But:         fmt.Sprintf("it was %q", res.Header.Get(header)),
			SubjectName: "the response",
		}
	}
}

// HaveJSONHeader is a convenience matcher for HaveHeader("content-type", "application/json").
func HaveJSONHeader(res *http.Response) pepper.MatchResult {
	return HaveHeader("content-type", "application/json")(res)
}

// HaveBody returns a matcher that checks if the response body meets the given matchers' criteria.
func HaveBody(bodyMatchers pepper.Matcher[string]) pepper.Matcher[*http.Response] {
	return func(res *http.Response) pepper.MatchResult {
		body, _ := io.ReadAll(res.Body)
		result := bodyMatchers(string(body))
		result.SubjectName = "the response body"
		return result
	}
}
