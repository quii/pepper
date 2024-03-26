package fs

import (
	"github.com/quii/pepper/matching"
	"io"
	"io/fs"
)

const subjectName = "file system"

func HaveFileCalled(name string, contentsMatchers ...matching.Matcher[string]) matching.Matcher[fs.FS] {
	return func(fileSystem fs.FS) matching.MatchResult {
		file, err := fileSystem.Open(name)

		if err != nil {
			return matching.MatchResult{
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

		return matching.MatchResult{
			Description: "have file called " + name,
			Matches:     true,
			SubjectName: subjectName,
		}
	}
}
