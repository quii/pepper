package pepper_test

import (
	"errors"
	"fmt"
	. "github.com/quii/pepper"
	. "github.com/quii/pepper/matchers/comparable"
	. "github.com/quii/pepper/matchers/string"
	"testing"
)

func ExampleInspector_To() {
	t := &SpyTB{}
	Expect(t, "Pepper").To(EqualTo("Stanley"))
	fmt.Println(t.LastError())
	//Output: expected Pepper to be equal to Stanley, but it was Pepper
}

func ExampleMatcher_Or() {
	t := &SpyTB{}
	tshirt := TShirt{Colour: "yellow"}

	Expect(t, tshirt).To(HaveColour("blue").Or(HaveColour("red")))
	fmt.Println(t.LastError())
	//Output: expected the t-shirt to have colour "blue" or have colour "red", but it was "yellow"
}

func ExampleNot() {
	t := &SpyTB{}

	tshirt := TShirt{Colour: "yellow"}

	Expect(t, tshirt).To(Not(HaveColour("yellow")))
	fmt.Println(t.LastError())
	//Output: expected the t-shirt to not have colour "yellow"
}

func ExampleMatcher_And() {
	t := &SpyTB{}
	player := Player{
		Name:   "Chris",
		Points: 11,
	}

	Expect(t, player).To(HaveScore(
		GreaterThan(5).And(LessThan(10)),
	))
	fmt.Println(t.LastError())
	//Output: expected Player Chris to score be greater than 5 and be less than 10, but it was 11
}

func ExampleExpectNoError() {
	t := &SpyTB{}

	err := errors.New("oh no")

	ExpectNoError(t, err)
	fmt.Println(t.LastError())
	//Output: unexpected error: oh no
}

func ExampleExpectError() {
	t := &SpyTB{}

	err := errors.New("oh no")

	ExpectError(t, err)
	fmt.Println(t.LastError())
	//Output:
}

func ExampleExpectErrorOfType() {
	t := &SpyTB{}

	unauthorised := errors.New("unauthorised")
	wrappedErr := fmt.Errorf("oh no: %w", unauthorised)

	ExpectErrorOfType(t, wrappedErr, unauthorised)
	fmt.Println(t.LastError())
	//Output:
}

func ExampleExpectErrorOfType_failing() {
	t := &SpyTB{}

	unauthorised := errors.New("unauthorised")
	wrappedErr := fmt.Errorf("oh no: %w", unauthorised)

	ExpectErrorOfType(t, wrappedErr, errors.New("not found"))
	fmt.Println(t.LastError())
	//Output: expected error of type *errors.errorString, but got "oh no: unauthorised"
}

func ExampleAct() {
	t := &SpyTB{}

	// Often we want to test functions that return a value and an error, Act is a helper to make this easier
	subject := func() (string, error) {
		return "hello", nil
	}

	Act(t, subject).AndAssertSubject(EqualTo("hello"))
	fmt.Println(t.LastError())
	//Output:
}

func TestMatching(t *testing.T) {
	t.Run("passing example", func(t *testing.T) {
		Expect(t, "hello").To(
			HaveLength(EqualTo(5)),
			EqualTo("hello"),
			HaveSubstring("ell"),
			Doesnt(HaveAllCaps),
		)
	})

	t.Run("combining failures", func(t *testing.T) {
		t.Run("when it has a but and both failed", func(t *testing.T) {
			someString := "goodbye"
			result1 := HaveLength(EqualTo(5))(someString)
			result2 := HaveAllCaps(someString)

			expected := MatchResult{
				Description: `have length be equal to 5 and in all caps`,
				Matches:     false,
				But:         "it was 7 and it was not in all caps",
			}

			actual := result1.Combine(result2)
			Expect(t, actual).To(EqualTo(expected))
		})

		t.Run("when nothing fails", func(t *testing.T) {
			someString := "HELLO"
			result1 := HaveLength(EqualTo(5))(someString)
			result2 := HaveAllCaps(someString)

			expected := MatchResult{
				Description: `have length be equal to 5 and in all caps`,
				Matches:     true,
			}

			actual := result1.Combine(result2)
			Expect(t, actual).To(EqualTo(expected))
		})

		t.Run("when first match is passing but second is failing", func(t *testing.T) {
			someString := "hello"
			result1 := HaveLength(EqualTo(5))(someString)
			result2 := HaveAllCaps(someString)

			expected := MatchResult{
				Description: `have length be equal to 5 and in all caps`,
				Matches:     false,
				But:         "it was not in all caps",
			}

			actual := result1.Combine(result2)
			Expect(t, actual).To(EqualTo(expected))
		})

		t.Run("when first match is failing but second is passing", func(t *testing.T) {
			someString := "GOODBYE"
			result1 := HaveLength(EqualTo(5))(someString)
			result2 := HaveAllCaps(someString)

			expected := MatchResult{
				Description: `have length be equal to 5 and in all caps`,
				Matches:     false,
				But:         "it was 7",
			}

			actual := result1.Combine(result2)
			Expect(t, actual).To(EqualTo(expected))
		})
	})
}

type TShirt struct {
	Colour string
}

func (t TShirt) String() string {
	return "the t-shirt"
}

func HaveColour(colour string) Matcher[TShirt] {
	return func(t TShirt) MatchResult {
		return MatchResult{
			Description: fmt.Sprintf("have colour %q", colour),
			Matches:     t.Colour == colour,
			But:         fmt.Sprintf("it was %q", t.Colour),
		}
	}
}

type Player struct {
	Name   string
	Points int
}

func (s Player) String() string {
	return fmt.Sprintf("Player %s", s.Name)
}

func HaveScore(matcher Matcher[int]) Matcher[Player] {
	return func(s Player) MatchResult {
		result := matcher(s.Points)
		return MatchResult{
			Description: "score " + result.Description,
			Matches:     result.Matches,
			But:         fmt.Sprintf("it was %d", s.Points),
		}
	}
}
