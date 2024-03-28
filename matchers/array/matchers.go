package array

import (
	"github.com/quii/pepper"
)

// HaveSize checks if an array's size meets a matcher's criteria.
func HaveSize[T any](matcher pepper.Matcher[int]) pepper.Matcher[[]T] {
	return func(items []T) pepper.MatchResult {
		return matcher(len(items))
	}
}

// ContainItem checks if an array contains an item that meets a matcher's criteria.
func ContainItem[T any](m pepper.Matcher[T]) pepper.Matcher[[]T] {
	return func(items []T) pepper.MatchResult {
		var exampleFailure pepper.MatchResult

		for _, item := range items {
			result := m(item)
			if result.Matches {
				return pepper.MatchResult{
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

// EveryItem checks if every item in an array meets a matcher's criteria.
func EveryItem[T any](m pepper.Matcher[T]) pepper.Matcher[[]T] {
	return func(items []T) pepper.MatchResult {
		for _, item := range items {
			result := m(item)
			if !result.Matches {
				return pepper.MatchResult{
					Description: "have every item " + result.Description,
					Matches:     false,
					But:         result.But,
				}
			}
		}

		return pepper.MatchResult{
			Matches: true,
		}
	}
}
