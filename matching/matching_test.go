package matching_test

import (
	"fmt"
	. "github.com/quii/pepper/matching"
	. "github.com/quii/pepper/matching/matchers/comparable"
	. "github.com/quii/pepper/matching/matchers/string"
	"testing"
)

func ExampleExpecter_To() {
	t := &SpyTB{}
	name := "Pepper"

	Expect(t, name).To(EqualTo("Stanley"))
	fmt.Println(t.ErrorCalls[0])
	//Output: expected Pepper to be equal to Stanley, but it was Pepper
}

func ExampleMatcher_Or() {
	t := &SpyTB{}
	tshirt := TShirt{Colour: "yellow"}

	Expect(t, tshirt).To(HaveColour("blue").Or(HaveColour("red")))
	fmt.Println(t.ErrorCalls[0])
	//Output: expected the t-shirt to have colour "blue" or have colour "red", but it was "yellow"
}

func ExampleNot() {
	t := &SpyTB{}

	tshirt := TShirt{Colour: "yellow"}

	Expect(t, tshirt).To(Not(HaveColour("yellow")))
	fmt.Println(t.ErrorCalls[0])
	//Output: expected the t-shirt to not have colour "yellow"
}

func ExampleMatcher_And() {
	t := &SpyTB{}
	score := Score{
		Name:   "Chris",
		Points: 11,
	}

	Expect(t, score).To(HaveScore(
		GreaterThan(5).And(LessThan(10))),
	)
	fmt.Println(t.ErrorCalls[0])
	//Output: expected Chris's score to be greater than 5 and be less than 10, but it was 11
}

func TestMatching(t *testing.T) {
	t.Run("passing example", func(t *testing.T) {
		Expect(t, "hello").To(
			HaveLength(5),
			EqualTo("hello"),
			HaveSubstring("ell"),
			Doesnt(HaveAllCaps),
		)
	})

	t.Run("combining failures", func(t *testing.T) {
		t.Run("when it has a but and both failed", func(t *testing.T) {
			someString := "goodbye"
			result1 := HaveLength(5)(someString)
			result2 := HaveAllCaps(someString)

			expected := MatchResult{
				Description: `have length of 5 and be in all caps`,
				Matches:     false,
				But:         "it has a length of 7 and it was not in all caps",
			}

			actual := result1.Combine(result2)
			Expect(t, actual).To(EqualTo(expected))
		})

		t.Run("when nothing fails", func(t *testing.T) {
			someString := "HELLO"
			result1 := HaveLength(5)(someString)
			result2 := HaveAllCaps(someString)

			expected := MatchResult{
				Description: `have length of 5 and be in all caps`,
				Matches:     true,
			}

			actual := result1.Combine(result2)
			Expect(t, actual).To(EqualTo(expected))
		})

		t.Run("when first match is passing but second is failing", func(t *testing.T) {
			someString := "hello"
			result1 := HaveLength(5)(someString)
			result2 := HaveAllCaps(someString)

			expected := MatchResult{
				Description: `have length of 5 and be in all caps`,
				Matches:     false,
				But:         "it was not in all caps",
			}

			actual := result1.Combine(result2)
			Expect(t, actual).To(EqualTo(expected))
		})

		t.Run("when first match is failing but second is passing", func(t *testing.T) {
			someString := "GOODBYE"
			result1 := HaveLength(5)(someString)
			result2 := HaveAllCaps(someString)

			expected := MatchResult{
				Description: `have length of 5 and be in all caps`,
				Matches:     false,
				But:         "it has a length of 7",
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

type Score struct {
	Name   string
	Points int
}

func (s Score) String() string {
	return fmt.Sprintf("%s's score", s.Name)
}

func HaveScore(matcher Matcher[int]) Matcher[Score] {
	return func(s Score) MatchResult {
		result := matcher(s.Points)
		return MatchResult{
			Description: result.Description,
			Matches:     result.Matches,
			But:         fmt.Sprintf("it was %d", s.Points),
		}
	}
}
