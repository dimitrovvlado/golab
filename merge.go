package main

import (
	"io"

	"github.com/spf13/cobra"
	git "gopkg.in/src-d/go-git.v4"
)

type mrCmd struct {
	gitlabURL    string
	gitlabToken  string
	baseBranch   string
	targetBranch string
	message      string
	assignee     string
	open         bool
	out          io.Writer
	repo         *git.Repository
}

var longMergeRequestHelp = `
This clean command will create a new merge request.
`

func newMergeRequestCommand(out io.Writer) *cobra.Command {
	i := &mrCmd{out: out}

	cmd := &cobra.Command{
		Use:     "mergerequest",
		Aliases: []string{"mr"},
		Short:   "Creates a new merge request",
		Long:    longMergeRequestHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			i.gitlabURL = settings.GitlabHost
			i.gitlabToken = settings.GitlabToken
			i.repo = settings.Repo

			return i.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&i.baseBranch, "base", "b", "", "Base branch name")
	f.StringVarP(&i.targetBranch, "target", "t", "", "Target branch name")
	f.StringVarP(&i.message, "message", "m", "", "Merge request message")
	f.BoolVarP(&i.open, "open", "o", false, "Opens edit page of merge request")

	return cmd
}

// run creates a merge request
func (i *mrCmd) run() error {
	if i.baseBranch == "" {
		ref, err := i.repo.Head()
		CheckIfErrorAndExit(err)
		i.baseBranch = ref.Name().Short()
	}
	return nil
}

// .option('-b, --base [optional]', 'Base branch name')
//   .option('-t, --target [optional]', 'Target branch name')
//   .option('-m, --message [optional]', 'Title of the merge request')
//   .option('-a, --assignee [optional]', 'User to assign merge request to')
//   .option('-l --labels [optional]', 'Comma separated list of labels to assign while creating merge request')
//   .option('-e, --edit [optional]', 'If supplied opens edit page of merge request. Opens merge request page otherwise')
//   .option('-p, --print [optional]', 'If supplied print the url of the merge request. Opens merge request page otherwise')
//   .option('-v, --verbose [optional]', 'Detailed logging emitted on console for debug purpose')
