package main

import (
	"context"
	"flag"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io"
	"log"
	"os"
	"text/template"
)

func main()  {

	var token, repo string
	flag.StringVar(&token, "token.txt", "", "token.txt for github")
	flag.StringVar(&repo, "repo", "", "name of github repository")
	flag.Parse()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	fileTemplate, err := os.Open("template_action.yml")
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		if err = fileTemplate.Close(); err != nil {
			os.Exit(1)
		}
	}()

	b, err := io.ReadAll(fileTemplate)
	if err != nil {
		os.Exit(1)
	}

	t:= template.Must(template.New("template_action").Delims("??", "??").Parse(string(b)))

	fileAction, err := os.Create("action.yml")
	if err != nil {
		os.Exit(1)
	}

	defer func() {
		if err = fileTemplate.Close(); err != nil {
			os.Exit(1)
		}
	}()


	pluginName:= struct {
		PluginName string
	}{repo}
	err = t.Execute(fileAction, pluginName)
	if err != nil {
		log.Println(err)
	}

	b, err = io.ReadAll(fileAction)
	if err != nil {
		os.Exit(1)
	}


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
		Content: b,
	}


	client.Repositories.CreateFile(context.Background(), "secure2work", "github-workflow-builder", ".github/workflows/action.yml", commitOption)

}
