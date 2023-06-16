package entities

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v53/github"
)

type Comment struct {
	StartLine   int
	EndLine     int
	CommentBody string
	FileAddress string
	Side        string
}

func NewComments(segment, comments, filename string, debug bool) ([]*Comment, error) {
	parsedComments := []*Comment{}
	splitComments := strings.Split(comments, "\n")
	for _, comment := range splitComments {
		if strings.Contains(comment, "No comment") {
			continue
		}
		comment := generateComment(comment, segment, filename, debug)
		if comment == nil {
			continue
		}
		parsedComments = append(parsedComments, comment)
	}
	return parsedComments, nil
}

func FilterReplies(ghComments []*github.PullRequestComment) []*github.PullRequestComment {
	filteredComments := []*github.PullRequestComment{}
	for _, comment := range ghComments {
		if comment.GetInReplyTo() == 0 {
			filteredComments = append(filteredComments, comment)
		}
	}
	return filteredComments
}

func (comment *Comment) ToGithubComment(commitID string) *github.PullRequestComment {
	// Remove "a/" or "b/" from file address
	if comment.FileAddress[:2] == "a/" || comment.FileAddress[:2] == "b/" {
		comment.FileAddress = comment.FileAddress[2:]
	}

	ghComment := &github.PullRequestComment{
		Body:     github.String(comment.CommentBody),
		Path:     github.String(comment.FileAddress),
		CommitID: github.String(commitID),
		Side:     github.String(comment.Side),
		Line:     github.Int(comment.EndLine),
	}

	if comment.EndLine == 0 {
		ghComment.Line = github.Int(comment.StartLine)
	} else {
		ghComment.Line = github.Int(comment.StartLine)
	}

	return ghComment
}

func rangeStrToInts(rangeStr string) (int, int) {
	rangeInts := []int{}
	rangeSplit := strings.Split(rangeStr, "-")
	for _, val := range rangeSplit {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return 0, 0
		}
		rangeInts = append(rangeInts, intVal)
	}
	return rangeInts[0], rangeInts[1]
}

func extractCodeLineAndComment(rawComment string) (int, int, string, error) {
	// Detect any ranges of lines
	rangeMatch := regexp.MustCompile(`^\[.*?(\d+-\d+)\](.*)`).FindStringSubmatch(rawComment)
	if len(rangeMatch) > 0 {
		rangeStr := rangeMatch[1]
		commentBody := rangeMatch[2]
		lineStart, lineEnd := rangeStrToInts(rangeStr)
		if lineStart != 0 && lineEnd != 0 {
			return lineStart, lineEnd, commentBody, nil
		}
	}

	// Extract numbers out of square brackets
	match := regexp.MustCompile(`^\[.*?(\d+)(?:-\d+)?\](.*)`).FindStringSubmatch(rawComment)
	if len(match) == 0 {
		fmt.Printf("Failed to find line of code in rawComment (skipping) %s\n", rawComment)
		return 0, 0, "", fmt.Errorf("failed to find line of code in rawComment (skipping)")
	}

	lineNumber := match[1]
	commentBody := match[2]
	lineInt, err := strconv.Atoi(lineNumber)
	if err != nil {
		fmt.Printf("Failed to convert line number to int (skipping) %s\n", rawComment)
		return 0, 0, "", fmt.Errorf("failed to convert line number to int (skipping)")
	}
	return lineInt, 0, commentBody, nil
}

func extractSide(rawComment string) string {
	// detect if there is a plus or minus in the square brackets
	sideMatch := regexp.MustCompile(`^\[.*([+-]).*]`).FindStringSubmatch(rawComment)
	if len(sideMatch) < 2 {
		return "RIGHT" // default to right
	}

	// detect if any values in sideMatch equal +
	side := "LEFT"
	for _, val := range sideMatch {
		if strings.Contains(val, "+") {
			side = "RIGHT"
			break
		}
	}
	return side
}

func generateComment(rawComment, originalCode, filename string, debug bool) *Comment {
	startLine, endLine, commentBody, err := extractCodeLineAndComment(rawComment)
	if err != nil {
		return nil
	}

	side := extractSide(rawComment)
	body := strings.Trim(commentBody, " ")
	if debug {
		body = body + fmt.Sprintf(
			"\nDEBUG INFO\n[File] %s", filename) + fmt.Sprintf(
			"\n\n[Side] %s", side) + fmt.Sprintf(
			"\n\n[Start Line] %d", startLine) + fmt.Sprintf(
			"\n\n[End Line] %d", endLine) + fmt.Sprintf(
			"\n\n[Raw Comment] %s", rawComment) + fmt.Sprintf(
			"\n[\n[Original Code] %s", originalCode)
	}

	// compile comment body and code
	return &Comment{
		StartLine:   startLine,
		EndLine:     endLine,
		CommentBody: body,
		FileAddress: filename,
		Side:        side,
	}
}
