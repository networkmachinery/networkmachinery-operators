package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apimachineryerror "github.com/networkmachinery/networkmachinery-operators/pkg/utils/apimachinery/error"

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
	FinalizerName       = "networkmachinery.io/networkmonitor"
	NetworkNotification = "network-notification-"
	ContentType         = "Content-Type"
	ApplicationJSON     = "application/json"
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
		return apimachinery.ReconcileErr(err)
	}

	if networkMonitor.DeletionTimestamp != nil {
		return r.delete(r.ctx, networkMonitor)
	}

	r.logger.Info("Reconciling NetworkMonitor", "Name", networkMonitor.Name)
	r.recorder.Event(networkMonitor, v1alpha1.EventTypeNormal, v1alpha1.EventTypeReconciliation, "Reconciling NetworkMonitor")

	return r.reconcile(r.ctx, networkMonitor)
}

func (r *ReconcileNetworkMonitor) delete(ctx context.Context, networkmonitor *v1alpha1.NetworkMonitor) (reconcile.Result, error) {
	hasFinalizer, err := apimachinery.HasFinalizer(networkmonitor, FinalizerName)
	if err != nil {
		r.logger.Error(err, "Could not instantiate finalizer deletion")
		return apimachinery.ReconcileErr(err)
	}
	if !hasFinalizer {
		r.logger.Info("Deleting NetworkMonitor causes a no-op as there is no finalizer.", "NetworkMonitor", networkmonitor.Name)
		return reconcile.Result{}, nil
	}

	r.logger.Info("Starting the deletion of network monitor", "NetworkMonitor", networkmonitor.Name)
	r.recorder.Event(networkmonitor, v1alpha1.EventTypeNormal, v1alpha1.EventTypeDeletion, "Deleting the Network Monitor")

	err = r.DeleteFlows(ctx, networkmonitor)
	if err != nil {
		return apimachinery.ReconcileErr(err)
	}

	err = r.DeleteThresholds(ctx, networkmonitor)
	if err != nil {
		return apimachinery.ReconcileErr(err)
	}

	r.logger.Info("Successfully cleaned up flows and thresholds", "NetworkMonitor", networkmonitor.Name)
	r.recorder.Event(networkmonitor, v1alpha1.EventTypeNormal, v1alpha1.EventTypeDeletion, "Deleted flows and thresholds")

	if err := apimachinery.DeleteFinalizer(ctx, r.client, FinalizerName, networkmonitor); err != nil {
		r.logger.Error(err, "Error removing finalizer from the NetworkMonitor resource", "Network Monitor", networkmonitor.Name)
		return apimachinery.ReconcileErr(err)
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileNetworkMonitor) reconcile(ctx context.Context, networkMonitor *v1alpha1.NetworkMonitor) (reconcile.Result, error) {
	if err := apimachinery.EnsureFinalizer(ctx, r.client, FinalizerName, networkMonitor); err != nil {
		return apimachinery.ReconcileErr(err)
	}

	err := r.InstallFlows(ctx, networkMonitor)
	if err != nil {
		return apimachinery.ReconcileErr(err)
	}

	err = r.InstallThresholds(ctx, networkMonitor)
	if err != nil {
		return apimachinery.ReconcileErr(err)
	}

	events, err := r.CheckEvents(ctx, networkMonitor)
	if err != nil {
		switch err.(type) {
		case *apimachineryerror.ZeroEventsError:
			r.logger.Info(err.Error(), "NetworkMonitor", networkMonitor.Name)
			return reconcile.Result{}, nil
		default:
			return apimachinery.ReconcileErr(err)
		}
	}

	for _, event := range events {
		err := r.createNetworkNotification(ctx, &event, networkMonitor)
		if err != nil {
			return apimachinery.ReconcileErr(err)
		}
	}
	// this is to allow for polling fresh events out of sFlow
	return reconcile.Result{
		RequeueAfter: 5 * time.Second,
	}, nil
}

func (r *ReconcileNetworkMonitor) InstallFlows(ctx context.Context, networkMonitor *v1alpha1.NetworkMonitor) error {
	var sflowAddress = net.JoinHostPort(networkMonitor.Spec.MonitoringEndpoint.IP, networkMonitor.Spec.MonitoringEndpoint.Port)

	flowExists := func(flowName string) (bool, error) {
		resp, err := resty.R().SetContext(ctx).Get("http://" + sflowAddress + fmt.Sprintf("/flow/%s/json", flowName))
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
		_, err = resty.R().
			SetContext(ctx).
			SetHeader(ContentType, ApplicationJSON).
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

func (r *ReconcileNetworkMonitor) DeleteFlows(ctx context.Context, networkMonitor *v1alpha1.NetworkMonitor) error {
	var sflowAddress = net.JoinHostPort(networkMonitor.Spec.MonitoringEndpoint.IP, networkMonitor.Spec.MonitoringEndpoint.Port)
	deleteFlow := func(flowName string, nm *v1alpha1.NetworkMonitor) error {
		resp, err := resty.R().SetContext(ctx).Delete("http://" + sflowAddress + fmt.Sprintf("/flow/%s/json", flowName))
		if err != nil {
			return err
		}
		// TODO: check other 20X codes
		if resp.StatusCode() != http.StatusOK {
			return nil
		}
		return nil
	}

	for _, flow := range networkMonitor.Spec.Flows {
		err := deleteFlow(flow.Name, networkMonitor)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReconcileNetworkMonitor) InstallThresholds(ctx context.Context, networkMonitor *v1alpha1.NetworkMonitor) error {
	var sflowAddress = net.JoinHostPort(networkMonitor.Spec.MonitoringEndpoint.IP, networkMonitor.Spec.MonitoringEndpoint.Port)
	thresholdExists := func(thresholdName string) (bool, error) {
		resp, err := resty.R().SetContext(ctx).Get("http://" + sflowAddress + fmt.Sprintf("/threshold/%s/json", thresholdName))
		if err != nil {
			return false, err
		}
		// TODO: check other 20X codes
		if resp.StatusCode() != http.StatusOK {
			return false, nil
		}

		var threshold v1alpha1.Threshold
		err = json.Unmarshal(resp.Body(), &threshold)
		if err != nil {
			return false, err
		}

		specThreshold, err := getThresholdForName(thresholdName, networkMonitor.Spec.Thresholds)
		if err != nil {
			return false, err
		}

		if !Equal(&threshold, specThreshold) {
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
		if len(threshold.ByFlow) > 0 {
			_, err = resty.R().
				SetContext(ctx).
				SetHeader(ContentType, ApplicationJSON).
				SetQueryParam("byFlow", "true").
				SetBody(jsonThreshold).
				Put(url)
		} else {
			_, err = resty.R().
				SetContext(ctx).
				SetHeader(ContentType, ApplicationJSON).
				SetBody(jsonThreshold).
				Put(url)
		}

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

func (r *ReconcileNetworkMonitor) DeleteThresholds(ctx context.Context, networkMonitor *v1alpha1.NetworkMonitor) error {
	var sflowAddress = net.JoinHostPort(networkMonitor.Spec.MonitoringEndpoint.IP, networkMonitor.Spec.MonitoringEndpoint.Port)
	deleteThreshold := func(thresholdName string, nm *v1alpha1.NetworkMonitor) error {
		resp, err := resty.R().SetContext(ctx).Delete("http://" + sflowAddress + fmt.Sprintf("/threshold/%s/json", thresholdName))
		if err != nil {
			return err
		}
		// TODO: check other 20X codes
		if resp.StatusCode() != http.StatusOK {
			return nil
		}
		return nil
	}

	for _, threshold := range networkMonitor.Spec.Thresholds {
		err := deleteThreshold(threshold.Name, networkMonitor)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReconcileNetworkMonitor) CheckEvents(ctx context.Context, networkMonitor *v1alpha1.NetworkMonitor) ([]v1alpha1.Event, error) {
	var (
		sflowAddress = net.JoinHostPort(networkMonitor.Spec.MonitoringEndpoint.IP, networkMonitor.Spec.MonitoringEndpoint.Port)
		url          = fmt.Sprintf("http://%s/events/json", sflowAddress)
	)
	resp, err := resty.R().SetContext(ctx).SetQueryParams(map[string]string{
		"maxEvents": networkMonitor.Spec.EventsConfig.MaxEvents,
		"timeout":   networkMonitor.Spec.EventsConfig.Timeout,
	}).Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("request returned with non-ok status code %v, reason %v", resp.StatusCode(), resp.Error())
	}

	var events []v1alpha1.Event
	err = json.Unmarshal(resp.Body(), &events)
	if err != nil {
		return nil, err
	}

	// we don't care about events that we didn't define thresholds for
	events = filterEvents(events, networkMonitor, func(e v1alpha1.Event, monitor *v1alpha1.NetworkMonitor) bool {
		var definedThresholds []string
		for _, threshold := range monitor.Spec.Thresholds {
			definedThresholds = append(definedThresholds, threshold.Name)
		}
		if in(e.ThresholdID, definedThresholds) {
			return true
		}
		return false
	})

	if len(events) == 0 {
		return nil, apimachineryerror.NewZeroEventsError("no related events for defined flows / thresholds found to return")
	}

	return events, err
}

func (r *ReconcileNetworkMonitor) createNetworkNotification(ctx context.Context, event *v1alpha1.Event, networkMonitor *v1alpha1.NetworkMonitor) error {
	eventThreshold, err := getThresholdForName(event.ThresholdID, networkMonitor.Spec.Thresholds)
	if err != nil {
		return err
	}

	eventFlow, err := getFlowForName(eventThreshold.FlowName, networkMonitor.Spec.Flows)
	if err != nil {
		return err
	}

	networkNotification := &v1alpha1.NetworkNotification{
		ObjectMeta: metav1.ObjectMeta{
			Name: NetworkNotification + strconv.Itoa(int(event.EventID)),
		},
		Spec: v1alpha1.NetworkNotificationSpec{
			Event: v1alpha1.NetworkEvent{
				Event: *event,
				Flow:  *eventFlow,
			},
		},
	}

	return apimachinery.CreateOrUpdate(ctx, r.client, networkNotification, nil)
}
