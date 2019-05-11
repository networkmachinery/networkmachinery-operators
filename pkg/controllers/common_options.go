package controllers

type LeaderElectionOptions struct {
	LeaderElection          bool
	LeaderElectionNamespace string
	LeaderElectionID        string
}

type ControllerOptions struct {
	MaxConcurrentReconciles int
}
