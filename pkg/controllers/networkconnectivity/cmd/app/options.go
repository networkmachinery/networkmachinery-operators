package app

import (
	"time"

	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers"
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkconnectivity/controller"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// NetworkMonitorCmdOptions necessary options to run the sFlowController
// the current context on a user's KUBECONFIG
type NetworkConnectivityTestCmdOpts struct {
	disableWebhookConfigInstaller bool
	leaderLelectionRetryPeriod    *time.Duration
	*genericclioptions.ConfigFlags
	controllers.LeaderElectionOptions
}

func (nct *NetworkConnectivityTestCmdOpts) InjectLeaderElectionOpts(mgrOpts *manager.Options) *manager.Options {
	mgrOpts.LeaderElectionID = nct.LeaderElectionID
	mgrOpts.LeaderElectionNamespace = nct.LeaderElectionNamespace
	mgrOpts.LeaderElection = nct.LeaderElection
	return mgrOpts
}

func (nct *NetworkConnectivityTestCmdOpts) InjectRetryOptions(mgrOpts *manager.Options) *manager.Options {
	mgrOpts.RetryPeriod = nct.leaderLelectionRetryPeriod
	return mgrOpts
}


func (nct *NetworkConnectivityTestCmdOpts) InitConfig() *rest.Config {
	config, err := nct.ToRESTConfig()
	if err != nil {
		utils.LogErrAndExit(err, "Error getting config")
	}
	config.UserAgent = controller.Name
	return config
}

func (nct *NetworkConnectivityTestCmdOpts) AddWebHookFlags(flags *pflag.FlagSet) {
	flags.BoolVar(&nct.disableWebhookConfigInstaller, "disable-webhook-config-installer", false,
		"disable the installer in the webhook server, so it won't install webhook configuration resources during bootstrapping")
}

func (nct *NetworkConnectivityTestCmdOpts) AddAllFlags(flags *pflag.FlagSet) {
	nct.AddWebHookFlags(flags)
	nct.AddFlags(flags)
}
