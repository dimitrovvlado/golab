package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

type mrCmd struct {
	gitlabURL   string
	gitlabToken string
	out         io.Writer
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
		Run: func(cmd *cobra.Command, args []string) {
			i.gitlabURL = settings.GitlabHost
			i.gitlabToken = settings.GitlabToken
			fmt.Printf("mr executed against %s\n", settings.GitlabHost)
		},
	}
	return cmd
}
