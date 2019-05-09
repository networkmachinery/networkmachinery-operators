package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"gopkg.in/resty.v1"

	"github.com/networkmachinery/networkmachinery-operators/pkg/utils/apimachinery"

	"github.com/go-logr/logr"
	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// FinalizerName is the controlplane controller finalizer.
	FinalizerName = "networkmachinery.io/networkmonitor"
)

// ReconcileMachineDeployment reconciles a MachineDeployment object.
type ReconcileNetworkMonitor struct {
	logger   logr.Logger
	client   client.Client
	ctx      context.Context
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

func (r *ReconcileNetworkMonitor) InjectClient(client client.Client) error {
	r.client = client
	return nil
}

func (r *ReconcileNetworkMonitor) InjectStopChannel(stopCh <-chan struct{}) error {
	r.ctx = utils.ContextFromStopChannel(stopCh)
	return nil
}

func (r *ReconcileNetworkMonitor) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	networkMonitor := &v1alpha1.NetworkMonitor{}
	if err := r.client.Get(r.ctx, request.NamespacedName, networkMonitor); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if networkMonitor.DeletionTimestamp != nil {
		return r.delete(r.ctx, networkMonitor)
	}

	r.logger.Info("Starting the reconciliation of NetworkMonitor", "Name", networkMonitor.Name)
	r.recorder.Event(networkMonitor, v1alpha1.EventTypeNormal, v1alpha1.EventTypeReconciliation, "Reconciling NetworkMonitor")

	return r.reconcile(r.ctx, networkMonitor)
}

func (r *ReconcileNetworkMonitor) delete(ctx context.Context, networkmonitor *v1alpha1.NetworkMonitor) (reconcile.Result, error) {
	hasFinalizer, err := apimachinery.HasFinalizer(networkmonitor, FinalizerName)
	if err != nil {
		r.logger.Error(err, "Could not instantiate finalizer deletion")
		return reconcile.Result{}, err
	}
	if !hasFinalizer {
		r.logger.Info("Deleting NetworkMonitor causes a no-op as there is no finalizer.", "NetworkMonitor", networkmonitor.Name)
		return reconcile.Result{}, nil
	}

	r.logger.Info("Starting the deletion of network monitor", "NetworkMonitor", networkmonitor.Name)
	r.recorder.Event(networkmonitor, v1alpha1.EventTypeNormal, v1alpha1.EventTypeDeletion, "Deleting the Network Monitor")

	if err := apimachinery.DeleteFinalizer(ctx, r.client, FinalizerName, networkmonitor); err != nil {
		r.logger.Error(err, "Error removing finalizer from the NetworkMonitor resource", "Network Monitor", networkmonitor.Name)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileNetworkMonitor) reconcile(ctx context.Context, networkMonitor *v1alpha1.NetworkMonitor) (reconcile.Result, error) {
	r.logger.Info("Adding finalizer", "NetworkMonitor", networkMonitor.Name)
	if err := apimachinery.EnsureFinalizer(ctx, r.client, FinalizerName, networkMonitor); err != nil {
		return reconcile.Result{}, err
	}

	err := installFlows(networkMonitor)
	if err != nil {
		return apimachinery.ReconcileErr(err)
	}
	r.logger.Info("Installed flows successfully", "NetworkMonitor", networkMonitor.Name)
	r.recorder.Event(networkMonitor, v1alpha1.EventTypeNormal, v1alpha1.EventTypeReconciliation, "Installed Flows")

	err = installThresholds(networkMonitor)
	if err != nil {
		return apimachinery.ReconcileErr(err)
	}
	r.logger.Info("Installed thresholds successfully", "NetworkMonitor", networkMonitor.Name)
	r.recorder.Event(networkMonitor, v1alpha1.EventTypeNormal, v1alpha1.EventTypeReconciliation, "Installed Thresholds")

	return reconcile.Result{}, nil
}

func getFlowForName(flowName string, flows []v1alpha1.Flow) (*v1alpha1.Flow, error) {
	for _, flow := range flows {
		if flowName == flow.Name {
			return &flow, nil
		}
	}
	return nil, fmt.Errorf("flow with key %s does not exist", flowName)
}

func hasParentFlow(threshold v1alpha1.Threshold) bool {
	return false
}

func installThresholds(networkMonitor *v1alpha1.NetworkMonitor) error {
	installThreshold := func(threshold v1alpha1.Threshold, nm *v1alpha1.NetworkMonitor) error {
		jsonFlow, err := json.Marshal(threshold)
		if err != nil {
			return err
		}
		url := fmt.Sprintf("http://%s/threshold/%s/json", net.JoinHostPort(nm.Spec.MonitoringEndpoint.IP, nm.Spec.MonitoringEndpoint.Port), threshold.Name)
		fmt.Println(url)
		resp, err := resty.R().
			SetHeader("Content-Type", "application/json").
			SetBody(jsonFlow).
			Put(url)

		fmt.Println(resp.RawResponse)
		return err
	}

	for _, threshold := range networkMonitor.Spec.Thresholds {
		if err := installThreshold(threshold, networkMonitor); err != nil {
			return err
		}
	}

	return nil
}

func installFlows(networkMonitor *v1alpha1.NetworkMonitor) error {
	installFlow := func(flow v1alpha1.Flow, nm *v1alpha1.NetworkMonitor) error {
		jsonFlow, err := json.Marshal(flow)
		if err != nil {
			return err
		}
		url := fmt.Sprintf("http://%s/flow/%s/json", net.JoinHostPort(nm.Spec.MonitoringEndpoint.IP, nm.Spec.MonitoringEndpoint.Port), flow.Name)
		fmt.Println(url)
		resp, err := resty.R().
			SetHeader("Content-Type", "application/json").
			SetBody(jsonFlow).
			Put(url)

		fmt.Println(resp.RawResponse)
		return err
	}

	for _, flow := range networkMonitor.Spec.Flows {
		if err := installFlow(flow, networkMonitor); err != nil {
			return err
		}
	}

	return nil
}
