package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/google/go-github/github"
	sshpkg "golang.org/x/crypto/ssh"
	"golang.org/x/oauth2"
)

func main() {
	// Replace these with your GitHub token and repository details
	accessToken := os.Args[1]
	owner := "dengsh12"
	repoName := "TrialRep"
	sourceBranch := "mybranch1"
	targetBranch := "main"
	prTitle := "My Pull Request"
	prBody := "This is a pull request created using Go."
	localPath := "." // Path to your cloned repository

	// Load SSH private key
	sshKeyPath := os.Getenv("HOME") + "/.ssh/id_ed25519"
	sshKey, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		log.Fatalf("Error reading SSH key: %v", err)
	}

	signer, err := sshpkg.ParsePrivateKey(sshKey)
	if err != nil {
		log.Fatalf("Error parsing SSH key: %v", err)
	}

	auth := &ssh.PublicKeys{
		User:   "git",
		Signer: signer,
	}

	// Open the existing repository
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		log.Fatalf("Error opening repository: %v", err)
	}

	// Pull the latest changes from the target branch
	w, err := repo.Worktree()
	if err != nil {
		log.Fatalf("Error getting worktree: %v", err)
	}

	err = w.AddWithOptions(&git.AddOptions{
		All: true,
	})
	if err != nil {
		log.Fatalf("Error adding changes: %v", err)
	}

	// Commit the changes
	commitMsg := "Committing all changes before pull"
	_, err = w.Commit(commitMsg, &git.CommitOptions{})
	if err != nil {
		log.Fatalf("Error committing changes: %v", err)
	}

	err = w.Pull(&git.PullOptions{
		RemoteName:    "origin",
		Auth:          auth,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + targetBranch),
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		log.Fatalf("Error pulling latest changes: %v", err)
	}

	// Create a new branch
	newBranch := plumbing.NewBranchReferenceName(sourceBranch)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: newBranch,
		Create: true,
	})
	if err != nil {
		log.Fatalf("Error creating new branch: %v", err)
	}

	// Add your changes (if any) and commit them
	// Example: w.Add(".") and w.Commit("Your commit message", &git.CommitOptions{})

	// Push the new branch to the remote repository
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	if err != nil {
		log.Fatalf("Error pushing new branch: %v", err)
	}

	// Create a new pull request
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	newPR := &github.NewPullRequest{
		Title: github.String(prTitle),
		Head:  github.String(sourceBranch),
		Base:  github.String(targetBranch),
		Body:  github.String(prBody),
	}

	pr, _, err := client.PullRequests.Create(ctx, owner, repoName, newPR)
	if err != nil {
		log.Fatalf("Error creating pull request: %v", err)
	}

	fmt.Printf("Pull request created!: %s\n", pr.GetHTMLURL())

}
