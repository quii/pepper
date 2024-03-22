package httptest

import (
	. "github.com/quii/pepper/matching"
	. "github.com/quii/pepper/matching/matchers/spy_tb"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPTestMatchers(t *testing.T) {
	t.Run("Status code matchers", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			t.Run("positive happy path", func(t *testing.T) {
				res := httptest.NewRecorder()
				spyTB := &SpyTB{}
				res.WriteHeader(http.StatusOK)
				ExpectThat(spyTB, StringedRes(*res)).To(BeOK)
				ExpectThat(t, spyTB).To(HaveNoErrors)
			})

			t.Run("negation on happy path", func(t *testing.T) {
				res := httptest.NewRecorder()
				spyTB := &SpyTB{}
				res.WriteHeader(http.StatusTeapot)
				ExpectThat(spyTB, StringedRes(*res)).To(Not(BeOK))
				ExpectThat(t, spyTB).To(HaveNoErrors)
			})

			t.Run("failure message", func(t *testing.T) {
				res := httptest.NewRecorder()
				spyTB := &SpyTB{}
				res.WriteHeader(http.StatusNotFound)
				ExpectThat(spyTB, StringedRes(*res)).To(BeOK)
				ExpectThat(t, spyTB).To(HaveError(`expected the response to have status of 200, but it was 404`))
			})
		})

		t.Run("user defined status", func(t *testing.T) {
			res := httptest.NewRecorder()
			spyTB := &SpyTB{}

			res.WriteHeader(http.StatusTeapot)
			ExpectThat(spyTB, StringedRes(*res)).To(HaveStatus(http.StatusTeapot))
			ExpectThat(t, spyTB).To(HaveNoErrors)
		})
	})

}
