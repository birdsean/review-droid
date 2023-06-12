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

const segmentLength = 2000

func (dt *DiffTransformer) Transform(rawDiff string) {
	dt.rawDiff = rawDiff

	// Split diff into files
	dt.fileDiffs = strings.Split(dt.rawDiff, "diff --git")
	dt.generateSegments()
}

func (dt *DiffTransformer) generateSegments() {
	segments := []string{}
	for _, file := range dt.fileDiffs {
		splitDiff := strings.Split(file, "\n")
		charCount := 0
		for i, line := range splitDiff {
			charCount += len(line)
			if charCount > segmentLength {
				pendingSegment := splitDiff[:i]
				numbered := dt.numberLines(pendingSegment)
				newSegment := strings.Join(numbered, "\n")
				segments = append(segments, newSegment)
				splitDiff = splitDiff[i:]
				charCount = 0
			}
			splitDiff[i] = fmt.Sprintf("%d %s", i, line)
		}
		segments = append(segments, strings.Join(splitDiff, "\n"))
	}
	dt.segments = segments
}

func (dt *DiffTransformer) numberLines(segments []string) []string {
	for i, segment := range segments {
		splitDiff := strings.Split(segment, "\n")
		for i, line := range splitDiff {
			splitDiff[i] = fmt.Sprintf("%d %s", i, line)
		}
		segments[i] = strings.Join(splitDiff, "\n")
	}
	return segments
}
