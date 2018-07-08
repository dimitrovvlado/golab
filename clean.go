package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var longCleanHelp = `
This clean command will delete all local configuration.
`

func newCleanCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "cleans all local configuration",
		Long:  longCleanHelp,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("clean executed")
		},
	}
	return cmd
}
