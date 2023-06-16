# Review Droid

This is a simple app that will utilize AI to do an initial pass on your code reviews. It will look for common issues and suggest fixes. It will also look for common patterns and suggest refactors.

## Features

* `Retrieves PRs`: Will pull down code updates from all open pull requests in a repository.
* `Auto-Reviews`: Feeds the code changes retrieved from the pull request into the OpenAI API. Extracts comments or suggestions related to code quality, potential bugs, or best practices from the LLM's output. 
* `Posts Review Comments`: Posts comments directly on the pull request within the GitHub interface.
* `Self Evaluates`: When the comments are generated, they are instantly ran back through the LLM to determine if they are good comments or not. Bad ones are discarded. After all comments are posted, each is ran through the LLM again in context of the code the comments they describe and bad ones are discareded.

## Get Started

### Requirements
* Go 1.20

### Installation
```bash
git clone git@github.com:birdsean/review-droid.git
cd review-droid
go mod install
```

### Usage
```bash
go run main.go
```

> Note: you will need to set the following environment variables:
> * REVIEW_DROID_TOKEN (GitHub Personal Access Token)
> * GITHUB_OWNER (Owner of the repository you'd like to review)
> * GITHUB_REPO (Name of the repository you'd like to review)
> * OPENAI_TOKEN (OpenAI API Secret Token)
