package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	// Replace these with your GitHub token and repository details
	accessToken := os.Args[0]
	owner := "owner"
	repo := "repository"
	sourceBranch := "source-branch"
	targetBranch := "main"
	prTitle := "My Pull Request"
	prBody := "This is a pull request created using Go."

	// Set up OAuth2 authentication
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	// Create a new GitHub client
	client := github.NewClient(tc)

	// Create a new pull request
	newPR := &github.NewPullRequest{
		Title: github.String(prTitle),
		Head:  github.String(sourceBranch),
		Base:  github.String(targetBranch),
		Body:  github.String(prBody),
	}

	// Call the GitHub API to create the pull request
	pr, _, err := client.PullRequests.Create(ctx, owner, repo, newPR)
	if err != nil {
		log.Fatalf("Error creating pull request: %v", err)
	}

	fmt.Printf("Pull request created: %s\n", pr.GetHTMLURL())
}
