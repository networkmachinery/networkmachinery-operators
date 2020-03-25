package controller

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils/apimachinery"
	authorizationv1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	LogKey = "NetworkDelayTest"
	// FinalizerName is the controlplane controller finalizer.
	FinalizerName = "networkmachinery.io/networkdelay"
)

// ReconcileMachineDeployment reconciles a MachineDeployment object.
type ReconcileNetworkDelayTest struct {
	config   *rest.Config
	logger   logr.Logger
	client   client.Client
	ctx      context.Context
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// InjectConfig implements inject.Config.
func (r *ReconcileNetworkDelayTest) InjectConfig(config *rest.Config) error {
	r.config = config
	return nil
}

func (r *ReconcileNetworkDelayTest) InjectClient(client client.Client) error {
	r.client = client
	return nil
}

func (r *ReconcileNetworkDelayTest) InjectStopChannel(stopCh <-chan struct{}) error {
	r.ctx = utils.ContextFromStopChannel(stopCh)
	return nil
}

func (r *ReconcileNetworkDelayTest) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	networkDelayTest := &v1alpha1.NetworkDelayTest{}
	if err := r.client.Get(r.ctx, request.NamespacedName, networkDelayTest); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return apimachinery.ReconcileErr(err)
	}

	if networkDelayTest.DeletionTimestamp != nil {
		return r.delete(r.ctx, networkDelayTest)
	}

	r.logger.Info("Reconciling Network Delay Test", "Name", networkDelayTest.Name)
	r.recorder.Event(networkDelayTest, v1alpha1.EventTypeNormal, v1alpha1.EventTypeReconciliation, "Reconciling NetworkDelayTest")

	return r.reconcile(r.ctx, networkDelayTest)
}

func (r *ReconcileNetworkDelayTest) reconcile(ctx context.Context, networkDelayTest *v1alpha1.NetworkDelayTest) (reconcile.Result, error) {
	// Finalizer is needed to make sure to clean up installed debugging tools after reconciliation
	if err := apimachinery.EnsureFinalizer(ctx, r.client, FinalizerName, networkDelayTest); err != nil {
		return apimachinery.ReconcileErr(err)
	}

	err := apimachinery.Can(ctx, r.client, &authorizationv1.ResourceAttributes{
		Namespace:   networkDelayTest.Spec.Source.Namespace,
		Verb:        "create",
		Resource:    "pods",
		Subresource: "exec",
		Group:       "",
		Name:        "",
	})
	if err != nil {
		return reconcile.Result{}, err
	}
	return r.reconcileNetworkDelay(ctx, networkDelayTest)

}

func (r *ReconcileNetworkDelayTest) delete(ctx context.Context, networkDelayTest *v1alpha1.NetworkDelayTest) (reconcile.Result, error) {
	hasFinalizer, err := apimachinery.HasFinalizer(networkDelayTest, FinalizerName)
	if err != nil {
		r.logger.Error(err, "Could not instantiate finalizer deletion")
		return apimachinery.ReconcileErr(err)
	}
	if !hasFinalizer {
		r.logger.Info("Deleting NetworkDelayTest causes a no-op as there is no finalizer.", LogKey, networkDelayTest.Name)
		return reconcile.Result{}, nil
	}

	r.logger.Info("Starting the deletion of the network Delay test ", LogKey, networkDelayTest.Name)
	r.recorder.Event(networkDelayTest, v1alpha1.EventTypeNormal, v1alpha1.EventTypeDeletion, "Deleting the network Delay test")

	// TODO: delete the installed tools from the pods

	if err := apimachinery.DeleteFinalizer(ctx, r.client, FinalizerName, networkDelayTest); err != nil {
		r.logger.Error(err, "Error removing finalizer from the NetworkMonitor resource", LogKey, networkDelayTest.Name)
		return apimachinery.ReconcileErr(err)
	}

	r.logger.Info("Deletion successful ", LogKey, networkDelayTest.Name)
	r.recorder.Event(networkDelayTest, v1alpha1.EventTypeNormal, v1alpha1.EventTypeDeletion, "Network Delay Test Deleted!")
	return reconcile.Result{}, nil
}
