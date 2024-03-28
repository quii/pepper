package fs

import (
	"github.com/quii/pepper"
	"io"
	"io/fs"
)

const subjectName = "file system"

func HaveFileCalled(name string, contentsMatchers ...pepper.Matcher[string]) pepper.Matcher[fs.FS] {
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

		if len(contentsMatchers) > 0 {
			all, _ := io.ReadAll(file)
			contents := string(all)
			for _, matcher := range contentsMatchers {
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
