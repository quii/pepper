package http_test

import (
	"encoding/json"
	"fmt"
	. "github.com/quii/pepper"
	. "github.com/quii/pepper/matchers/comparable"
	. "github.com/quii/pepper/matchers/http"
	. "github.com/quii/pepper/matchers/spytb"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPTestMatchers(t *testing.T) {
	t.Run("Body matching", func(t *testing.T) {
		t.Run("simple string match", func(t *testing.T) {
			res := httptest.NewRecorder()

			res.Body.WriteString("Hello, world")

			// see how we can compose matchers together!
			Expect(t, res.Result()).To(HaveBody(EqualTo("Hello, world")))
		})

		t.Run("simple string mismatch", func(t *testing.T) {
			res := httptest.NewRecorder()
			spyTB := &SpyTB{}

			res.Body.WriteString("Hello, world")

			Expect(spyTB, res.Result()).To(HaveBody(EqualTo("Goodbye, world")))
			Expect(t, spyTB).To(HaveError("expected the response body to be equal to Goodbye, world, but it was Hello, world"))
		})

		t.Run("example of matching JSON", func(t *testing.T) {
			type Todo struct {
				Name        string    `json:"name"`
				Completed   bool      `json:"completed"`
				LastUpdated time.Time `json:"last_updated"`
			}

			WithCompletedTODO := func(body string) MatchResult {
				var todo Todo
				_ = json.Unmarshal([]byte(body), &todo)
				return MatchResult{
					Description: "have a completed todo",
					Matches:     todo.Completed,
					But:         "it wasn't complete",
				}
			}
			WithTodoNameOf := func(todoName string) Matcher[string] {
				return func(body string) MatchResult {
					var todo Todo
					_ = json.Unmarshal([]byte(body), &todo)
					return MatchResult{
						Description: fmt.Sprintf("have a todo name of %q", todoName),
						Matches:     todo.Name == todoName,
						But:         fmt.Sprintf("it was %q", todo.Name),
					}
				}
			}

			t.Run("with completed todo", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.Body.WriteString(`{"name": "Finish the side project", "completed": true}`)
				Expect(t, res.Result()).To(HaveBody(WithCompletedTODO))
			})

			t.Run("with incomplete todo", func(t *testing.T) {
				res := httptest.NewRecorder()
				spyTB := &SpyTB{}

				res.Body.WriteString(`{"name": "Finish the side project", "completed": false}`)
				Expect(spyTB, res.Result()).To(HaveBody(WithCompletedTODO))
				Expect(t, spyTB).To(HaveError("expected the response body to have a completed todo, but it wasn't complete"))
			})

			t.Run("with a todo name", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.Body.WriteString(`{"name": "Finish the side project", "completed": false}`)
				Expect(t, res.Result()).To(HaveBody(WithTodoNameOf("Finish the side project")))
			})

			t.Run("with incorrect todo name and not completed", func(t *testing.T) {
				res := httptest.NewRecorder()
				spyTB := &SpyTB{}

				res.Body.WriteString(`{"name": "Egg", "completed": false}`)
				Expect(spyTB, res.Result()).To(HaveBody(WithTodoNameOf("Bacon"), WithCompletedTODO))
				Expect(t, spyTB).To(
					HaveError(`expected the response body to have a todo name of "Bacon" and have a completed todo, but it was "Egg" and it wasn't complete`),
				)
			})

			t.Run("with incorrect todo name", func(t *testing.T) {
				res := httptest.NewRecorder()
				spyTB := &SpyTB{}

				res.Body.WriteString(`{"name": "Egg", "completed": false}`)
				Expect(spyTB, res.Result()).To(HaveBody(WithTodoNameOf("Bacon")))
				Expect(t, spyTB).To(
					HaveError(`expected the response body to have a todo name of "Bacon", but it was "Egg"`),
				)
			})

			t.Run("compose the matchers", func(t *testing.T) {
				res := httptest.NewRecorder()

				res.Body.WriteString(`{"name": "Egg", "completed": false}`)
				res.Header().Add("content-type", "application/json")

				Expect(t, res.Result()).To(
					BeOK,
					HaveJSONHeader,
					HaveBody(WithTodoNameOf("Egg"), Not(WithCompletedTODO)),
				)
			})

			t.Run("compose the matchers to show an unexpected response", func(t *testing.T) {
				res := httptest.NewRecorder()
				spyTB := &SpyTB{}

				res.Body.WriteString(`{"name": "Egg", "completed": true}`)

				Expect(spyTB, res.Result()).To(
					BeOK,
					HaveJSONHeader,
					HaveBody(WithTodoNameOf("Egg"), Not(WithCompletedTODO)),
				)
				Expect(t, spyTB).To(HaveError(`expected the response body to not have a completed todo`))
				Expect(t, spyTB).To(HaveError(`expected the response to have header "content-type" of "application/json", but it was ""`))
			})
		})
	})

	t.Run("Status code matchers", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			t.Run("positive happy path", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.WriteHeader(http.StatusOK)
				Expect(t, res.Result()).To(BeOK)
			})

			t.Run("negation on happy path", func(t *testing.T) {
				res := httptest.NewRecorder()
				res.WriteHeader(http.StatusTeapot)
				Expect(t, res.Result()).To(Not(BeOK))
			})

			t.Run("failure message", func(t *testing.T) {
				res := httptest.NewRecorder()
				spyTB := &SpyTB{}
				res.WriteHeader(http.StatusNotFound)
				Expect(spyTB, res.Result()).To(BeOK)
				Expect(t, spyTB).To(HaveError(`expected the response to have status of 200, but it was 404`))
			})
		})

		t.Run("user defined status", func(t *testing.T) {
			res := httptest.NewRecorder()
			res.WriteHeader(http.StatusTeapot)
			Expect(t, res.Result()).To(HaveStatus(http.StatusTeapot))
		})
	})

	t.Run("Header matchers", func(t *testing.T) {
		t.Run("happy path multiple headers", func(t *testing.T) {
			res := httptest.NewRecorder()
			res.Header().Add("Content-Encoding", "gzip")
			res.Header().Add("Content-Type", "text/html")

			Expect(t, res.Result()).To(
				HaveHeader("Content-Encoding", "gzip"),
				HaveHeader("Content-Type", "text/html"),
			)
		})

		t.Run("unhappy path with multiple headers", func(t *testing.T) {
			res := httptest.NewRecorder()
			spyTB := &SpyTB{}
			res.Header().Add("Content-Type", "text/xml")

			Expect(spyTB, res.Result()).To(
				HaveHeader("Content-Encoding", "gzip"),
				HaveHeader("Content-Type", "text/html"),
			)
			Expect(t, spyTB).To(
				HaveError(`expected the response to have header "Content-Encoding" of "gzip", but it was ""`),
				HaveError(`expected the response to have header "Content-Type" of "text/html", but it was "text/xml"`),
			)
		})
	})
}
