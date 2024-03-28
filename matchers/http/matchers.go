package http

import (
	"fmt"
	"github.com/quii/pepper"
	"io"
	"net/http"
)

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

func BeOK(res *http.Response) pepper.MatchResult {
	return HaveStatus(http.StatusOK)(res)
}

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

func HaveJSONHeader(res *http.Response) pepper.MatchResult {
	return HaveHeader("content-type", "application/json")(res)
}

func HaveBody(bodyMatchers ...pepper.Matcher[string]) pepper.Matcher[*http.Response] {
	bodyExtractor := func(res *http.Response) string {
		body, _ := io.ReadAll(res.Body)
		return string(body)
	}

	return pepper.Compose(bodyExtractor, "the response body", bodyMatchers...)
}
