package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"gopkg.in/src-d/go-git.v4"
)

// EnvSettings describes all of the environment settings.
type EnvSettings struct {
	//Full URL of gitlab host
	GitlabHost string
	//Personal access token
	GitlabToken string
	//Current repository
	Repo *git.Repository
}

//Init environment
func (s *EnvSettings) Init() {
	opt := &git.PlainOpenOptions{DetectDotGit: true}
	repo, err := git.PlainOpenWithOptions(".", opt)
	if err != nil {
		switch err {
		case git.ErrRepositoryNotExists:
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("Not a git repository."))
		default:
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("Error: %s", err))
		}
		os.Exit(1)
	}
	config, err := repo.Config()

	CheckIfErrorAndExit(err)
	gitlab := config.Raw.Section("gitlab")

	if len(strings.TrimSpace(gitlab.Option("url"))) == 0 {
		remote := config.Raw.Section("remote").Subsection("origin").Option("url")
		var def string
		u, err := url.Parse("//" + remote)
		if err == nil && len(remote) > 0 {
			def = "https://" + u.Hostname()
		}
		url := promptForString(&promptui.Prompt{Label: "Gitlab URL", Default: def})
		gitlab.SetOption("url", url)

		repo.Storer.SetConfig(config)
	}

	if len(strings.TrimSpace(gitlab.Option("token"))) == 0 {
		token := promptForString(&promptui.Prompt{Label: "Gitlab token"})
		gitlab.SetOption("token", token)

		repo.Storer.SetConfig(config)
	}
	s.GitlabHost = gitlab.Option("url")
	s.GitlabToken = gitlab.Option("token")
	s.Repo = repo
}

func promptForString(prompt *promptui.Prompt) string {
	result, err := prompt.Run()
	CheckIfErrorAndExit(err)
	return result
}
