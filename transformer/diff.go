package transformer

import (
	"fmt"
	"regexp"
	"strings"
)

type DiffTransformer struct {
	rawDiff         string
	numberedRawDiff string
	fileDiffs       map[string]string
	fileSegments    map[string][]string
}

const segmentLength = 4000

func (dt *DiffTransformer) Transform(rawDiff string) {
	dt.rawDiff = rawDiff

	// Split diff into files
	dt.numberRawDiff()
	dt.splitIntoFiles()
	dt.generateSegments()
}

func (dt *DiffTransformer) GetFileSegments() map[string][]string {
	return dt.fileSegments
}

func (dt *DiffTransformer) splitIntoFiles() {
	fileDiffs := strings.Split(dt.numberedRawDiff, "diff --git")
	dt.fileDiffs = make(map[string]string)

	for _, file := range fileDiffs {
		// extract file name from line "+++ b/path/to/file" with regex
		match := regexp.MustCompile(`\+\+\+ b/(.*)`).FindStringSubmatch(file)
		if len(file) == 0 || len(match) == 0 {
			continue
		}

		fileName := match[1]

		// remove all lines that start with "+++", "---", or "@@"
		lines := strings.Split(file, "\n")
		for j := 0; j < len(lines); j++ {
			if strings.HasPrefix(lines[j], "+++") || strings.HasPrefix(lines[j], "---") || strings.HasPrefix(lines[j], "@@") {
				lines = append(lines[:j], lines[j+1:]...)
				j--
			} else {
				dt.fileDiffs[fileName] = strings.Join(lines, "\n")
			}
		}
	}
}

func (dt *DiffTransformer) generateSegments() {
	dt.fileSegments = make(map[string][]string)
	for filename, diff := range dt.fileDiffs {
		segments := []string{}
		splitDiff := strings.Split(diff, "\n")
		wordCount := 0
		for i, line := range splitDiff {
			wordCount += len(strings.Split(line, " "))
			if wordCount > segmentLength {
				newSegment := strings.Join(splitDiff[:i], "\n")
				segments = append(segments, newSegment)
				splitDiff = splitDiff[i:]
				wordCount = 0
			}
		}
		newSegment := strings.Join(splitDiff, "\n")
		dt.fileSegments[filename] = append(segments, newSegment)
	}
}

func (dt *DiffTransformer) numberRawDiff() {
	segments := []string{}
	for i, line := range strings.Split(dt.rawDiff, "\n") {
		segments = append(segments, fmt.Sprintf("%d %s", i+1, line))
	}
	dt.numberedRawDiff = strings.Join(segments, "\n")
}
