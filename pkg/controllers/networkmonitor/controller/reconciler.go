package controller

import (
	"context"
	"fmt"

	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"

	"github.com/go-logr/logr"
	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
	r.logger.Info("Reconciling NetworkMonitor!")

	networkMonitor := &v1alpha1.NetworkMonitor{}
	if err := r.client.Get(r.ctx, request.NamespacedName, networkMonitor); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	for _, flow := range networkMonitor.Spec.Flows {
		fmt.Println(flow.Name, flow.Value, flow.ActiveTimeout, flow.Filter, flow.FlowStart, flow.Keys)
	}

	return reconcile.Result{}, nil
}
