package json

import (
	"bytes"
	"fmt"
	. "github.com/quii/pepper"
	. "github.com/quii/pepper/matchers/comparable"
	"io"
)

func ExampleParse() {
	t := &SpyTB{}

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	HasName := func(want string) Matcher[Person] {
		return func(p Person) MatchResult {
			return MatchResult{
				Description: fmt.Sprintf("name is %s", p.Name),
				Matches:     p.Name == want,
				But:         fmt.Sprintf("name is %s", want),
				SubjectName: "Person",
			}
		}
	}

	someJSON := bytes.NewBuffer([]byte(`{"name": "John", "age": 42}`))

	Expect[io.Reader](t, someJSON).To(Parse[Person](HasName("John")))
	fmt.Println(t.Result())
	//Output: Test passed
}

func ExampleParse_fail() {
	t := &SpyTB{}

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	someJSON := bytes.NewBuffer([]byte(`invalid json`))

	Expect[io.Reader](t, someJSON).To(Parse[Person](EqualTo(Person{
		Name: "Pepper",
		Age:  14,
	})))
	fmt.Println(t.Result())
	//Output: Test failed: [expected JSON to be parseable into json.Person, but it could not be parsed: invalid character 'i' looking for beginning of value]
}
