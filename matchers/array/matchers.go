package array

import (
	"github.com/quii/pepper"
)

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

func EveryItem[T any](m pepper.Matcher[T]) pepper.Matcher[[]T] {
	panic("not yet implemented")
}
