package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

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
	r.logger.Info("Ensuring finalizer", "NetworkMonitor", networkMonitor.Name)
	if err := apimachinery.EnsureFinalizer(ctx, r.client, FinalizerName, networkMonitor); err != nil {
		return reconcile.Result{}, err
	}

	err := r.installFlows(networkMonitor)
	if err != nil {
		return apimachinery.ReconcileErr(err)
	}

	err = r.installThresholds(networkMonitor)
	if err != nil {
		return apimachinery.ReconcileErr(err)
	}

	_, err = r.checkEvents(networkMonitor)
	if err != nil {
		return apimachinery.ReconcileErr(err)
	}
	// this is to allow for polling fresh events out of sFlow
	return reconcile.Result{
		RequeueAfter: 60 * time.Second,
	}, nil
}

func (r *ReconcileNetworkMonitor) installFlows(networkMonitor *v1alpha1.NetworkMonitor) error {
	var sflowAddress = net.JoinHostPort(networkMonitor.Spec.MonitoringEndpoint.IP, networkMonitor.Spec.MonitoringEndpoint.Port)

	flowExists := func(flowName string) (bool, error) {
		resp, err := resty.R().Get("http://" + sflowAddress + fmt.Sprintf("/flow/%s/json", flowName))
		if err != nil {
			return false, err
		}
		// TODO: check other 20X codes
		if resp.StatusCode() != http.StatusOK {
			return false, nil
		}
		return true, nil
	}

	installFlow := func(flow v1alpha1.Flow, nm *v1alpha1.NetworkMonitor) error {
		jsonFlow, err := json.Marshal(flow)
		if err != nil {
			return err
		}
		url := fmt.Sprintf("http://%s/flow/%s/json", net.JoinHostPort(nm.Spec.MonitoringEndpoint.IP, nm.Spec.MonitoringEndpoint.Port), flow.Name)
		fmt.Println(url)
		_, err = resty.R().
			SetHeader("Content-Type", "application/json").
			SetBody(jsonFlow).
			Put(url)

		return err
	}

	for _, flow := range networkMonitor.Spec.Flows {
		exists, err := flowExists(flow.Name)
		if err != nil {
			return err
		}
		if !exists {
			if err := installFlow(flow, networkMonitor); err != nil {
				return err
			}
			r.logger.Info("Installed flow successfully", "NetworkMonitor", flow.Name)
			r.recorder.Event(networkMonitor, v1alpha1.EventTypeNormal, v1alpha1.EventTypeReconciliation, fmt.Sprintf("Installed flow %s", flow.Name))
		}
	}

	return nil
}

func (r *ReconcileNetworkMonitor) installThresholds(networkMonitor *v1alpha1.NetworkMonitor) error {
	var sflowAddress = net.JoinHostPort(networkMonitor.Spec.MonitoringEndpoint.IP, networkMonitor.Spec.MonitoringEndpoint.Port)
	thresholdExists := func(thresholdName string) (bool, error) {
		resp, err := resty.R().Get("http://" + sflowAddress + fmt.Sprintf("/threshold/%s/json", thresholdName))
		if err != nil {
			return false, err
		}
		// TODO: check other 20X codes
		if resp.StatusCode() != http.StatusOK {
			return false, nil
		}
		return true, nil
	}

	installThreshold := func(threshold v1alpha1.Threshold, nm *v1alpha1.NetworkMonitor) error {
		jsonThreshold, err := json.Marshal(threshold)
		if err != nil {
			return err
		}

		url := fmt.Sprintf("http://%s/threshold/%s/json", sflowAddress, threshold.Name)
		_, err = resty.R().
			SetHeader("Content-Type", "application/json").
			SetBody(jsonThreshold).
			Put(url)

		return err
	}

	for _, threshold := range networkMonitor.Spec.Thresholds {
		exists, err := thresholdExists(threshold.Name)
		if err != nil {
			return err
		}
		if !exists {
			if err := installThreshold(threshold, networkMonitor); err != nil {
				return err
			}
			r.logger.Info("Installed thresholds successfully", "NetworkMonitor", threshold.Name)
			r.recorder.Event(networkMonitor, v1alpha1.EventTypeNormal, v1alpha1.EventTypeReconciliation, fmt.Sprintf("Installed threshold %s", threshold.Name))
		}
	}
	return nil
}

func (r *ReconcileNetworkMonitor) checkEvents(networkMonitor *v1alpha1.NetworkMonitor) ([]v1alpha1.NetworkEvent, error) {
	var (
		sflowAddress = net.JoinHostPort(networkMonitor.Spec.MonitoringEndpoint.IP, networkMonitor.Spec.MonitoringEndpoint.Port)
		url          = fmt.Sprintf("http://%s/events/json", sflowAddress)
	)
	resp, err := resty.R().SetQueryParams(map[string]string{
		"maxEvents": networkMonitor.Spec.EventsConfig.MaxEvents,
		"timeout":   networkMonitor.Spec.EventsConfig.Timeout,
	}).Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("request returned with non-ok status code %v", resp.StatusCode())
	}

	fmt.Println(string(resp.Body()))

	return nil, err
}

func getFlowForName(flowName string, flows []v1alpha1.Flow) (*v1alpha1.Flow, error) {
	for _, flow := range flows {
		if flowName == flow.Name {
			return &flow, nil
		}
	}
	return nil, fmt.Errorf("flow with key %s does not exist", flowName)
}
