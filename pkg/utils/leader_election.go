package utils

import "fmt"

// LeaderElectionNameID returns a leader election ID for the given name.
func LeaderElectionNameID(name string) string {
	return fmt.Sprintf("%s-leader-election", name)
}
