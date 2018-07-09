package main

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	settings EnvSettings
)

var globalUsage = `The Usage Meter Services command for Gitlab
To begin working with UMS, run the 'ums init' command:
	$ ums init --token
This will set up any necessary local configuration.
Common actions from this point include:
- ums push:    push HEAD to a 
- ums clean:   clean the configuration
Environment:
  $GITLAB_TOKEN          set an alternative Gitlab token
  $GITLAB_URL            set an alternative Gitlab URL
`

func main() {
	cmd := newRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

}

func newRootCmd(args []string) *cobra.Command {

	cmd := &cobra.Command{
		Use:          "ums",
		Short:        "The Usage Meter Services command for merge requests.",
		Long:         globalUsage,
		SilenceUsage: true,
	}
	flags := cmd.PersistentFlags()
	flags.Parse(args)

	out := cmd.OutOrStdout()

	cmd.AddCommand(
		newPushCommand(out),
		newInitCommand(out),
		newCleanCommand(out),
		newMergeRequestCommand(out),
	)

	settings.Init()

	return cmd
}

// CheckIfErrorAndExit should be used to naively panics if an error is not nil.
func CheckIfErrorAndExit(err error) {
	if err != nil {
		switch err {
		case promptui.ErrInterrupt:
			//Exit gracefully
			os.Exit(0)
		default:
			fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("Error: %s", err))
			os.Exit(1)
		}
	}
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
