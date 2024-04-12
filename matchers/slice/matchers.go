package slice

import (
	"fmt"
	"github.com/quii/pepper"
	"slices"
)

var passingResult = pepper.MatchResult{
	Matches: true,
}

// HaveSize checks if an array's size meets a matcher's criteria.
func HaveSize[T any](matcher pepper.Matcher[int]) pepper.Matcher[[]T] {
	return func(items []T) pepper.MatchResult {
		result := matcher(len(items))
		result.Description = "have a size " + result.Description
		return result
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
					Description: "contain an item",
					Matches:     true,
				}
			} else {
				exampleFailure = result
			}
		}

		exampleFailure.But = "it did not"
		exampleFailure.Description = "contain an item " + exampleFailure.Description
		exampleFailure.SubjectName = fmt.Sprintf("%+v", items)

		return exampleFailure
	}
}

// EveryItem checks if every item in an array meets a matcher's criteria.
func EveryItem[T any](m pepper.Matcher[T]) pepper.Matcher[[]T] {
	return func(items []T) pepper.MatchResult {
		for _, item := range items {
			if result := m(item); !result.Matches {
				return everyItemFailure(result)
			}
		}

		return passingResult
	}
}

// ShallowEquals checks if two slices are equal, only works with slices of comparable types.
func ShallowEquals[T comparable](other []T) pepper.Matcher[[]T] {
	return func(ts []T) pepper.MatchResult {
		equal := slices.Equal(ts, other)
		but := ""
		if !equal {
			but = "the slice is not equal"
		}
		return pepper.MatchResult{
			Matches:     equal,
			Description: fmt.Sprintf("be equal to %v", other),
			But:         but,
		}
	}
}

func everyItemFailure(result pepper.MatchResult) pepper.MatchResult {
	return pepper.MatchResult{
		Description: "have every item " + result.Description,
		Matches:     false,
		But:         result.But,
	}
}
