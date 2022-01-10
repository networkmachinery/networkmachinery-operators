package app

import (
	"context"

	networkconnectivitycmd "github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkconnectivity/cmd/app"
	networkdelaycmd "github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkdelay/cmd/app"

	networkcontrolcmd "github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkcontrol/cmd/app"
	networktrafficshapercmd "github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networktrafficshaper/cmd/app"

	networkmonitorcmd "github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkmonitor/cmd/app"
	versioncmd "github.com/networkmachinery/networkmachinery-operators/version/cmd"
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
		networkcontrolcmd.NewNetworkContrlCmd(ctx),
		networkconnectivitycmd.NewNetworkConnectivityTestCmd(ctx),
		networkdelaycmd.NewNetworkDelayTestCmd(ctx),
		networktrafficshapercmd.NewNetworkTrafficShaperCmd(ctx),
	)

	return cmd
}
