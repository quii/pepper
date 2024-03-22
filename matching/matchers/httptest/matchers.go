package httptest

import (
	"fmt"
	"github.com/quii/pepper/matching"
	"net/http"
	"net/http/httptest"
)

type StringedRes httptest.ResponseRecorder

func (s StringedRes) String() string {
	return "the response"
}

func HaveStatus(status int) matching.Matcher[StringedRes] {
	return func(res StringedRes) matching.MatchResult {
		return matching.MatchResult{
			Description: fmt.Sprintf("have status of %d", status),
			But:         fmt.Sprintf("it was %d", res.Code),
			Matches:     res.Code == status,
		}
	}
}

func BeOK(res StringedRes) matching.MatchResult {
	return HaveStatus(http.StatusOK)(res)
}

func HaveBody(bodyMatchers matching.Matcher[string]) matching.Matcher[StringedRes] {
	return func(res StringedRes) matching.MatchResult {
		return bodyMatchers(res.Body.String())
	}
}
