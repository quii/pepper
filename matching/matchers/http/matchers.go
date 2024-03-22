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

func HaveJSONHeader(res *http.Response) matching.MatchResult {
	return matching.MatchResult{
		Description: "have JSON header",
		Matches:     res.Header.Get("Content-Type") == "application/json",
		But:         fmt.Sprintf("it was %q", res.Header.Get("Content-Type")),
		SubjectName: "the response",
	}
}

func HaveBody(bodyMatchers ...matching.Matcher[string]) matching.Matcher[*http.Response] {
	return func(res *http.Response) matching.MatchResult {
		body, _ := io.ReadAll(res.Body)
		for _, matcher := range bodyMatchers {
			result := matcher(string(body))
			result.SubjectName = "the response body"
			if !result.Matches {
				return result
			}
		}

		return matching.MatchResult{
			Description: "have body matching",
			Matches:     true,
			SubjectName: "the response body",
		}
	}
}
