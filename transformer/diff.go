package transformer

import (
	"fmt"
	"strings"
)

type DiffTransformer struct {
	rawDiff   string
	fileDiffs []string
	segments  []string
}

const segmentLength = 4000

func (dt *DiffTransformer) Transform(rawDiff string) {
	dt.rawDiff = rawDiff

	// Split diff into files
	dt.splitIntoFiles()
	dt.generateSegments()
}

func (dt *DiffTransformer) GetLastSegment() string {
	return dt.segments[len(dt.segments)-1]
}

func (dt *DiffTransformer) GetSegments() []string {
	return dt.segments
}

func (dt *DiffTransformer) splitIntoFiles() {
	dt.fileDiffs = strings.Split(dt.rawDiff, "diff --git")
	// remove all empty strings from dt.fileDiffs
	for i := 0; i < len(dt.fileDiffs); i++ {
		dt.fileDiffs[i] = strings.TrimSpace(dt.fileDiffs[i])
		if dt.fileDiffs[i] == "" {
			dt.fileDiffs = append(dt.fileDiffs[:i], dt.fileDiffs[i+1:]...)
			i--
		}
	}
}

func (dt *DiffTransformer) generateSegments() {
	segments := []string{}
	for _, file := range dt.fileDiffs {
		splitDiff := strings.Split(file, "\n")
		wordCount := 0
		for i, line := range splitDiff {
			wordCount += len(strings.Split(line, " "))
			if wordCount > segmentLength {
				pendingSegment := splitDiff[:i]
				numbered := dt.numberLines(pendingSegment)
				newSegment := strings.Join(numbered, "\n")
				segments = append(segments, newSegment)
				splitDiff = splitDiff[i:]
				wordCount = 0
			}
		}
		segments = append(segments, strings.Join(splitDiff, "\n"))
	}
	dt.segments = segments
}

func (dt *DiffTransformer) numberLines(segments []string) []string {
	for i, segment := range segments {
		segments[i] = fmt.Sprintf("%d %s", i, segment)
	}
	return segments
}
