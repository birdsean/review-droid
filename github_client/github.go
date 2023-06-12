package github_client

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/birdsean/review-droid/comments"
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

type GithubRepoClient struct {
	client *github.Client
	ctx    context.Context
	owner  string
	repo   string
}

func (grc *GithubRepoClient) Init() {
	// Set up authentication using a personal access token
	token := os.Getenv("REVIEW_DROID_TOKEN")
	if token == "" {
		log.Fatalf("REVIEW_DROID_TOKEN environment variable must be set")
	}
	grc.ctx = context.Background()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(grc.ctx, ts)

	client := github.NewClient(tc)
	grc.client = client

	grc.owner = os.Getenv("GITHUB_OWNER")
	if grc.owner == "" {
		log.Fatalf("GITHUB_OWNER environment variable must be set")
	}

	grc.repo = os.Getenv("GITHUB_REPO")
	if grc.repo == "" {
		log.Fatalf("GITHUB_REPO environment variable must be set")
	}
}

func (grc *GithubRepoClient) GetPrs() ([]*github.PullRequest, error) {

	// List pull requests for the specified repository
	prs, _, err := grc.client.PullRequests.List(grc.ctx, grc.owner, grc.repo, &github.PullRequestListOptions{})
	return prs, err
}

func (grc *GithubRepoClient) GetPrDiff(pr *github.PullRequest) (string, error) {
	diff, _, err := grc.client.PullRequests.GetRaw(grc.ctx, grc.owner, grc.repo, pr.GetNumber(), github.RawOptions{Type: github.Diff})
	return diff, err
}

func (grc *GithubRepoClient) PostComment(pr *github.PullRequest, comment *github.PullRequestComment) error {
	_, _, err := grc.client.PullRequests.CreateComment(grc.ctx, grc.owner, grc.repo, pr.GetNumber(), comment)
	return err
}

func (grc *GithubRepoClient) ParsedCommentToGithubComment(parsed *comments.Comment, commitID string) *github.PullRequestComment {
	// Remove "a/" or "b/" from file address
	if parsed.FileAddress[:2] == "a/" || parsed.FileAddress[:2] == "b/" {
		parsed.FileAddress = parsed.FileAddress[2:]
	}

	fmt.Printf("Line of Code: %d\n", parsed.CodeLine)
	fmt.Printf("File address: %s\n", parsed.FileAddress)

	comment := &github.PullRequestComment{
		Body:     github.String(parsed.CommentBody),
		Path:     github.String(parsed.FileAddress),
		CommitID: github.String(commitID),
		Side:     github.String(parsed.Side),
		Line:     github.Int(parsed.CodeLine),
	}

	return comment
}
