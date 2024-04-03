package pepper

type Matcher[T any] func(T) MatchResult

// Doesnt is a helper function to negate a matcher.
func Doesnt[T any](matcher Matcher[T]) Matcher[T] {
	return negate(matcher)
}

// Not is a helper function to negate a matcher.
func Not[T any](matcher Matcher[T]) Matcher[T] {
	return negate(matcher)
}

// Or combines matchers with a boolean OR.
func (m Matcher[T]) Or(matchers ...Matcher[T]) Matcher[T] {
	return func(got T) MatchResult {
		result := m(got)

		for _, matcher := range matchers {
			result.Description += " or " + matcher(got).Description
		}

		if result.Matches {
			return result
		}

		for _, matcher := range matchers {
			if r := matcher(got); r.Matches {
				result.Matches = true
				return result
			}
		}

		return result
	}
}

// And combines matchers with a boolean AND.
func (m Matcher[T]) And(matchers ...Matcher[T]) Matcher[T] {
	return func(got T) MatchResult {
		result := m(got)

		for _, matcher := range matchers {
			r := matcher(got)
			result = result.Combine(r)
		}

		return result
	}
}

func negate[T any](matcher Matcher[T]) Matcher[T] {
	return func(got T) MatchResult {
		result := matcher(got)
		return MatchResult{
			Description: "not " + result.Description,
			Matches:     !result.Matches,
		}
	}
}
