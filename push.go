package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
)

type pushCmd struct {
	gitlabURL   string
	gitlabToken string
	out         io.Writer
	repo        *git.Repository
}

var longPushHelp = `
This command pushes the HEAD commit to a new remote branch 
and creates/updates a merge request.
`

func newPushCommand(out io.Writer) *cobra.Command {
	i := &pushCmd{out: out}

	cmd := &cobra.Command{
		Use:   "push",
		Short: "creates or updates a merge request",
		Long:  longPushHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			i.gitlabURL = settings.GitlabHost
			i.gitlabToken = settings.GitlabToken
			i.repo = settings.Repo

			return i.run()
		},
	}
	return cmd
}

// run creates a merge request
func (i *pushCmd) run() error {
	ref, err := i.repo.Head()
	CheckIfErrorAndExit(err)
	commit, err := i.repo.CommitObject(ref.Hash())
	CheckIfErrorAndExit(err)

	var cid string
	for _, line := range strings.Split(strings.TrimSuffix(commit.Message, "\n"), "\n") {
		i := strings.Index(line, "Change-Id")
		if i > -1 {
			cid = strings.TrimSpace(line[i+10:])
		}
	}
	if cid == "" {
		return errors.New("Commit message contains no Change-Id")
	}
	fmt.Printf("\x1b[32m%s\x1b[0m\n", fmt.Sprintf("Creating merge request with internal Change-Id: %s", cid))

	spec := config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/remotes/origin/%s", cid, cid))
	refSpecs := []config.RefSpec{spec}
	err = i.repo.Push(&git.PushOptions{RefSpecs: refSpecs})

	if err != nil {
		switch err {
		case git.NoErrAlreadyUpToDate:
			fmt.Printf("\x1b[32m%s\x1b[0m\n", fmt.Sprintf(err.Error()))
		default:
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("Error: %s", err))
			os.Exit(1)
		}
	}

	return nil
}
