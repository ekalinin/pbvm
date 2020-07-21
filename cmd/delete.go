package cmd

import (
	"errors"

	"github.com/ekalinin/pbvm/utils"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Aliases: []string{"rm"},
	Use:     "delete <version>",
	Short:   "Delete version",
	Long:    `Delete version. Version should be installed.`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		installed, _, err := utils.IsInstalledVersion(pbName, version)
		if err != nil {
			return err
		}
		if !installed {
			// suppress help output
			// https://github.com/spf13/cobra/issues/340
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			return errors.New("Version " + version + " is not installed")
		}
		active, err := utils.IsActiveVersion(pbName, version)
		if err != nil {
			return err
		}
		if active {
			// suppress help output
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			return errors.New("Version " + version + " is active at the moment")
		}

		return utils.DeleteVersion(pbName, version)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
