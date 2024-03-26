package maps

import (
	"fmt"
	"github.com/quii/pepper/matching"
)

func HaveKey[K comparable, V any](key K, valueMatchers []matching.Matcher[V]) matching.Matcher[map[K]V] {
	return func(m map[K]V) matching.MatchResult {
		if value, ok := m[key]; ok {
			for _, valueMatcher := range valueMatchers {
				result := valueMatcher(value)
				if !result.Matches {
					return matching.MatchResult{
						Description: fmt.Sprintf("have key %v with value %v", key, result.Description),
						Matches:     false,
						But:         result.But,
					}
				}
			}

			return matching.MatchResult{
				Matches: true,
			}
		}

		return matching.MatchResult{
			Description: fmt.Sprintf("have key %v", key),
			Matches:     false,
			But:         "it did not",
		}
	}
}

func WithValue[T any](valueMatcher ...matching.Matcher[T]) []matching.Matcher[T] {
	for i := range valueMatcher {
		valueMatcher[i] = matching.Compose(func(v T) T { return v }, "value", valueMatcher[i])
	}
	return valueMatcher
}

func WithAnyValue[T any]() []matching.Matcher[T] {
	return nil
}
