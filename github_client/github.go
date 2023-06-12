package github_client

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/github"
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
