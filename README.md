# Review Droid

This is a simple app that will utilize AI to do an initial pass on your code reviews. It will look for common issues and suggest fixes. It will also look for common patterns and suggest refactors.

## Features

Code Retrieval: Will pull down code updates from all open pull requests in a repository.
Language Model Integration: Integrated with OpenAI's API
Code Analysis: Feeds the code changes retrieved from the pull request into the LLM. Extracts comments or suggestions related to code quality, potential bugs, or best practices from the LLM's output.
Comment Posting: Posts comments directly on the pull request within the GitHub interface.

### Backlog
- [ ] GitHub App OAuth: Implement the OAuth flow to allow users to authenticate with GitHub and authorize the plugin to access their repositories. Implement as a GitHub App.
- [ ] Webhook Listener: Implement the functionality to listen for pull request events from GitHub.

## Entities
* TBD