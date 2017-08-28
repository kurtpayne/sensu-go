package check

import (
	"fmt"

	"github.com/sensu/sensu-go/cli"
	"github.com/spf13/cobra"
)

// UpdateCommand adds command that allows user to create new checks
func UpdateCommand(cli *cli.SensuCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "update NAME",
		Short:        "update new checks",
		SilenceUsage: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Fetch handlers from API
			checkID := args[0]
			check, err := cli.Client.FetchCheck(checkID)
			if err != nil {
				return err
			}

			// Administer questionnaire
			opts := newCheckOpts()
			opts.withCheck(check)
			opts.administerQuestionnaire(true)

			// Apply given arguments to check
			opts.Copy(check)

			if err := check.Validate(); err != nil {
				return err
			}

			//
			// TODO:
			//
			// Current validation is a bit too laissez faire. For usability we should
			// determine whether there are assets / handlers / mutators associated w/
			// the check and warn the user if they do not exist yet.
			if err := cli.Client.CreateCheck(check); err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "OK")
			return nil
		},
	}

	cmd.Flags().StringP("command", "c", "", "the command the check should run")
	cmd.Flags().StringP("interval", "i", intervalDefault, "interval, in second, at which the check is run")
	cmd.Flags().StringP("subscriptions", "s", "", "comma separated list of topics check requests will be sent to")
	cmd.Flags().String("handlers", "", "comma separated list of handlers to invoke when check fails")
	cmd.Flags().StringP("runtime-assets", "r", "", "comma separated list of assets this check depends on")

	// Mark flags are required for bash-completions
	cmd.MarkFlagRequired("command")
	cmd.MarkFlagRequired("interval")
	cmd.MarkFlagRequired("subscriptions")

	return cmd
}
