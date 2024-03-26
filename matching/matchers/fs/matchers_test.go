package fs

import (
	. "github.com/quii/pepper/matching"
	"github.com/quii/pepper/matching/matchers/spytb"
	. "github.com/quii/pepper/matching/matchers/string"
	"testing"
	"testing/fstest"
)

func HaveFileCalled(name string, contentsMatchers ...Matcher[string]) Matcher[fstest.MapFS] {
	return func(fs fstest.MapFS) MatchResult {
		_, ok := fs[name]

		if !ok {
			return MatchResult{
				Description: "have file called " + name,
				Matches:     false,
				But:         "it did not",
				SubjectName: "file system",
			}
		}

		if len(contentsMatchers) > 0 {
			contents := string(fs[name].Data)
			for _, matcher := range contentsMatchers {
				result := matcher(contents)
				result.SubjectName = "file called " + name
				if !result.Matches {
					return result
				}
			}
		}

		return MatchResult{
			Description: "have file called " + name,
			Matches:     true,
			SubjectName: "file system",
		}
	}
}

func TestFSMatching(t *testing.T) {
	t.Run("FileContains", func(t *testing.T) {
		t.Run("file existence check", func(t *testing.T) {
			t.Run("passing", func(t *testing.T) {
				stubFS := fstest.MapFS{
					"someFile.txt": {
						Data: []byte("hello world"),
					},
				}

				Expect(t, stubFS).To(HaveFileCalled("someFile.txt"))
			})

			t.Run("failing", func(t *testing.T) {
				stubFS := fstest.MapFS{
					"someFile.txt": {
						Data: []byte("hello world"),
					},
				}

				spytb.VerifyFailingMatcher(
					t,
					stubFS,
					HaveFileCalled("anotherFile.txt"),
					"expected file system to have file called anotherFile.txt, but it did not",
				)
			})
		})
	})

	t.Run("FileContains with contents", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			stubFS := fstest.MapFS{
				"someFile.txt": {
					Data: []byte("hello world"),
				},
			}

			Expect(t, stubFS).To(HaveFileCalled("someFile.txt", HaveSubstring("world")))
		})

		t.Run("failing", func(t *testing.T) {
			stubFS := fstest.MapFS{
				"someFile.txt": {
					Data: []byte("hello world"),
				},
			}

			spytb.VerifyFailingMatcher(
				t,
				stubFS,
				HaveFileCalled("someFile.txt", HaveSubstring("goodbye")),
				`expected file called someFile.txt to contain "goodbye"`,
			)
		})
	})
}
