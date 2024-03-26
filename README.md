# Pepper

[![Go Reference](https://pkg.go.dev/badge/github.com/quii/pepper.svg)](https://pkg.go.dev/github.com/quii/pepper)

##  _Type-safe_, composable, re-usable test matchers for Go

Inspired by [Hamcrest](https://hamcrest.org)

## Examples

Simple examples

```go
Expect(t, []string{"hello", "WORLD"}).To(ContainItem(HaveAllCaps))

Expect(spyTB, []string{"hello", "world"}).To(ContainItem(EqualTo("goodbye")))

Expect(t, "hello").To(HaveLength(5))
```

HTTP response examples

```go
t.Run("simple string match", func(t *testing.T) {
    res := httptest.NewRecorder()

    res.Body.WriteString("Hello, world")

    // see how we can compose and re-use matchers
    Expect(t, res.Result()).To(HaveBody(EqualTo("Hello, world")))
})

t.Run("compose http matchers", func(t *testing.T) {
    res := httptest.NewRecorder()

    res.Body.WriteString(`{"name": "Egg", "completed": false}`)
    res.Header().Add("content-type", "application/json")

    Expect(t, res.Result()).To(
        BeOK,
        HaveJSONHeader,
        HaveBody(WithTodoNameOf("Egg"), Not(WithCompletedTODO)),
    )
})
```

## Trade-offs and optimisations

### Type-safety

Now Go has generics, we now have a more expressive language which lets us make matchers that are type-safe, rather than relying on reflection. The problem with reflection is it can lead to lazy test writing, especially when you're dealing with complex types. Developers can lazily assert on complex types, which makes the tests harder to follow and more brittle. See the [curse of asserting on irrelevant detail](#the-curse-of-asserting-on-irrelevant-detail) for more on this.

The trade-off we're making here though is you will have to make your own matchers for your own types at times. This is a good thing, as it forces you to think about what you're actually testing, and it makes your tests more readable and less brittle. There will be plenty of examples to show how, and you can read the existing standard library of matchers to see how it's done. Due to the compositional nature of the library though, you _should_ be able to leverage existing matchers for re-use. 

### Composition

Matchers should be designed with composition in mind. For instance, let's take a look at the body matcher for an HTTP response:

```go
func HaveBody(bodyMatchers ...matching.Matcher[string]) matching.Matcher[*http.Response]
```

This allows the user to re-use string matchers that are already defined, but also lets you define your own matchers for the specific _kind_ of body you're interested in. Here's an example of some `matching.Matcher[string]` for a hypothetical JSON API that returns todo items:

```go
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
            But:         "it wasn't",
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

    t.Run("with a todo name", func(t *testing.T) {
        res := httptest.NewRecorder()
        res.Body.WriteString(`{"name": "Finish the side project", "completed": false}`)
        Expect(t, res.Result()).To(HaveBody(WithTodoNameOf("Finish the side project")))
    })

    t.Run("compose various matchers", func(t *testing.T) {
        res := httptest.NewRecorder()

        res.Body.WriteString(`{"name": "Egg", "completed": false}`)
        res.Header().Add("content-type", "application/json")

        Expect(t, res.Result()).To(
            BeOK,
            HaveJSONHeader,
            HaveBody(WithTodoNameOf("Egg"), Not(WithCompletedTODO)),
        )
    })
})
```

Note in the final example how we can compose built-in matchers like `BeOK`, `HaveJSONHeader` and `Not`, with the custom-built matchers to easily write very expressive tests that fail with very clear error messages.

In practice, your matchers should live outside your test code. Someone could argue writing these matchers adds more code as if that's a bad thing, but I would argue that the tests read far better, and don't suffer the problems you can run in to if you lazily assert on complex types. 

In my experience of using matchers, over time as you find yourself testing more and more permutations of behaviour, the effort behind the matchers pays off in terms of making tests easier to write, read and maintain. 

## Test failure readability

One of the most frustrating areas of working in a codebase with automated tests is how often test failure quality is poor. I'm sure every developer has into this scenario:

> `test_foo.go:123` - `true was not equal to false`

Computer, I already know that true is not equal to false. What was not false? What was true? What was the context? This is another reason why following TDD's process of "inspecting the failing test message" is so important, but sadly often overlooked.

This library should make it easy for you to write tests that give you a clear, concise message when they fail. Here's an example of a failing test:

```go
t.Run("failing test", func(t *testing.T) {
    res := httptest.NewRecorder()

    res.Body.WriteString(`{"name": "Bacon", "completed": false}`)
    res.Header().Add("content-type", "application/xml")

    Expect(t, res.Result()).To(
        BeOK,
        HaveJSONHeader,
        HaveBody(WithTodoNameOf("Egg"), Not(WithCompletedTODO)),
    )
})
```

Here is the failing output

```
=== RUN   TestHTTPTestMatchers/Body_matching/example_of_matching_JSON/compose_the_matchers
    matchers_test.go:100: expected the response to have header "content-type" of "application/json", but it was "application/xml"
    matchers_test.go:100: expected the response body to have a todo name of "Egg", but it was "Bacon"
```

Embracing this approach with well-written matchers means you get readable test failures for free. 

## Benefits of matchers

A lot of the time people zero-in on the "fluency" of matchers. Whilst it's true that the ease of use does make using matchers attractive, I think there's a larger, perhaps less obvious benefit. 

### The curse of asserting on irrelevant detail

A trap many fall in to when they write tests is they end up writing tests with [greedy assertions](https://quii.gitbook.io/learn-go-with-tests/meta/anti-patterns#asserting-on-irrelevant-detail) where you end up lazily writing tests where you check one complex object equals another. 

Often when we write a test, we only really care about the state of one field in a struct, yet when we are greedy, we end up coupling our tests to other data needlessly. This makes the test:

- Brittle, if domain logic changes elsewhere that happens to change values you weren't interested in, your test will fail
- Difficult to read, it's less obvious which effect you were hoping to exercise

Matchers allow you to write _domain specific_ matching code, focused on the _specific effects_ you're looking for. When used well, with a domain-centric, well-designed codebase, you tend to build a useful library of matchers that you can **re-use and compose** to write clear, consistently written, less brittle tests.

## Writing your own matchers

Writing your own matchers is a great way of building a re-usable, composable library for your domain code.

### How to

A matcher is a function that takes the _thing_ you're trying to match against, returning a result

```go
type Matcher[T any] func(T) MatchResult
```

Here is the definition of `MatchResult`

```go
type MatchResult struct {
    Description string
    Matches     bool
    But         string
    SubjectName string
}
```

Here is how `HaveAllCaps` is defined

```go
func HaveAllCaps(in string) matching.MatchResult {
	return matching.MatchResult{
		Description: "be in all caps",
		Matches:     strings.ToUpper(in) == in,
		But:         "it was not in all caps",
	}
}
```

#### Higher-order matchers

This is fine for simple matchers where you want to assert on a static property of `T`. Often though, you'll want to write matchers where you want to check a particular _ value of a property_. 

For this, no magic is required, just create a higher-order function that _returns_ a `Matcher[T]`. 

A simple example is with `EqualTo`

```go
func EqualTo[T comparable](in T) matching.Matcher[T] {
	return func(got T) matching.MatchResult {
		return matching.MatchResult{
			Description: fmt.Sprintf("be equal to %v", in),
			Matches:     got == in,
			But:         fmt.Sprintf("it was %v", got),
		}
	}
}
```

### Test support

It might be overkill to unit test matchers, but if you intend for them to be re-used a lot (like the library included in this package), it might be worth doing. I find it's a useful way to see in isolation the test output on failure. 

Thankfully, Pepper makes this easy

Check out some of the unit tests for some of the comparison matchers

```go
func TestComparisonMatchers(t *testing.T) {
	t.Run("Less than", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			Expect(t, 5).To(LessThan(6))
		})

		t.Run("failing", func(t *testing.T) {
			spytb.VerifyFailingMatcher(t, 6, LessThan(6), "expected 6 to be less than 6")
			spytb.VerifyFailingMatcher(t, 6, LessThan(3), "expected 6 to be less than 3")
		})
	})

	t.Run("Greater than", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			Expect(t, 5).To(GreaterThan(4))
		})

		t.Run("failing", func(t *testing.T) {
			spytb.VerifyFailingMatcher(t, 6, GreaterThan(6), "expected 6 to be greater than 6")
			spytb.VerifyFailingMatcher(t, 2, GreaterThan(10), "expected 2 to be greater than 10")
		})
	})
}
```
