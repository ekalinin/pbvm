package cmd

import (
	"os"
	"strconv"

	"github.com/ekalinin/pbvm/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listLocalCmd represents the listLocal command
var listLocalCmd = &cobra.Command{
	Aliases: []string{"ls"},
	Use:     "list-local",
	Short:   "List local (previously installed) versions",
	Long:    `Shows list of installed versions.`,
	Run: func(cmd *cobra.Command, args []string) {
		// version, install date (folder stat), active?
		versions, err := utils.ListInstalledVersions(pbName)
		if err != nil {
			panic(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Version", "Install date", "Active"})

		for _, v := range versions {
			table.Append([]string{
				v.Version,
				v.Date.Format(pbDateFormat),
				strconv.FormatBool(v.Active),
			})
		}
		table.SetBorder(false)
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listLocalCmd)
}
