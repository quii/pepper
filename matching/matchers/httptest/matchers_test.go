package httptest

import (
	"github.com/quii/pepper/assert"
	. "github.com/quii/pepper/matching"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPTestMatchers(t *testing.T) {
	t.Run("Status code matchers", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			res := httptest.NewRecorder()
			spyTB := &SpyTB{}

			res.WriteHeader(http.StatusOK)
			ExpectThat(spyTB, StringedRes(*res)).To(BeOK)

			res.WriteHeader(http.StatusTeapot)
			ExpectThat(spyTB, StringedRes(*res)).To(Not(BeOK))

			t.Run("failure message", func(t *testing.T) {
				res := httptest.NewRecorder()
				spyTB := &SpyTB{}
				res.WriteHeader(http.StatusNotFound)
				ExpectThat(spyTB, StringedRes(*res)).To(BeOK)
				assert.True(t, len(spyTB.ErrorCalls) == 1)
				assert.Equal(t, spyTB.ErrorCalls[0], "expected the response to have status of 200, but it was 404")
			})
		})

		t.Run("user defined status", func(t *testing.T) {
			res := httptest.NewRecorder()
			spyTB := &SpyTB{}

			res.WriteHeader(http.StatusTeapot)
			ExpectThat(spyTB, StringedRes(*res)).To(HaveStatus(http.StatusTeapot))
		})
	})

}
