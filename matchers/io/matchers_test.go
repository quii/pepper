package io

import (
	"bytes"
	"fmt"
	. "github.com/quii/pepper"
	. "github.com/quii/pepper/matchers/comparable"
	"github.com/quii/pepper/matchers/spytb"
	"io"
	"testing"
)

func ExampleContainingByte() {
	t := &SpyTB{}

	buf := &bytes.Buffer{}
	buf.WriteString("hello")
	buf.WriteString("world")

	Expect[io.Reader](t, buf).To(HaveData(
		ContainingByte([]byte("hello")).And(ContainingByte([]byte("world"))),
	))
	fmt.Println(t.LastError())
	//Output:
}

func ExampleContainingByte_fail() {
	t := &SpyTB{}

	buf := &bytes.Buffer{}
	buf.WriteString("hello")
	buf.WriteString("world")

	Expect[io.Reader](t, buf).To(HaveData(
		ContainingByte([]byte("goodbye")),
	))
	fmt.Println(t.LastError())
	//Output: expected the reader to contain "goodbye", but it didn't have "goodbye"
}

func ExampleContainingString() {
	t := &SpyTB{}

	buf := &bytes.Buffer{}
	buf.WriteString("hello")
	buf.WriteString("world")

	Expect[io.Reader](t, buf).To(HaveData(
		ContainingString("world"),
	))
	fmt.Println(t.LastError())
	//Output:
}

func ExampleContainingString_fail() {
	t := &SpyTB{}

	buf := &bytes.Buffer{}
	buf.WriteString("hello")
	buf.WriteString("world")

	Expect[io.Reader](t, buf).To(HaveData(
		ContainingString("goodbye"),
	))
	fmt.Println(t.LastError())
	//Output: expected the reader to contain "goodbye", but it was "helloworld"
}

func ExampleHaveString() {
	t := &SpyTB{}

	buf := &bytes.Buffer{}
	buf.WriteString("hello")
	buf.WriteString("world")
	Expect[io.Reader](t, buf).To(HaveString(EqualTo("helloworld")))
	fmt.Println(t.LastError())
	//Output:
}

func ExampleHaveString_fail() {
	t := &SpyTB{}
	buf := &bytes.Buffer{}
	buf.WriteString("hello")
	buf.WriteString("world")
	Expect[io.Reader](t, buf).To(HaveString(EqualTo("Poo")))
	fmt.Println(t.LastError())
	//Output: expected "helloworld" to be equal to "Poo", but it was "helloworld"
}

func TestIOMatchers(t *testing.T) {
	t.Run("passing", func(t *testing.T) {
		buf := &bytes.Buffer{}
		buf.WriteString("hello")
		buf.WriteString("world")

		Expect[io.Reader](t, buf).To(HaveData(
			ContainingByte([]byte("hello")).And(ContainingByte([]byte("world"))),
		))

		buf.WriteString("goodbye")
		Expect[io.Reader](t, buf).To(HaveData(
			ContainingString("goodbye"),
		))
	})

	t.Run("failing", func(t *testing.T) {
		buf := &bytes.Buffer{}
		buf.WriteString("hello")
		buf.WriteString("world")

		spytb.VerifyFailingMatcher[io.Reader](
			t,
			buf,
			HaveData(ContainingByte([]byte("goodbye"))),
			`expected the reader to contain "goodbye", but it didn't have "goodbye"`,
		)

		spytb.VerifyFailingMatcher[io.Reader](
			t,
			buf,
			HaveData(ContainingString("goodbye")),
			`expected the reader to contain "goodbye", but it was ""`,
		)
	})
}
