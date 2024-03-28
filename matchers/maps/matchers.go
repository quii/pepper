package maps

import (
	"fmt"
	"github.com/quii/pepper"
)

// HaveKey checks if a map has a key with a specific value.
func HaveKey[K comparable, V any](key K, valueMatcher pepper.Matcher[V]) pepper.Matcher[map[K]V] {
	return func(m map[K]V) pepper.MatchResult {
		if value, ok := m[key]; ok {
			result := valueMatcher(value)
			if !result.Matches {
				return pepper.MatchResult{
					Description: fmt.Sprintf("have key %v with value %v", key, result.Description),
					Matches:     false,
					But:         result.But,
				}
			}

			return pepper.MatchResult{
				Matches: true,
			}
		}

		return pepper.MatchResult{
			Description: fmt.Sprintf("have key %v", key),
			Matches:     false,
			But:         "it did not",
		}
	}
}

// WithAnyValue lets you match any value, useful if you're just looking for the presence of a key.
func WithAnyValue[T any]() pepper.Matcher[T] {
	return func(T) pepper.MatchResult {
		return pepper.MatchResult{
			Matches: true,
		}
	}
}
