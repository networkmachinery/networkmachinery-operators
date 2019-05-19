package app

import (
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers"
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkmonitor/controller"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// NetworkTrafficShaperCmdOptions options for the network traffic shaper cmd
type NetworkTrafficShaperCmdOptions struct {
	*genericclioptions.ConfigFlags
	controllers.LeaderElectionOptions
	controllers.ControllerOptions
}

func (nm *NetworkTrafficShaperCmdOptions) ApplyLeaderElection(mgr *manager.Options) *manager.Options {
	mgr.LeaderElectionID = nm.LeaderElectionID
	mgr.LeaderElectionNamespace = nm.LeaderElectionNamespace
	mgr.LeaderElection = nm.LeaderElection
	return mgr
}

func (nm *NetworkTrafficShaperCmdOptions) InitConfig() *rest.Config {
	config, err := nm.ToRESTConfig()
	if err != nil {
		utils.LogErrAndExit(err, "Error getting config")
	}
	config.UserAgent = controller.Name
	return config
}
