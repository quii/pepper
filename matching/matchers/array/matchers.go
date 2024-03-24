package array

import "github.com/quii/pepper/matching"

func ContainItem(m matching.Matcher[string]) matching.Matcher[[]string] {
	return func(items []string) matching.MatchResult {
		var exampleFailure matching.MatchResult
		for _, item := range items {
			result := m(item)
			if result.Matches {
				return matching.MatchResult{
					Description: "contain item",
					Matches:     true,
				}
			} else {
				exampleFailure = result
			}
		}

		exampleFailure.But = "it did not"
		exampleFailure.Description = "contain item " + exampleFailure.Description

		return exampleFailure
	}
}
