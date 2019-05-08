package cmd

import (
	"context"
	controllercmd "github.com/gardener/gardener-extensions/pkg/controller/cmd"
	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/install"
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)



const Name = "sflow"

func NewSFlowControllerCmd(ctx context.Context){
	sFlowOpts := SFlowControllerOptions{
		LeaderElectionOptions: LeaderElectionOptions{
			LeaderElection:          true,
			LeaderElectionNamespace: "default",
			LeaderElectionID:        controllers.LeaderElectionNameID(Name),
		},
		ControllerOptions: ControllerOptions{
			MaxConcurrentReconciles: 5,
		},
	}

	cmd := &cobra.Command{
		Use: "sflow-controller-manager",
		Run: func(cmd *cobra.Command, args []string) {
			config, err := sFlowOpts.ToRESTConfig()
			if err != nil {
				controllercmd.LogErrAndExit(err, "Error getting config")
			}
			config.UserAgent = Name
			mgrOptions := &manager.Options{}
			mgr, err := manager.New(config, *sFlowOpts.ApplyLeaderElection(mgrOptions))
			if err != nil {
				controllercmd.LogErrAndExit(err, "Could not instantiate manager")
			}

			if err := install.AddToScheme(mgr.GetScheme()); err != nil {
				controllercmd.LogErrAndExit(err, "Could not update manager scheme")
			}

			if err := coreos.AddToManager(mgr); err != nil {
				controllercmd.LogErrAndExit(err, "Could not add controller to manager")
			}

			if err := mgr.Start(ctx.Done()); err != nil {
				controllercmd.LogErrAndExit(err, "Error running manager")
			}
		},
	}
	sFlowOpts.AddFlags(cmd.Flags())
}