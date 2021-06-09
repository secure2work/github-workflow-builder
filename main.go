package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

)

func main() {

	var token, repo string
	flag.StringVar(&token, "token", "", "token for github")
	flag.StringVar(&repo, "repo", "", "name of github repository")
	flag.Parse()

	fmt.Println("token is", token)
	fmt.Println("repo is", repo)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token}, // repos - gmlewis account
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	commitOption:= &github.RepositoryContentFileOptions{
		Branch:  github.String("main"),
		Message: github.String("testing this"),
		Committer: &github.CommitAuthor{
			Name:  github.String("bruteforce1414"),
			Email: github.String("bruteforce1414@gmail.com"),
		},
		Author: &github.CommitAuthor{
			Name:  github.String("bruteforce1414"),
			Email: github.String("bruteforce1414@gmail.com"),
		},
		Content: []byte("this is my content"),
	}


	client.Repositories.CreateFile(context.Background(), "secure2work", "github-workflow-builder", "test.txt", commitOption)

}
