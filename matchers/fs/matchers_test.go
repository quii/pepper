package fs

import (
	"fmt"
	. "github.com/quii/pepper"
	"github.com/quii/pepper/matchers/spytb"
	. "github.com/quii/pepper/matchers/string"
	"io/fs"
	"testing"
	"testing/fstest"
)

func ExampleHaveFileCalled() {
	t := &SpyTB{}
	stubFS := fstest.MapFS{
		"someFile.txt": {
			Data: []byte("hello world"),
		},
	}

	Expect[fs.FS](t, stubFS).To(HaveFileCalled("someFile.txt", HaveSubstring("Pluto")))

	fmt.Println(t.LastError())
	//Output: expected file called someFile.txt to contain "Pluto"
}

func ExampleHaveDir() {
	t := &SpyTB{}
	stubFS := fstest.MapFS{
		"someDir": {
			Mode: fs.ModeDir,
		},
	}

	Expect[fs.FS](t, stubFS).To(HaveDir("someDir"))

	fmt.Println(t.LastError())
	//Output:
}

func ExampleHaveDir_withErr() {
	t := &SpyTB{}
	stubFS := fstest.MapFS{
		"someFile.txt": {
			Data: []byte("hello world"),
		},
	}

	Expect[fs.FS](t, stubFS).To(HaveDir("someFile.txt"))

	fmt.Println(t.LastError())
	//Output: expected file system to have directory called "someFile.txt", but it was not a directory
}

func TestFSMatching(t *testing.T) {
	stubFS := fstest.MapFS{
		"someFile.txt": {
			Data: []byte("hello world"),
		},
		"someDir": {
			Mode: fs.ModeDir,
		},
		"nested/someFile.txt": {
			Data: []byte("hello world"),
		},
	}

	t.Run("HasDir", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			Expect[fs.FS](t, stubFS).To(HaveDir("someDir"))
		})

		t.Run("failing", func(t *testing.T) {
			spytb.VerifyFailingMatcher[fs.FS](
				t,
				stubFS,
				HaveDir("someFile.txt"),
				`expected file system to have directory called "someFile.txt", but it was not a directory`,
			)
			spytb.VerifyFailingMatcher[fs.FS](
				t,
				stubFS,
				HaveDir("non-existent-file"),
				`expected file system to have directory called "non-existent-file", but it did not`,
			)
		})
	})

	t.Run("FileContains", func(t *testing.T) {
		t.Run("file existence check", func(t *testing.T) {
			t.Run("passing", func(t *testing.T) {
				Expect[fs.FS](t, stubFS).To(HaveFileCalled("someFile.txt"))
				Expect[fs.FS](t, stubFS).To(HaveFileCalled("nested/someFile.txt"))
			})
		})
	})

	t.Run("FileContains with contents", func(t *testing.T) {
		t.Run("passing", func(t *testing.T) {
			Expect[fs.FS](t, stubFS).To(HaveFileCalled("someFile.txt", HaveSubstring("world")))
		})

		t.Run("failing", func(t *testing.T) {
			spytb.VerifyFailingMatcher[fs.FS](
				t,
				stubFS,
				HaveFileCalled("someFile.txt", HaveSubstring("goodbye")),
				`expected file called someFile.txt to contain "goodbye"`,
			)

			t.Run("failing", func(t *testing.T) {
				failingFS := FailToReadFS{Error: fmt.Errorf("could not read file")}
				spytb.VerifyFailingMatcher[fs.FS](
					t,
					failingFS,
					HaveFileCalled("anotherFile.txt", HaveSubstring("BLAH")),
					"expected file system to have file called anotherFile.txt, but it could not be read",
				)
			})
		})
	})
}

type FailToReadFS struct {
	Error error
}

func (f FailToReadFS) Open(name string) (fs.File, error) {
	return FailingFile(f), nil
}

type FailingFile struct {
	Error error
}

func (f FailingFile) Stat() (fs.FileInfo, error) {
	return nil, f.Error
}

func (f FailingFile) Read(bytes []byte) (int, error) {
	return 0, f.Error
}

func (f FailingFile) Close() error {
	return f.Error
}
