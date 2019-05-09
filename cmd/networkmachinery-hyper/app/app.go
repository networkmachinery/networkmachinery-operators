package app

import (
	"context"

	versioncmd "github.com/networkmachinery/networkmachinery-operators/version/cmd"

	networkmonitorcmd "github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkmonitor/cmd"
	"github.com/spf13/cobra"
)

// NewHyperCommand creates a new Hyper command consisting of all controllers under this repository.
func NewHyperCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "networkmachinery-hyper",
	}

	cmd.AddCommand(
		versioncmd.NewVersionCmd(),
		networkmonitorcmd.NewNetworkMonitorCmd(ctx),
	)

	return cmd
}
