package networkmonitor

import (
	"context"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
)

type NetworkMonitorer interface {
	InstallFlows(context.Context, *v1alpha1.NetworkMonitor) error
	InstallThresholds(context.Context, *v1alpha1.NetworkMonitor) error
	DeleteFlows(context.Context, *v1alpha1.NetworkMonitor) error
	DeleteThresholds(context.Context, *v1alpha1.NetworkMonitor) error
	CheckEvents(context.Context, *v1alpha1.NetworkMonitor) ([]v1alpha1.Event, error)
}
