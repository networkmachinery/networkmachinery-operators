package utils

import (
	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	"strings"
)

type Netcat struct {
	status v1alpha1.NetcatResultState
}

func (n *Netcat) State() v1alpha1.NetcatResultState {
	return n.status
}

func ParseNetcatOutput(outs []byte, nc *Netcat) {
	outputString := string(outs)
	switch {
	case strings.Contains(outputString, "refused"):
		nc.status = v1alpha1.Refused

	case strings.Contains(outputString, "succeeded"):
		nc.status = v1alpha1.Succeeded
	}
}
