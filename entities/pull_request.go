package entities

import (
	"fmt"
	"log"
	"os"

	"github.com/birdsean/review-droid/clients"
	"github.com/google/go-github/v53/github"
)

var DEBUG = os.Getenv("DEBUG") == "true"

type PullRequest struct {
	core         *github.PullRequest
	githubClient clients.GithubRepoClient
}

func GetAllPullRequests(client clients.GithubRepoClient) ([]*PullRequest, error) {
	prs, err := client.GetPrs()
	if err != nil {
		return nil, err
	}
	wrappedPrs := []*PullRequest{}
	for _, pr := range prs {
		wrappedPrs = append(wrappedPrs, NewPullRequest(pr, client))
	}
	return wrappedPrs, nil
}

func NewPullRequest(core *github.PullRequest, client clients.GithubRepoClient) *PullRequest {
	return &PullRequest{core, client}
}

func (pr *PullRequest) Comment(comments []*Comment) error {
	commitId := pr.core.GetHead().GetSHA()
	fmt.Printf("Posting %d comments to PR #%d\n", len(comments), pr.core.GetNumber())
	for _, comment := range comments {
		ghComment := comment.ToGithubComment(commitId)
		err := pr.githubClient.PostComment(pr.core.GetNumber(), ghComment)
		if err != nil {
			fmt.Printf("Failed to post comment: %v\n", err)
			// TODO post comment to entire file if failed on a line.
		}
	}
	return nil
}

func (pr *PullRequest) Review() []*Comment {
	fmt.Printf("PR #%d: %s\n", pr.core.GetNumber(), pr.core.GetTitle())

	rawDiff, err := pr.githubClient.GetPrDiff(pr.core.GetNumber())
	if err != nil {
		log.Fatalf("Failed to get raw diff: %v", err)
	}

	diff := NewDiff(rawDiff)
	fileSegments := diff.GetFileSegments()
	allComments := []*Comment{}

	fmt.Printf("Getting comments for %d segments\n", len(fileSegments))
	for filename, segments := range fileSegments {
		for _, segment := range segments {
			reviewStatements := generateReviewStatements(segment)
			if reviewStatements == nil {
				continue
			}

			comments, err := NewComments(segment, *reviewStatements, filename, DEBUG)
			if err != nil {
				log.Fatalf("Failed to zip comment: %v", err)
			}

			allComments = append(allComments, comments...)
			if DEBUG && len(allComments) >= 1 {
				break
			}
		}
		if DEBUG && len(allComments) >= 1 {
			break
		}
	}

	return allComments
}

func (pr *PullRequest) Evaluate() error {
	critic := NewCritic(pr.githubClient)
	return critic.EvaluateReviewQuality(pr.core.GetNumber())
}

func generateReviewStatements(segment string) *string {
	openAiClient := clients.NewOpenAiClient()
	completion, err := openAiClient.GetCompletion(segment, DEBUG)
	if err != nil {
		fmt.Printf("Failed to get completion: %v\n", err)
	}

	if DEBUG && completion != nil {
		fmt.Println("********************")
		fmt.Println(*completion)
		fmt.Println("********************")
	}
	return completion
}
