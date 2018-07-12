package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var longInitHelp = `
This init command will set up any necessary local configuration.
`

type initCmd struct {
	token string
	out   io.Writer
}

func newInitCommand(out io.Writer) *cobra.Command {
	init := &initCmd{
		out: out,
	}

	cmd := &cobra.Command{
		Use:   "init",
		Short: "inits any necessary configuration",
		Long:  longInitHelp,
		RunE: func(cmd *cobra.Command, args []string) error {

			return init.run()
		},
	}

	f := cmd.Flags()

	f.StringVar(&init.token, "token", "", "Your token is required, you can find it on https://gitlab.eng.vmware.com/profile/account in the Private Tokens section.', type=str, help='Gitlab personal access token. Get it from https://gitlab.eng.vmware.com/profile/account")

	cmd.Flags()
	return cmd
}

func (i *initCmd) run() error {
	fmt.Println("init executed: " + i.token)
	return nil
}
