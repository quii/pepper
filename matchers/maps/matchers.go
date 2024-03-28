package maps

import (
	"fmt"
	"github.com/quii/pepper"
)

func HaveKey[K comparable, V any](key K, valueMatchers []pepper.Matcher[V]) pepper.Matcher[map[K]V] {
	return func(m map[K]V) pepper.MatchResult {
		if value, ok := m[key]; ok {
			for _, valueMatcher := range valueMatchers {
				result := valueMatcher(value)
				if !result.Matches {
					return pepper.MatchResult{
						Description: fmt.Sprintf("have key %v with value %v", key, result.Description),
						Matches:     false,
						But:         result.But,
					}
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

func WithValue[T any](valueMatcher ...pepper.Matcher[T]) []pepper.Matcher[T] {
	for i := range valueMatcher {
		valueMatcher[i] = pepper.Compose(func(v T) T { return v }, "value", valueMatcher[i])
	}
	return valueMatcher
}

func WithAnyValue[T any]() []pepper.Matcher[T] {
	return nil
}
