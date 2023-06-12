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

	// detect if there is a plus or minus in the square brackets
	sideMatch := regexp.MustCompile(`^\[.*([+-]).*]`).FindStringSubmatch(rawComment)
	if len(sideMatch) < 2 {
		// find first + or - in original code
		plusIdx := strings.Index(originalCode, "+")
		minusIdx := strings.Index(originalCode, "-")
		if plusIdx == -1 && minusIdx == -1 {
			fmt.Printf("Failed to find side of code in rawComment (skipping) %s\n", rawComment)
			return nil
		}
		if plusIdx == -1 {
			sideMatch = []string{"", "", "-"}
		}
		if minusIdx == -1 {
			sideMatch = []string{"", "", "+"}
		}
	}

	sideSymbol := sideMatch[1]
	side := "RIGHT"
	if sideSymbol == "-" {
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
