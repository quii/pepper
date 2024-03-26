package array

import "github.com/quii/pepper/matching"

func ContainItem[T any](m matching.Matcher[T]) matching.Matcher[[]T] {
	return func(items []T) matching.MatchResult {
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

func EveryItem[T any](m matching.Matcher[T]) matching.Matcher[[]T] {
	panic("not yet implemented")
}
