package fs

import (
	"fmt"
	"github.com/quii/pepper"
	"io"
	"io/fs"
)

const subjectName = "file system"

// HaveFileCalled checks if a file exists in the file system, and can run additional matchers on its contents.
func HaveFileCalled(name string, contentMatcher ...pepper.Matcher[io.Reader]) pepper.Matcher[fs.FS] {
	return func(fileSystem fs.FS) pepper.MatchResult {
		file, err := fileSystem.Open(name)

		if err != nil {
			return pepper.MatchResult{
				Description: "have file called " + name,
				Matches:     false,
				But:         "it did not",
				SubjectName: subjectName,
			}
		}

		defer file.Close()

		if len(contentMatcher) > 0 {
			for _, matcher := range contentMatcher {
				result := matcher(file)
				result.SubjectName = fmt.Sprintf("file called %q", name)
				if !result.Matches {
					if result.But == "" {
						result.But = "while the file existed, the contents did not match"
					}
					return result
				}
			}
		}

		return pepper.MatchResult{
			Description: "have file called " + name,
			Matches:     true,
			SubjectName: subjectName,
		}
	}
}

// HaveDir checks if a directory exists in the file system.
func HaveDir(name string) pepper.Matcher[fs.FS] {
	return func(fileSystem fs.FS) pepper.MatchResult {
		f, err := fileSystem.Open(name)

		result := pepper.MatchResult{
			Description: fmt.Sprintf("have directory called %q", name),
			SubjectName: subjectName,
			Matches:     true,
		}

		if err != nil {
			result.Matches = false
			result.But = "it did not"
			return result
		}

		stat, err := f.Stat()
		if err != nil {
			result.Matches = false
			result.But = "it could not be read"
			return result
		}

		if !stat.IsDir() {
			result.Matches = false
			result.But = "it was not a directory"
			return result
		}

		return result
	}
}
