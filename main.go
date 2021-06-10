package main

import (
	"bytes"
	"context"
	"flag"
	"io"
	"log"
	"os"
	"text/template"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	var token, owner, repo string
	flag.StringVar(&owner, "owner", "", "owner for github")
	flag.StringVar(&token, "token", "", "token for github")
	flag.StringVar(&repo, "repo", "", "name of github repository")
	flag.Parse()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	b, err := actionCreating(repo)
	if err != nil {
		log.Println("action creating error", err.Error())
		os.Exit(1)
	}

	commitOption, err := getCommitOptions(ctx, client, b, owner, repo)
	if err != nil {
		log.Println("action creating error", err.Error())
		os.Exit(1)
	}

	if commitOption.SHA == nil {
		log.Println("CREATING")
		_, _, err = client.Repositories.CreateFile(ctx, owner, repo, ".github/workflows/action.yml", commitOption)
		if err != nil {
			log.Println(err.Error())
		}
		os.Exit(0)
	}
	_, _, err = client.Repositories.UpdateFile(ctx, owner, repo, ".github/workflows/action.yml", commitOption)
	log.Println("UPDATING")
	if err != nil {
		log.Println(err.Error())
	}

}

func actionCreating(repo string) (*bytes.Buffer, error) {
	fileTemplate, err := os.Open("templates/action.yml")
	if err != nil {
		log.Println("action.yml open error", err)
		return nil, err
	}
	defer func() {
		if err = fileTemplate.Close(); err != nil {
			log.Println("templates/action.yml", err)
			os.Exit(1)
		}
	}()

	var b []byte
	b, err = io.ReadAll(fileTemplate)
	if err != nil {
		log.Println("action.yml readAll error", err)
		return nil, err
	}
	t := template.Must(template.New("action").Delims("??", "??").Parse(string(b)))
	buf := bytes.Buffer{}
	pluginName := struct {
		PluginName string
	}{repo}
	err = t.Execute(&buf, pluginName)
	if err != nil {
		log.Println("t.Execute", err)
		return nil, err
	}

	return &buf, nil
}

func getCommitOptions(ctx context.Context, client *github.Client, b *bytes.Buffer, owner, repo string) (*github.RepositoryContentFileOptions, error) {
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, ".github/workflows/action.yml", nil)
	if err != nil {
		switch e := err.(type) {
		case *github.ErrorResponse:
			if e.Response.StatusCode == 404 {
				log.Println("Not found")
				return nil, err
			}
		default:
			log.Println(err.Error())
			return nil, err
		}
	}
	var sha *string
	if fileContent != nil {
		sha = fileContent.SHA
	}
	commitOption := &github.RepositoryContentFileOptions{
		Branch:  github.String("main"),
		Message: github.String("testing this 6"),
		Committer: &github.CommitAuthor{
			Name:  github.String("bruteforce1414"),
			Email: github.String("bruteforce1414"),
		},
		Author: &github.CommitAuthor{
			Name:  github.String("bruteforce1414"),
			Email: github.String("bruteforce1414"),
		},
		Content: b.Bytes(),
		SHA:     sha,
	}
	return commitOption, nil
}
