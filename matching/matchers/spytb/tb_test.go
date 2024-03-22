package spytb

import (
	. "github.com/quii/pepper/matching"
	"testing"
)

func TestSpyTB(t *testing.T) {
	t.Run("correctly had errors", func(t *testing.T) {
		spyTB := &SpyTB{}
		subject := &SpyTB{ErrorCalls: []string{"oh no"}}

		Expect(spyTB, subject).To(HaveError("oopsie"))
		Expect(t, spyTB).To(HaveError(`expected Spy TB to have error "oopsie", but has [oh no]`))
	})

	t.Run("complains if it has errors when none expected", func(t *testing.T) {
		spyTB := &SpyTB{}
		subject := &SpyTB{ErrorCalls: []string{"oh no"}}

		Expect(spyTB, subject).To(HaveNoErrors)
		Expect(t, spyTB).To(HaveError(`expected Spy TB to have no errors, but it had errors [oh no]`))
	})
}
