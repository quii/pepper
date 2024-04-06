package http

import (
	"fmt"
	"github.com/quii/pepper"
	"io"
	"net/http"
)

const (
	subjectName             = "the response"
	responseBodySubjectName = subjectName + " body"
)

// HaveStatus returns a matcher that checks if the response status code is equal to the given status code.
func HaveStatus(status int) pepper.Matcher[*http.Response] {
	return func(res *http.Response) pepper.MatchResult {
		return pepper.MatchResult{
			Description: fmt.Sprintf("have status of %d", status),
			But:         fmt.Sprintf("it was %d", res.StatusCode),
			Matches:     res.StatusCode == status,
			SubjectName: subjectName,
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
			SubjectName: subjectName,
		}
	}
}

// HaveJSONHeader is a convenience matcher for HaveHeader("content-type", "application/json").
func HaveJSONHeader(res *http.Response) pepper.MatchResult {
	return HaveHeader("content-type", "application/json")(res)
}

// HaveBody returns a matcher that checks if the response body meets the given matchers' criteria. Note this will read the entire body using io.ReadAll.
func HaveBody(bodyMatchers pepper.Matcher[io.Reader]) pepper.Matcher[*http.Response] {
	return func(res *http.Response) pepper.MatchResult {
		result := bodyMatchers(res.Body)
		result.SubjectName = responseBodySubjectName
		return result
	}
}
