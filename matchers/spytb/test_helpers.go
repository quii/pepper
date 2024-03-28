package spytb

import (
	"github.com/quii/pepper"
)

func VerifyFailingMatcher[T any](t pepper.TB, subject T, matcher pepper.Matcher[T], expectedError string) {
	t.Helper()
	spyTB := &pepper.SpyTB{}
	pepper.Expect(spyTB, subject).To(matcher)
	pepper.Expect(t, spyTB).To(HaveError(expectedError))
}
