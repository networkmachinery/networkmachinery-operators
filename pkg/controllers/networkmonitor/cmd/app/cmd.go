package app

import (
	"context"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers"

	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkmonitor/controller"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/install"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func NewNetworkMonitorCmd(ctx context.Context) *cobra.Command {
	networkMonitorCmdOpts := NetworkMonitorCmdOptions{
		ConfigFlags: genericclioptions.NewConfigFlags(),
		LeaderElectionOptions: controllers.LeaderElectionOptions{
			LeaderElection:          true,
			LeaderElectionNamespace: "default",
			LeaderElectionID:        utils.LeaderElectionNameID(controller.Name),
		},
		ControllerOptions: controllers.ControllerOptions{
			MaxConcurrentReconciles: 5,
		},
	}

	cmd := &cobra.Command{
		Use: "networkmonitor-controller",
		Run: func(cmd *cobra.Command, args []string) {
			mgrOptions := &manager.Options{}
			mgr, err := manager.New(networkMonitorCmdOpts.InitConfig(), *networkMonitorCmdOpts.ApplyLeaderElection(mgrOptions))
			if err != nil {
				utils.LogErrAndExit(err, "Could not instantiate manager")
			}
			if err := install.AddToScheme(mgr.GetScheme()); err != nil {
				utils.LogErrAndExit(err, "Could not update manager scheme")
			}

			if err := controllers.AddToManager(mgr); err != nil {
				utils.LogErrAndExit(err, "Could not add controller to manager")
			}

			if err := mgr.Start(ctx.Done()); err != nil {
				utils.LogErrAndExit(err, "Error running manager")
			}
		},
	}

	networkMonitorCmdOpts.AddFlags(cmd.Flags())
	return cmd
}
