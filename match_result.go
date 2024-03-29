package pepper

// MatchResult describes the result of a match against a "subject".
type MatchResult struct {
	Description string
	Matches     bool
	But         string
	SubjectName string
}

// Zero returns true if the MatchResult is the zero value.
func (m MatchResult) Zero() bool {
	return m.Description == "" && m.But == "" && !m.Matches
}

// Combine merges two MatchResults into one.
func (m MatchResult) Combine(other MatchResult) MatchResult {
	if m.Zero() {
		return other
	}

	but := m.But + " and " + other.But

	if m.Matches && other.Matches {
		but = ""
	}

	if m.Matches && !other.Matches {
		but = other.But
	}

	if !m.Matches && other.Matches {
		but = m.But
	}

	return MatchResult{
		Description: m.Description + " and " + other.Description,
		Matches:     m.Matches && other.Matches,
		But:         but,
		SubjectName: m.SubjectName,
	}
}
