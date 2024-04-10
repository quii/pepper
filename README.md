# Pepper

##  _Type-safe_, composable, extensible test matching for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/quii/pepper.svg)](https://pkg.go.dev/github.com/quii/pepper)
[![Go Report Card](https://goreportcard.com/badge/github.com/quii/pepper)](https://goreportcard.com/report/github.com/quii/pepper)
![Test suite status](https://github.com/quii/pepper/actions/workflows/test.yaml/badge.svg)

Inspired by [Hamcrest](https://hamcrest.org)

Out of the box, Pepper can work just like other test libraries in Go like [is](https://github.com/matryer/is), where you can make basic assertions on data.

### Simple examples

```go
Expect(t, "Pepper").To(Equal("Pepper"))
```

And Pepper has lots of built-in **matchers**, which you pass in to `To`, for common testing operations such as examining `comparable`, `string`, `io` and `*http.Response`.

What is a `Matcher[T]` ? It's a function that takes a `T`, and returns a `MatchResult`

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

```go
Expect(t, "HELLO").To(HaveAllCaps)
```

Quite nice, but still, not all that different from libraries you already use. Pepper starts to come into its own when you start taking advantage of _composing matchers_.

### Composing matchers

```go
Expect(t, score).To(GreaterThan(5).And(LessThan(10)))
```

The method `And` on `Matcher[T]`, lets you compose matchers. `And` _returns_ the composed `Matcher[T]`, so you can continue to chain more matchers however you like.

```go
Expect(t, score).To(GreaterThan(5).And(Not(GreaterThan(opponentScore))))
```

`Not` negates a matcher. By using matchers, and composing them with `And`, `Not`, `Or`, you can write very expressive tests, cheaply.

### Defining your own matchers

You can define your own matchers for your own types. Over time, the investment in writing matchers for your tests pays dividends, the cost of writing your tests decrease, as you reuse, mix and match the standard matchers and composition tools with your own. 

Some will argue writing these matchers adds more code as if that's inherently a bad thing, but I would argue that the tests read far better, and don't suffer the problems you can run in to if you lazily assert on complex types.

In my experience of using matchers, over time as you find yourself testing more and more permutations of behaviour, the effort behind the matchers pays off in terms of making tests easier to write, read and maintain.

Here is an example of testing a todo-list

```go
type Todo struct {
    Name        string    `json:"name"`
    Completed   bool      `json:"completed"`
    LastUpdated time.Time `json:"last_updated"`
}

func WithCompletedTODO(todo Todo) MatchResult {
    return MatchResult{
        Description: "have a completed todo",
        Matches:     todo.Completed,
        But:         "it wasn't complete",
    }
}
func WithTodoNameOf(todoName string) Matcher[Todo] {
    return func(todo Todo) MatchResult {
        return MatchResult{
            Description: fmt.Sprintf("have a todo name of %q", todoName),
            Matches:     todo.Name == todoName,
            But:         fmt.Sprintf("it was %q", todo.Name),
        }
    }
}

func TestTodos(t *testing.T) {
    t.Run("with completed todo", func(t *testing.T) {
        res := httptest.NewRecorder()
        res.Body.WriteString(`{"name": "Finish the side project", "completed": true}`)
        Expect(t, res.Result()).To(HaveBody(Parse[Todo](WithCompletedTODO)))
    })

    t.Run("with a todo name", func(t *testing.T) {
        res := httptest.NewRecorder()
        res.Body.WriteString(`{"name": "Finish the side project", "completed": false}`)
        Expect(t, res.Result()).To(HaveBody(Parse[Todo](WithTodoNameOf("Finish the side project"))))
    })

    t.Run("compose the matchers", func(t *testing.T) {
        res := httptest.NewRecorder()

        res.Body.WriteString(`{"name": "Egg", "completed": false}`)
        res.Header().Add("content-type", "application/json")

        Expect(t, res.Result()).To(
            BeOK,
            HaveJSONHeader,
            HaveBody(Parse[Todo](WithTodoNameOf("Egg").And(Not(WithCompletedTODO)))),
        )
    })
})
```

Note how we can compose built-in matchers like `BeOK`, `HaveJSONHeader` and `Not`, with the custom-built matchers to easily write very expressive tests that fail with very clear error messages. Pepper makes it **really easy** to check JSON responses of your HTTP handlers.

Also note, due to the compositional nature of `Matcher[T]`, we can re-use our `Matcher[Todo]` for tests at different abstraction levels; these matchers are not coupled to HTTP, we _composed_ the matchers for this context. For instance, if you have a `TodoRepository`, you could use _these same matchers_ in the tests for that too. 

### Test failure readability

One of the most frustrating areas working with automated tests is how often test failure quality is poor. I'm sure every developer has into this scenario:

> `test_foo.go:123` - `true was not equal to false`

Computer, I already know that true is not equal to false. What was not false? What was true? What was the context? 

Pepper makes it easy for you to write tests that give you a clear message when they fail.

```go
t.Run("failure message", func(t *testing.T) {
    res := httptest.NewRecorder()
    res.WriteHeader(http.StatusNotFound)
    Expect(t, res.Result()).To(BeOK)
})
```

Here is the failing output

```
=== RUN   TestHTTPTestMatchers/Status_code_matchers/OK/failure_message
    matchers_test.go:292: expected the response to have status of 200, but it was 404
```

Embracing this approach with well-written matchers means you get readable test failures for free.

### Summary

Pepper brings the following to the table

- ✅ Type-safe tests. No `interface{}`
- ✅ Composition to reduce boilerplate
- ✅ Clear test output as a first-class citizen
- ✅ A "standard library" of matchers to let you quickly write expressive tests for common scenarios out of the box
- ✅ Extensibility. You can write rich, re-useable matchers for your domain to help you write high quality, low maintenance tests
- ❌ **Not a framework**. Pepper does not dictate how you set up your tests or how you design your code. All it does is help you write better assertions with less effort. 


## Examples

You can find high-level examples in the [GoDoc](https://pkg.go.dev/github.com/quii/pepper#pkg-examples). 

The matchers, which you can find in the [directories section](https://pkg.go.dev/github.com/quii/pepper#pkg-examples) also have examples.

You'll notice in the examples the following line:

```go
t := &SpyTB{}
```

This is a test spy that is used to verify the output of the matches made. The examples call `LastError()` to see what test output would happen, so you can see what the failures look like. 

Finally, I have worked through my course [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests), using Pepper to write assertions, so you can find more examples [at the Github repository](https://github.com/quii/learn-go-with-pepper)

## Trade-offs and optimisations

### Type-safety

Now Go has generics, we now have a more expressive language which lets us make matchers that are type-safe, rather than relying on reflection. The problem with reflection is it can lead to lazy test writing, especially when you're dealing with complex types. Developers can lazily assert on complex types, which makes the tests harder to follow and more brittle. See the [curse of asserting on irrelevant detail](#the-curse-of-asserting-on-irrelevant-detail) for more on this.

The trade-off we're making here though is you will have to make your own matchers for your own types at times. This is a good thing, as it forces you to think about what you're actually testing, and it makes your tests more readable and less brittle. There will be plenty of examples to show how, and you can read the existing standard library of matchers to see how it's done. Due to the compositional nature of the library though, you _should_ be able to leverage existing matchers for re-use.

Due to Go having some constraints on _where_ you can use generics, such as function types not being allowed to have type parameters, the API isn't as friendly as it would be, if you used `interface{}`/`any`. However, this is a trade-off I am OK with, in the name of type-safety. 

### Composition and re-use

Matchers should be designed with composition in mind. For instance, let's take a look at the body matcher for an HTTP response:

```go
func HaveBody(bodyMatchers pepper.Matcher[io.Reader]) pepper.Matcher[*http.Response]
```

This allows the user to re-use `io.Reader` matchers that are already defined, compose them with `And`/`Or`/`Not`, and of course users can define their own `Matcher[io.Reader]` too. 

### Test failure readability

Often the most expensive part of a test suite is trying to understand _what_ has failed when a test goes red. This is why TDD emphasises the first step of writing a failing test and inspecting the output. It's a chance for you to see what it's like if the test fails 6 months later, and you've lost all context. Pepper strives to make it easy for you to write tests that explain exactly what has gone wrong.

Please note though that this library will not bend over backwards to write _perfect_ English. It's important that the reason for test failure is clear, but perfect grammar is not needed for this; and the complexity cost involved to make matchers "write" different sentences depending on how they are used, is not worth it. 

## Benefits of matchers

A lot of the time people zero-in on the "fluency" of matchers. Whilst it's true that the ease of use does make using matchers attractive, I think there's a larger, perhaps less obvious benefit. 

### The curse of asserting on irrelevant detail

A trap many fall in to when they write tests is they end up writing tests with [greedy assertions](https://quii.gitbook.io/learn-go-with-tests/meta/anti-patterns#asserting-on-irrelevant-detail) where you end up lazily writing tests where you check one complex object equals another. 

Often when we write a test, we only really care about the state of one field in a struct, yet when we are greedy, we end up coupling our tests to other data needlessly. This makes the test:

- Brittle, if domain logic changes elsewhere that happens to change values you weren't interested in, your test will fail
- Difficult to read, it's less obvious which effect you were hoping to exercise

Matchers allow you to write _domain specific_ matching code, focused on the _specific effects_ you're looking for. When used well, with a domain-centric, well-designed codebase, you tend to build a useful library of matchers that you can **re-use and compose** to write clear, consistently written, less brittle tests.

## How to write your own matchers

Reminder: a matcher is a function that takes the _thing_ you're trying to match against, returning a result

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

This is fine for simple matchers where you want to assert on a static property of `T`. Often though, you'll want to write matchers where you want to check a particular _ value of a property_ .

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

#### Leveraging composition

When designing your higher-order matchers, think about how the value you are matching against could be matched with other matchers you may not have thought of. 

For example, `HasSize`, I could've written like this:

```go
func HaveSize[T any](size int) pepper.Matcher[[]T]
```

However, this needlessly couples the matcher to the specific matching I was currently catering for (logically `EqualTo`). Instead, we can design our matcher to be combined with other matchers that work on `int`.

```go
func HaveSize[T any](matcher pepper.Matcher[int]) pepper.Matcher[[]T] {
	return func(items []T) pepper.MatchResult {
		return matcher(len(items))
	}
}
```

This way, users can use this matcher in different ways, like checking if a slice has a size `LessThan(5)` or `GreaterThan(3)`.

With this simple change, users can leverage the _other_ composition tools like `And`:

```go
Expect(t, catsInAHotel).To(HaveSize(GreaterThan(3).And(LessThan(10))))
```

_Tip_: When designing your matcher, consider changing the argument(s) from `T` to `Matcher[T]`. 


### Test support

Pepper makes testing matchers easy because you inject in the testing framework into `Expect`, so we can _spy_ on it. 

I have found writing [testable examples](https://go.dev/blog/examples) though to be a satisfying way of both documenting and testing matchers.

```go
func ExampleContainItem() {
	t := &SpyTB{}

	anArray := []string{"HELLO", "WORLD"}
	Expect(t, anArray).To(ContainItem(HaveAllCaps))

	fmt.Println(t.LastError())
	//Output:
}

func ExampleContainItem_fail() {
	t := &SpyTB{}

	anArray := []string{"hello", "world"}
	Expect(t, anArray).To(ContainItem(HaveAllCaps))

	fmt.Println(t.LastError())
	//Output: expected [hello world] to contain an item in all caps, but it did not
}
```

However, if you wish to check multiple scenarios, polluting the go doc with lots of examples may not be appropriate, in which case, write some unit tests.

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

### Contributing your own matchers

If you have a matcher you think would be useful to others, please consider contributing it to this library. 

**Please only submit matchers that work against types in the standard library**. This keeps the library focused and backward compatible. It would be fantastic if over time this library matured into a rich suite of matchers so any dev can pick up Go and start writing excellent tests against the standard library, which already gets you so far in terms of getting work done. 

Your PR will need the following

- At least two [testable examples](https://go.dev/blog/examples), one showing the matcher passing (with an empty output) and one showing the matcher failing with the expected failing output. This will help users understand how to use your matcher. 
- Automated tests in general
- Go doc comments for the matcher

As discussed above, try to keep them "open" in terms of their design, so they can be composed with other matchers. 
