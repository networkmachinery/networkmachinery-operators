package app

import (
	"context"
	"time"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkcontrol/controller"

	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/install"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var log = logf.Log.WithName("networkdelay-test-controller")

func NewNetworkDelayTestCmd(ctx context.Context) *cobra.Command {
	var (
		retryDuration           = 100 * time.Millisecond
		networkDelayTestCmdOpts = NetworkDelayTestCmdOpts{
			ConfigFlags: genericclioptions.NewConfigFlags(true),
			LeaderElectionOptions: controllers.LeaderElectionOptions{
				LeaderElection:          true,
				LeaderElectionNamespace: "default",
				LeaderElectionID:        utils.LeaderElectionNameID(controller.Name),
			},
			leaderLelectionRetryPeriod: &retryDuration,
		}
	)

	cmd := &cobra.Command{
		Use: "networkdelay-test-controller",
		Run: func(cmd *cobra.Command, args []string) {
			mgrOptions := &manager.Options{}
			mgr, err := manager.New(networkDelayTestCmdOpts.InitConfig(),
				*networkDelayTestCmdOpts.InjectRetryOptions(networkDelayTestCmdOpts.InjectLeaderElectionOpts(mgrOptions)))
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
	networkDelayTestCmdOpts.AddAllFlags(cmd.Flags())
	return cmd
}
