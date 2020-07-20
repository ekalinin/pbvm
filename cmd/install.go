package cmd

import (
	"context"
	"runtime"

	"github.com/ekalinin/pbvm/utils"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install <version>",
	Short: "Install a version",
	Long: `Install a version.

If entered version was installled before, then that version will be
just enabled. In another case, entered version will be downloaded,
installed and enabled.

To get all available versions use "list-remote" command.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		client := github.NewClient(nil)
		tag := args[0]

		d("Searching release: ", tag, " ...")
		release, _, err := client.Repositories.GetReleaseByTag(ctx, pbOwner, pbRepo, tag)
		if err != nil {
			panic(err)
		}
		d(" ... found: ", release.URL)

		d("Searching asset: ...")
		arch := utils.GetArch()
		var asset *github.ReleaseAsset
		for _, a := range release.Assets {
			if !utils.IsSuitableAsset(*a.Name, arch, runtime.GOOS) {
				d(" ... skip:", *a.Name)
				continue
			}
			d(" ... found:", *a.BrowserDownloadURL)
			asset = a
			break
		}

		if asset == nil {
			panic("Could not found asset.")
		}

		d("Downloading version: ", tag, " ...")
		downloaded, err := utils.DownloadVersion(pbName, tag, asset)
		if err != nil {
			panic(err)
		}
		d(" ... is realy downloaded: ", downloaded)

		d("Activating version: ", tag, " ...")
		if err := utils.ActivateVersion(pbName, tag); err != nil {
			panic(err)
		}

		d("Done.")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
