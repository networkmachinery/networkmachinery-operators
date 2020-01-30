package app

import (
	"time"

	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers"
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkdelay/controller"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type NetworkDelayTestCmdOpts struct {
	leaderLelectionRetryPeriod *time.Duration
	*genericclioptions.ConfigFlags
	controllers.LeaderElectionOptions
}

func (ndt *NetworkDelayTestCmdOpts) InjectLeaderElectionOpts(mgrOpts *manager.Options) *manager.Options {
	mgrOpts.LeaderElectionID = ndt.LeaderElectionID
	mgrOpts.LeaderElectionNamespace = ndt.LeaderElectionNamespace
	mgrOpts.LeaderElection = ndt.LeaderElection
	return mgrOpts
}

func (ndt *NetworkDelayTestCmdOpts) InjectRetryOptions(mgrOpts *manager.Options) *manager.Options {
	mgrOpts.RetryPeriod = ndt.leaderLelectionRetryPeriod
	return mgrOpts
}

func (ndt *NetworkDelayTestCmdOpts) InitConfig() *rest.Config {
	config, err := ndt.ToRESTConfig()
	if err != nil {
		utils.LogErrAndExit(err, "Error getting config")
	}
	config.UserAgent = controller.Name
	return config
}

func (ndt *NetworkDelayTestCmdOpts) AddAllFlags(flags *pflag.FlagSet) {
	ndt.AddFlags(flags)
}
