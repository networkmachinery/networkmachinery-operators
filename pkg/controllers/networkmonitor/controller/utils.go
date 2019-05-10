package controller

import (
	"fmt"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
)

func filterEvents(events []v1alpha1.Event, monitor *v1alpha1.NetworkMonitor, test func(v1alpha1.Event, *v1alpha1.NetworkMonitor) bool) (ret []v1alpha1.Event) {
	for _, event := range events {
		if test(event, monitor) {
			ret = append(ret, event)
		}
	}
	return
}

func in(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getThresholdForName(thresholdName string, thresholds []v1alpha1.Threshold) (*v1alpha1.Threshold, error) {
	for _, threshold := range thresholds {
		if thresholdName == threshold.Name {
			return &threshold, nil
		}
	}
	return nil, fmt.Errorf("threshold with key %s does not exist", thresholdName)
}

func getFlowForName(flowName string, flows []v1alpha1.Flow) (*v1alpha1.Flow, error) {
	for _, flow := range flows {
		if flowName == flow.Name {
			return &flow, nil
		}
	}
	return nil, fmt.Errorf("flow with key %s does not exist", flowName)
}
