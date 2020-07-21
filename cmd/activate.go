package cmd

import (
	"fmt"

	"github.com/ekalinin/pbvm/utils"
	"github.com/spf13/cobra"
)

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate <version>",
	Short: "Activate version",
	Long:  `Activate version. Version should be installed.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		installed, _, err := utils.IsInstalledVersion(pbName, version)
		if err != nil {
			panic(err)
		}
		if !installed {
			fmt.Printf("Version %s is not installed.\n"+
				"Please, run: '%s install %[1]s'\n",
				version, pbName)
			return
		}
		if err := utils.ActivateVersion(pbName, version); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
}
