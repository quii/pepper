package spytb

import "github.com/quii/pepper/matching"

func VerifyFailingMatcher[T any](t matching.TB, subject T, matcher matching.Matcher[T], expectedError string) {
	t.Helper()
	spyTB := &matching.SpyTB{}
	matching.Expect(spyTB, subject).To(matcher)
	matching.Expect(t, spyTB).To(HaveError(expectedError))
}
