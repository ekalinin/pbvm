package cmd

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/ekalinin/pbvm/utils"
	"github.com/spf13/cobra"
)

var version string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:           "run <command>",
	Args:          cobra.ExactArgs(1),
	Short:         "Run a command under a version",
	Long:          `Run a command under a version`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(version) < 5 {
			return errors.New("Version is incorrect: " + version)
		}

		installed, _, err := utils.IsInstalledVersion(pbName, version)
		if err != nil {
			return err
		}
		if !installed {
			return errors.New("Version is not installed: " + version)
		}

		originVersion, err := utils.GetActiveVersion(pbName)
		if err != nil {
			return err
		}

		// TODO: run command without affecting other sessions
		// (change PATH & syscall.Exec?)

		if err := utils.ActivateVersion(pbName, version); err != nil {
			return err
		}
		defer utils.ActivateVersion(pbName, originVersion)

		cs := strings.Split(args[0], " ")
		command := exec.Command(cs[0], cs[1:]...)
		out, err := command.CombinedOutput()
		if err != nil {
			return err
		}
		println(string(out))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVar(&version, "version", "v",
		"Version used for command execution")
}
