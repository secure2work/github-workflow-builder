package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
)

func main() {

	var token, repo string
	flag.StringVar(&token, "token", "", "token for github")
	flag.StringVar(&repo, "repo", "", "name of github repository")
	flag.Parse()

	fmt.Println("token is", token)
	fmt.Println("repo is", repo)

	client := github.NewClient(nil)
	branch:="main"
	rcfo:=github.RepositoryContentFileOptions{
		Message:   nil,
		Content:   nil,
		SHA:       nil,
		Branch:    &branch,
		Author:    nil,
		Committer: nil,
	}
	client.Repositories.CreateFile(context.Background(), "secure2work", "github-workflow-builder", "test.txt", &rcfo)

}
