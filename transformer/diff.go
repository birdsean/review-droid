package transformer

import (
	"fmt"
	"regexp"
	"strconv"
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
		match := regexp.MustCompile(`[\+\-]{3} [a|b]/(.*)`).FindStringSubmatch(file)
		if len(file) == 0 || len(match) == 0 {
			continue
		}

		fileName := match[1]

		// remove all lines that start with "+++", "---", or "@@"
		lines := strings.Split(file, "\n")
		for j := 0; j < len(lines); j++ {
			if strings.Contains(lines[j], "+++") || strings.Contains(lines[j], "---") || strings.Contains(lines[j], "@@") {
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
	// number each line of the diff, looking at the header like "@@ -2,13 +2,15 @@" to know that the first line is line 2
	splitDiff := strings.Split(dt.rawDiff, "\n")
	rmLineNumber := 0
	addLineNumber := 0
	for i, line := range splitDiff {
		if strings.HasPrefix(line, "@@") {
			rmLineNumber, addLineNumber = getLineNumbers(line)
		} else if strings.HasPrefix(line, "-") {
			splitDiff[i] = fmt.Sprintf("%d %s", rmLineNumber, line)
			rmLineNumber++
		} else if strings.HasPrefix(line, "+") {
			splitDiff[i] = fmt.Sprintf("%d %s", addLineNumber, line)
			addLineNumber++
		} else {
			rmLineNumber++
			addLineNumber++
		}
	}
	dt.numberedRawDiff = strings.Join(splitDiff, "\n")
}

func getLineNumbers(line string) (int, int) {
	match := regexp.MustCompile(`@@ -(\d+),(\d+) \+(\d+),(\d+) @@`).FindStringSubmatch(line)
	if len(match) == 0 {
		return 0, 0
	}
	rmStart := match[1]
	addStart := match[3]
	rmStartInt, err := strconv.Atoi(rmStart)
	if err != nil {
		fmt.Println(err)
	}
	addStartInt, err := strconv.Atoi(addStart)
	if err != nil {
		fmt.Println(err)
	}
	return rmStartInt, addStartInt
}
