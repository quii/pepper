package fs

import (
	"github.com/quii/pepper"
	"io"
	"io/fs"
)

const subjectName = "file system"

// HaveFileCalled checks if a file exists in the file system, and can run additional matchers on its contents.
func HaveFileCalled(name string, contentMatcher ...pepper.Matcher[string]) pepper.Matcher[fs.FS] {
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
			all, _ := io.ReadAll(file)
			contents := string(all)
			for _, matcher := range contentMatcher {
				result := matcher(contents)
				result.SubjectName = "file called " + name
				if !result.Matches {
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
