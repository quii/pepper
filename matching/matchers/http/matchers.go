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
	return func(res *http.Response) matching.MatchResult {
		body, _ := io.ReadAll(res.Body)
		var combinedMatchResults *matching.MatchResult = nil

		for _, matcher := range bodyMatchers {
			result := matcher(string(body))
			result.SubjectName = "the response body"
			if !result.Matches {
				if combinedMatchResults == nil {
					combinedMatchResults = &result
				} else {
					//todo: there's gonna be a nicer way than this, but make it work, make it right...
					foo := combinedMatchResults.Combine(result)
					combinedMatchResults = &foo
				}
			}
		}

		if combinedMatchResults == nil {
			return matching.MatchResult{
				Description: "have body matching",
				Matches:     true,
				SubjectName: "the response body",
			}
		}

		return *combinedMatchResults
	}
}
