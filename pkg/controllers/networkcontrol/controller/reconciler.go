package controller

import (
	"context"
	"time"

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
type ReconcileNetworkController struct {
	logger   logr.Logger
	client   client.Client
	ctx      context.Context
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

func (r *ReconcileNetworkController) InjectClient(client client.Client) error {
	r.client = client
	return nil
}

func (r *ReconcileNetworkController) InjectStopChannel(stopCh <-chan struct{}) error {
	r.ctx = utils.ContextFromStopChannel(stopCh)
	return nil
}

func (r *ReconcileNetworkController) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	networkNotification := &v1alpha1.NetworkNotification{}
	if err := r.client.Get(r.ctx, request.NamespacedName, networkNotification); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return apimachinery.ReconcileErr(err)
	}

	r.logger.Info("Received a Network Notification", "Name", networkNotification.Name)
	r.recorder.Event(networkNotification, v1alpha1.EventTypeNormal, v1alpha1.EventTypeReconciliation, "Received a Network Notification")

	return r.reconcile(r.ctx, networkNotification)
}

// TODO: implement real actions
func (r *ReconcileNetworkController) reconcile(ctx context.Context, networkNotification *v1alpha1.NetworkNotification) (reconcile.Result, error) {
	r.logger.Info("Who DARES disrupt my Network, YOU SHALL NOT PASS", "Name", networkNotification.Name)
	r.logger.Info("I CAN Use OpenFlow now, I WILL FIND YOU AND I WILL KILL YOU")
	r.logger.Info("I Know where you live", "Your SWITCH", networkNotification.Spec.Event.Event.Agent, "Your PORT", networkNotification.Spec.Event.Event.DataSource)
	return reconcile.Result{
		RequeueAfter: 30 * time.Second,
	}, nil
}
