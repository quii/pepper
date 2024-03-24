package array

import (
	"fmt"
	"github.com/quii/pepper/matching"
)

func Contain[T comparable](needle T) matching.Matcher[[]T] {
	return func(haystack []T) matching.MatchResult {
		for _, h := range haystack {
			if h == needle {
				return matching.MatchResult{
					Description: fmt.Sprintf("contain %v", needle),
					Matches:     true,
				}
			}
		}

		return matching.MatchResult{
			Description: fmt.Sprintf("contain %v", needle),
			But:         "it did not",
		}
	}
}
