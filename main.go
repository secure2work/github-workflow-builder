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
	//@todo owner as a flag?
	ctx:=context.Background()
	var token, owner, repo string
	flag.StringVar(&owner, "owner", "", "owner for github")
	flag.StringVar(&token, "token", "", "token for github")
	flag.StringVar(&repo, "repo", "", "name of github repository")
	flag.Parse()

	var b []byte
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	fileTemplate, err := os.Open("template_action.yml")
	if err != nil {
		log.Println("template_action.yml open error", err)
		os.Exit(1)
	}
	defer func() {
		if err = fileTemplate.Close(); err != nil {
			log.Println("template_action.yml close error", err)
			os.Exit(1)
		}
	}()

	b, err = io.ReadAll(fileTemplate)
	if err != nil {
		log.Println("template_action.yml readAll error", err)
		os.Exit(1)
	}

	t:= template.Must(template.New("template_action").Delims("??", "??").Parse(string(b)))

	fileAction, err := os.Create("action.yml")
	if err != nil {
		log.Println("action.yml readAll error", err)
		os.Exit(1)
	}




	pluginName:= struct {
		PluginName string
	}{repo}
	err = t.Execute(fileAction, pluginName)
	if err != nil {
		log.Println("t.Execute", err)
		os.Exit(1)

	}
	if err = fileAction.Close(); err != nil {
		log.Println("action.yml close error", err)

		os.Exit(1)
	}
	fileAction2, err := os.Open("action.yml")
	if err != nil {
		log.Println("action.yml readAll error", err)
		os.Exit(1)
	}
	b, err = io.ReadAll(fileAction2)
	if err != nil {
		log.Println("action.yml readall", err)
		os.Exit(1)
	}
	log.Println("b", b)
	if err = fileAction2.Close(); err != nil {
		log.Println("action.yml close error", err)

		os.Exit(1)
	}


/*	repositories, _,err:=client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		os.Exit(1)
	}*/

	//log.Println("files", repositories.GetBlobsURL())
	commitOption:= &github.RepositoryContentFileOptions{
		Branch:  github.String("main"),
		Message: github.String("testing this1"),
		Committer: &github.CommitAuthor{
			Name:  github.String("bruteforce1414"),
			Email: github.String("bruteforce1414@gmail.com"),
		},
		Author: &github.CommitAuthor{
			Name:  github.String("bruteforce1414"),
			Email: github.String("bruteforce1414@gmail.com"),
		},
		Content: b,
		//SHA:
	}


	client.Repositories.CreateFile(ctx, owner, repo, ".github/workflows/action.yml", commitOption)

}
