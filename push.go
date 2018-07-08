package main

import (
	"fmt"
	"io"
	"github.com/spf13/cobra"
)

var longPushHelp = `
This command pushes the HEAD commit to a new remote branch 
and creates/updates a merge request.
`

func newPushCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "creates or updates a merge request",
		Long:  longPushHelp,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("push executed")
		},
	}
	return cmd
}
