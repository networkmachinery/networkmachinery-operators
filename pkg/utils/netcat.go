package utils

import (
	"strings"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
)

type Netcat struct {
	status v1alpha1.NetcatResultState
}

func (n *Netcat) State() v1alpha1.NetcatResultState {
	return n.status
}

func ParseNetcatOutput(out string, nc *Netcat) {
	switch {
	case strings.Contains(out, "refused"), len(out) == 0:
		nc.status = v1alpha1.Refused
	case strings.Contains(out, "succeeded"), strings.Contains(out, "open"):
		nc.status = v1alpha1.Succeeded
	default:
		nc.status = v1alpha1.Refused
	}
}
