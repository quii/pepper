# Pepper

Very much playing around, work in progress.

## Hamcrest-style matchers for Go

Like [gocrest](https://github.com/corbym/gocrest) but less mature and battle-tested. I'm making this purely to scratch an itch, but I do hope it is useful for some. The main purpose of writing this is for material for a chapter of [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests)

## Trade-offs and optimisations

### Type-safety

Now Go has generics, we now have a more expressive language which lets us make matchers that are type-safe, rather than relying on reflection. The problem with reflection is it can lead to lazy test writing, especially when you're dealing with complex types. It can lead developers into lazily asserting on complex types, which makes the tests harder to follow and more brittle. See the [curse of asserting on irrelevant detail](#the-curse-of-asserting-on-irrelevant-detail) for more on this.

The trade-off we're making here though is you will have to make your own matchers for your own types at times. This is a good thing, as it forces you to think about what you're actually testing, and it makes your tests more readable and less brittle. There will be plenty of examples to show how, and you can read the existing standard library of matchers to see how it's done.

### Composition

Matchers should be designed with composition in mind. For instance, let's take a look at the body matcher for an HTTP response:

```go
func HaveBody(bodyMatchers ...matching.Matcher[string]) matching.Matcher[*http.Response]
```

This allows the user to re-use string matchers that are already defined, but also lets you define your own matchers for the specific _kind_ of body you're interested in. Here's an example of some `matching.Matcher[string]` for a hypothetical JSON API that returns todo items:

```go
t.Run("example of matching JSON", func(t *testing.T) {
    type Todo struct {
        Name      string `json:"name"`
        Completed bool   `json:"completed"`
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
    

    t.Run("compose the matchers", func(t *testing.T) {
        res := httptest.NewRecorder()

        res.Body.WriteString(`{"name": "Egg", "completed": false}`)
        res.Header().Add("content-type", "application/json")

        Expect(t, res.Result()).To(
            HaveJSONHeader,
            HaveBody(WithTodoNameOf("Egg"), Not(WithCompletedTODO)),
        )
    })
})
```

In practice, your matchers should live outside your test code. Someone could argue writing these matchers adds more code as if that's a bad thing, but I would argue that the tests read far better, and don't suffer the problems you can run in to if you lazily assert on complex types. 

In my experience of using matchers, over time as you find yourself testing more and more permutations of behaviour, the effort behind the matchers pays off in terms of making tests easier to write, read and maintain. 

## Benefits of matchers

A lot of the time people zero-in on the "fluency" of matchers. Whilst it's true that the ease of use does make using matchers attractive, I think there's a larger, perhaps less obvious benefit. 

### The curse of asserting on irrelevant detail

A trap many fall in to when they write tests is they end up writing tests with [greedy assertions](https://quii.gitbook.io/learn-go-with-tests/meta/anti-patterns#asserting-on-irrelevant-detail) where you end up lazily writing tests where you check one complex object equals another. 

Often when we write a test, we only really care about the state of one field in a struct, yet when we are greedy, we end up coupling our tests to other data needlessly. This makes the test:

- Brittle, if domain logic changes elsewhere that happens to change values you weren't interested in, your test will fail
- Difficult to read, it's less obvious which effect you were hoping to exercise

Matchers allow you to write _domain specific_ matching code, focused on the _specific effects_ you're looking for. When used well, with a domain-centric, well-designed codebase, you tend to build a useful library of matchers that you can **re-use and compose** to write clear, consistently written, less brittle tests.
