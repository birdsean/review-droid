package clients

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

type GithubRepoClient struct {
	client *github.Client
	ctx    context.Context
	owner  string
	repo   string
}

// constructor
func NewGithubRepoClient() GithubRepoClient {
	grc := GithubRepoClient{}
	grc.Init()
	return grc
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

func (grc *GithubRepoClient) GetPrDiff(prId int) (string, error) {
	diff, _, err := grc.client.PullRequests.GetRaw(grc.ctx, grc.owner, grc.repo, prId, github.RawOptions{Type: github.Diff})
	return diff, err
}

func (grc *GithubRepoClient) PostComment(prId int, comment *github.PullRequestComment) error {
	_, _, err := grc.client.PullRequests.CreateComment(grc.ctx, grc.owner, grc.repo, prId, comment)
	return err
}

func (grc *GithubRepoClient) GetPrComments(prId int) ([]*github.PullRequestComment, error) {
	comments, _, err := grc.client.PullRequests.ListComments(grc.ctx, grc.owner, grc.repo, prId, &github.PullRequestListCommentsOptions{})
	return comments, err
}

func (grc *GithubRepoClient) DeleteComment(comment *github.PullRequestComment) error {
	_, err := grc.client.PullRequests.DeleteComment(grc.ctx, grc.owner, grc.repo, comment.GetID())
	return err
}

func (grc *GithubRepoClient) ReplyToComment(prId int, comment *github.PullRequestComment, body string) error {
	_, _, err := grc.client.PullRequests.CreateCommentInReplyTo(grc.ctx, grc.owner, grc.repo, prId, body, comment.GetID())
	return err
}
