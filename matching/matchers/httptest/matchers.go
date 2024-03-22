package httptest

import (
	"fmt"
	"github.com/quii/pepper/matching"
	"net/http"
	"net/http/httptest"
)

type MatchableRes httptest.ResponseRecorder

func (s MatchableRes) String() string {
	return "the response"
}

func HaveStatus(status int) matching.Matcher[MatchableRes] {
	return func(res MatchableRes) matching.MatchResult {
		return matching.MatchResult{
			Description: fmt.Sprintf("have status of %d", status),
			But:         fmt.Sprintf("it was %d", res.Code),
			Matches:     res.Code == status,
		}
	}
}

func BeOK(res MatchableRes) matching.MatchResult {
	return HaveStatus(http.StatusOK)(res)
}

func HaveBody(bodyMatchers ...matching.Matcher[string]) matching.Matcher[MatchableRes] {
	return func(res MatchableRes) matching.MatchResult {
		body := res.Body.String()
		for _, matcher := range bodyMatchers {
			result := matcher(body)
			if !result.Matches {
				return result
			}
		}

		return matching.MatchResult{
			Description: "have body matching",
			Matches:     true,
		}
	}
}
