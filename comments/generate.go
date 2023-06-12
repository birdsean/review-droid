package comments

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Comment struct {
	CodeLine    int
	CommentBody string
	FileAddress string
	Side        string
}

func ZipComment(segment, comments, filename string) ([]*Comment, error) {
	parsedComments := []*Comment{}
	splitComments := strings.Split(comments, "\n")
	for _, comment := range splitComments {
		if strings.Contains(comment, "No comment") {
			continue
		}
		comment := generateComment(comment, segment, filename)
		if comment == nil {
			continue
		}
		parsedComments = append(parsedComments, comment)
	}
	return parsedComments, nil
}

func generateComment(rawComment, originalCode, filename string) *Comment {
	// Extract numbers out of square brackets
	match := regexp.MustCompile(`^\[.*?(\d+)(?:-\d+)?\](.*)`).FindStringSubmatch(rawComment)
	if len(match) == 0 {
		fmt.Printf("Failed to find line of code in rawComment (skipping) %s\n", rawComment)
		return nil
	}

	lineNumber := match[1]
	commentBody := match[2]
	lineInt, err := strconv.Atoi(lineNumber)
	if err != nil {
		fmt.Printf("Failed to convert line number to int (skipping) %s\n", rawComment)
		return nil
	}

	// if has minus after line number like: 1 -package github, then side = LEFT
	// if has plus after line number like:  1 +package github_client, then side = RIGHT
	side := "RIGHT"
	col := strings.Split(originalCode, fmt.Sprintf("%d ", lineInt))[1]
	if strings.HasPrefix(col, "-") {
		side = "LEFT"
	}

	// compile comment body and code
	return &Comment{
		CodeLine:    lineInt,
		CommentBody: strings.Trim(commentBody, " "),
		FileAddress: filename,
		Side:        side,
	}
}
