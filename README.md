# Pepper

## Hamcrest-style matchers for Go

Like [gocrest](https://github.com/corbym/gocrest) but less mature and battle-tested. I'm making this purely to scratch an itch, but I do hope it is useful for some. The main purpose of writing this is for material for a chapter of [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests)

## Benefits of matchers

A lot of the time people zero-in on the "fluency" of matchers. Whilst it's true that the ease of use does make using matchers attractive, I think there's a larger, perhaps less obvious benefit. 

### The curse of asserting on irrelevant detail

A trap many fall in to when they write tests is they end up writing tests with [greedy assertions](https://quii.gitbook.io/learn-go-with-tests/meta/anti-patterns#asserting-on-irrelevant-detail) where you end up lazily writing tests where you check one complex object equals another. 

Often when we write a test, we only really care about the state of one field in a struct, yet when we are greedy, we end up coupling our tests to other data needlessly. This makes the test:

- Brittle, if domain logic changes elsewhere that happens to change values you weren't interested in, your test will fail
- Difficult to read, it's less obvious which effect you were hoping to exercise

Matchers allow you to write _domain specific_ matching code, focused on the _specific effects_ you're looking for. When used well, with a domain-centric, well-designed codebase, you tend to build a useful library of matchers that you can **re-use and compose** to write clear, consistently written, less brittle tests. 

### TODO: Worked example

- Include "zeroing time" as an example