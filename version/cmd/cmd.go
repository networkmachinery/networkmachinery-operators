package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of networkmachinery-sflow",
	Long:  `All software has versions. This is networkmachinery-sflow`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Build Date:", BuildDate)
		fmt.Println("Git Commit:", GitCommit)
		fmt.Println("Version:", Version)
		fmt.Println("Go Version:", GoVersion)
		fmt.Println("OS / Arch:", OsArch)
	},
}
