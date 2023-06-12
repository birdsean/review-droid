# Review Droid

This is a simple app that will utilize AI to do an initial pass on your code reviews. It will look for common issues and suggest fixes. It will also look for common patterns and suggest refactors.

## Features

### MVP
- [ ] GitHub OAuth: Implement the OAuth flow to allow users to authenticate with GitHub and authorize the plugin to access their repositories.
- [ ] Webhook Listener: Implement the functionality to listen for pull request events from GitHub.
- [ ] Code Retrieval: Implement the functionality to fetch the code changes from the pull request.
- [ ] Language Model Integration: Integrate with OpenAI API.
- [ ] Code Analysis: Feed the code changes retrieved from the pull request into the LLM. Extract comments or suggestions related to code quality, potential bugs, or best practices from the LLM's output.
[] Comment Posting: Develop the functionality to post the extracted comments and suggestions as comments directly on the pull request within the GitHub interface.

## Entities
* TBD