package http

import (
	"fmt"
	"github.com/quii/pepper/matching"
	"io"
	"net/http"
)

func HaveStatus(status int) matching.Matcher[*http.Response] {
	return func(res *http.Response) matching.MatchResult {
		return matching.MatchResult{
			Description: fmt.Sprintf("have status of %d", status),
			But:         fmt.Sprintf("it was %d", res.StatusCode),
			Matches:     res.StatusCode == status,
			SubjectName: "the response",
		}
	}
}

func BeOK(res *http.Response) matching.MatchResult {
	return HaveStatus(http.StatusOK)(res)
}

func HaveHeader(header, value string) matching.Matcher[*http.Response] {
	return func(res *http.Response) matching.MatchResult {
		return matching.MatchResult{
			Description: fmt.Sprintf("have header %q of %q", header, value),
			Matches:     res.Header.Get(header) == value,
			But:         fmt.Sprintf("it was %q", res.Header.Get(header)),
			SubjectName: "the response",
		}
	}
}

func HaveJSONHeader(res *http.Response) matching.MatchResult {
	return HaveHeader("content-type", "application/json")(res)
}

func HaveBody(bodyMatchers ...matching.Matcher[string]) matching.Matcher[*http.Response] {
	bodyExtractor := func(res *http.Response) string {
		body, _ := io.ReadAll(res.Body)
		return string(body)
	}

	return matching.Compose(bodyExtractor, "the response body", bodyMatchers...)
}
