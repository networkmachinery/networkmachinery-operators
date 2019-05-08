package cmd

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)


var managerOptions = manager.Options{}

type LeaderElectionOptions struct {
	LeaderElection bool
	LeaderElectionNamespace string
	LeaderElectionID string
}

type ControllerOptions struct {
	MaxConcurrentReconciles int
}
// SFlowControllerOptions necessary options to run the sFlowController
// the current context on a user's KUBECONFIG
type SFlowControllerOptions struct {
	*genericclioptions.ConfigFlags
	LeaderElectionOptions
	ControllerOptions
}

func (l *LeaderElectionOptions) ApplyLeaderElection(mgr *manager.Options) *manager.Options{
	mgr.LeaderElectionID = l.LeaderElectionID
	mgr.LeaderElectionNamespace= l.LeaderElectionNamespace
	mgr.LeaderElection = l.LeaderElection
	return mgr
}


