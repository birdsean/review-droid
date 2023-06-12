package comments

import (
	"regexp"
	"strconv"
	"strings"
)

type Comment struct {
	Code        string
	CommentBody string
	FileAddress string
}

func ZipComment(segment, comments, filename string) ([]*Comment, error) {
	// create struct with keys: line number value: segment line
	lineReference := make(map[int]string)
	codeLines := strings.Split(segment, "\n")
	for _, line := range codeLines {
		lineNumber := regexp.MustCompile(`^(\d+)`).FindStringSubmatch(line)[1]
		lineInt, err := strconv.Atoi(lineNumber)
		if err != nil {
			return nil, err
		}
		lineReference[(lineInt)] = line
	}

	parsedComments := []*Comment{}
	splitComments := strings.Split(comments, "\n")
	for _, comment := range splitComments {
		if strings.Contains(comment, "No comment") {
			continue
		}

		lineNumber := "0"
		commentBody := comment
		match := regexp.MustCompile(`\[Line (\d+)\](.*)`).FindStringSubmatch(comment)
		if len(match) > 0 {
			lineNumber = match[1]
			commentBody = match[2]
		}

		lineInt, err := strconv.Atoi(lineNumber)
		if err != nil {
			return nil, err
		}

		// compile comment body and code
		comment := Comment{
			Code:        lineReference[lineInt],
			CommentBody: commentBody,
			FileAddress: filename,
		}
		parsedComments = append(parsedComments, &comment)
	}
	return parsedComments, nil
}
