package cmd

import (
	"context"
	"os"
	"strconv"

	"github.com/ekalinin/pbvm/utils"
	"github.com/google/go-github/v32/github"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var numberOfVersions int

// listRemoteCmd represents the listRemote command
var listRemoteCmd = &cobra.Command{
	Use:   "list-remote",
	Short: "List available versions",
	Long:  `List available versions`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		client := github.NewClient(nil)
		opts := &github.ListOptions{Page: 1, PerPage: numberOfVersions}
		releases, _, err := client.Repositories.ListReleases(ctx, pbOwner, pbRepo, opts)
		if err != nil {
			panic(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Version", "Pre-release", "Date", "Installed"})

		for _, r := range releases {
			installed, _, err := utils.IsInstalledVersion(pbName, *r.TagName)
			if err != nil {
				panic(err)
			}
			table.Append([]string{
				*r.TagName,
				strconv.FormatBool(*r.Prerelease),
				(*r.PublishedAt).Format("2006.01.02"),
				strconv.FormatBool(installed),
			})
		}
		table.SetBorder(false)
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(listRemoteCmd)

	// TODO: add options:
	// 	- hide pre-releases
	// 	- only pre-releases
	//	- show only latest
	//  - show descrition
	//  - show header
	// 	- show column "installed" (if version already installed)
	listRemoteCmd.Flags().IntVarP(&numberOfVersions, "number", "n", 10,
		"Number of last vertions to show")
}
