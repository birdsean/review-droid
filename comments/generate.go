package comments

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Comment struct {
	Code        string
	CommentBody string
	FileAddress string
}

func ZipComment(segment string, comments string) ([]*Comment, error) {
	// create struct with keys: line number value: segment line
	lineReference := make(map[int]string)
	codeLines := strings.Split(segment, "\n")
	for _, line := range codeLines {
		// extract line number with regex. ex: "89 +    return &resp.Choices[0].Message.Content, nil"
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

		match := regexp.MustCompile(`\[Line (\d+)\](.*)`).FindStringSubmatch(comment)
		lineNumber := match[1]
		commentBody := match[2]
		lineInt, err := strconv.Atoi(lineNumber)
		if err != nil {
			return nil, err
		}

		// compile comment body and code
		comment := Comment{
			Code:        lineReference[lineInt],
			CommentBody: commentBody,
			FileAddress: "TODO",
		}
		parsedComments = append(parsedComments, &comment)
	}
	return parsedComments, nil
}

func Main() {
	// Read results.json.test JSON file
	results := []string{}
	file, err := ioutil.ReadFile("results.json.test")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	err = json.Unmarshal(file, &results)
	if err != nil {
		log.Fatalf("Failed to unmarshal file: %v", err)
	}

}
