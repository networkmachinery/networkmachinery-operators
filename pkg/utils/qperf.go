package utils

import (
	"regexp"
	"time"

	networkmachineryv1alpha1 "github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
)

type Qperf struct {
	tcpdelay, udpdelay time.Duration
}

func (q *Qperf) TCPDelay() string {
	return q.tcpdelay.String()
}

func (q *Qperf) UDPDelay() string {
	return q.udpdelay.String()
}

func ParseQperfOutput(outs []byte, qperf *Qperf, protocol networkmachineryv1alpha1.Protocol) {
	var latency = regexp.MustCompile(`latency\s+=\s+(\d+).(\d+)\sus`)
	result := latency.FindAllStringSubmatch(string(outs), -1)
	if len(result) > 0 {
		switch protocol {
		case networkmachineryv1alpha1.TCP:
			qperf.tcpdelay, _ = time.ParseDuration(result[0][1] + "us")
		case networkmachineryv1alpha1.UDP:
			qperf.udpdelay, _ = time.ParseDuration(result[0][2] + "us")
		}
	}
}
