package app

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/webhook"

	networkmachineryhandlers "github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkconnectivity/webhook"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers"

	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkconnectivity/controller"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/install"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var log = logf.Log.WithName("example-controller")

const (
	layerValidationServerPath       = "/validate-layer-v1alpha1-networkconnectivitytest"
	destinationValidationServerPath = "/validate-destination-v1alpha1-networkconnectivitytest"

	webhookServerPort = 9876
)

func NewNetworkConnectivityTestCmd(ctx context.Context) *cobra.Command {
	entryLog := log.WithName("networkconnectivity-test-cmd")

	networkConnectivityTestCmdOpts := NetworkConnectivityTestCmdOpts{
		ConfigFlags: genericclioptions.NewConfigFlags(true),
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
		Use: "networkconnectivity-test-controller",
		Run: func(cmd *cobra.Command, args []string) {
			mgrOptions := &manager.Options{
				Port: webhookServerPort,
			}
			mgr, err := manager.New(networkConnectivityTestCmdOpts.InitConfig(), *networkConnectivityTestCmdOpts.ApplyLeaderElection(mgrOptions))
			if err != nil {
				utils.LogErrAndExit(err, "Could not instantiate manager")
			}
			if err := install.AddToScheme(mgr.GetScheme()); err != nil {
				utils.LogErrAndExit(err, "Could not update	 manager scheme")
			}

			entryLog.Info("setting up webhook server")
			admissionServer := mgr.GetWebhookServer()

			entryLog.Info("registering webhooks to the webhook server")
			admissionServer.Register(layerValidationServerPath, &webhook.Admission{Handler: &networkmachineryhandlers.LayerValidator{}})
			admissionServer.Register(destinationValidationServerPath, &webhook.Admission{Handler: &networkmachineryhandlers.DestinationValidator{}})

			if err := controllers.AddToManager(mgr); err != nil {
				utils.LogErrAndExit(err, "Could not add controller to manager")
			}

			if err := mgr.Start(ctx.Done()); err != nil {
				utils.LogErrAndExit(err, "Error running manager")
			}
		},
	}
	networkConnectivityTestCmdOpts.AddAllFlags(cmd.Flags())
	return cmd
}
