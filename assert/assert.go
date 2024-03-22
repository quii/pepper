package assert

import "testing"

func NoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Error(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error")
	}
}

func Equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func SliceEqual[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length of slices different: got %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("got  %v\n want %v", got, want)
		}
	}
}

func NotEqual[T comparable](t *testing.T, got, notWanted T) {
	t.Helper()
	if got == notWanted {
		t.Fatalf("didn't want this to equal %v", notWanted)
	}
}

func True(t *testing.T, got bool, msg ...string) {
	t.Helper()
	if !got {
		t.Fatal("expected true", msg)
	}
}

func False(t *testing.T, got bool, msg ...string) {
	t.Helper()
	if got {
		t.Fatal("expected false", msg)
	}
}
