package controller

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils/apimachinery"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	LogKey = "NetworkConnnectivityTest"
	// FinalizerName is the controlplane controller finalizer.
	FinalizerName       = "networkmachinery.io/networkconnectivity"
)

var (
	reconcilePeriod = 10* time.Second

	// StatusTypeMeta is the TypeMeta of the GCP InfrastructureStatus
	StatusTypeMeta = metav1.TypeMeta{
		APIVersion: v1alpha1.SchemeGroupVersion.String(),
		Kind:       "InfrastructureStatus",
	}
)



type PingOutput struct {
	state v1alpha1.PingResultState
	min,avg,max string
}
// ReconcileMachineDeployment reconciles a MachineDeployment object.
type ReconcileNetworkConnectivityTest struct {
	config *rest.Config
	logger   logr.Logger
	client   client.Client
	ctx      context.Context
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// InjectConfig implements inject.Config.
func (r *ReconcileNetworkConnectivityTest) InjectConfig(config *rest.Config) error {
	r.config = config
	return nil
}


func (r *ReconcileNetworkConnectivityTest) InjectClient(client client.Client) error {
	r.client = client
	return nil
}

func (r *ReconcileNetworkConnectivityTest) InjectStopChannel(stopCh <-chan struct{}) error {
	r.ctx = utils.ContextFromStopChannel(stopCh)
	return nil
}

func (r *ReconcileNetworkConnectivityTest) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	networkConnectivityTest := &v1alpha1.NetworkConnectivityTest{}
	if err := r.client.Get(r.ctx, request.NamespacedName, networkConnectivityTest); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return apimachinery.ReconcileErr(err)
	}

	if networkConnectivityTest.DeletionTimestamp != nil {
		return r.delete(r.ctx, networkConnectivityTest)
	}

	r.logger.Info("Reconciling Network Connectivity Test", "Name", networkConnectivityTest.Name)
	r.recorder.Event(networkConnectivityTest, v1alpha1.EventTypeNormal, v1alpha1.EventTypeReconciliation, "Reconciling NetworkConnectivityTest")

	return r.reconcile(r.ctx, networkConnectivityTest)
}

func (r *ReconcileNetworkConnectivityTest) delete(ctx context.Context, networkConnectivityTest *v1alpha1.NetworkConnectivityTest) (reconcile.Result, error) {
	hasFinalizer, err := apimachinery.HasFinalizer(networkConnectivityTest, FinalizerName)
	if err != nil {
		r.logger.Error(err, "Could not instantiate finalizer deletion")
		return apimachinery.ReconcileErr(err)
	}
	if !hasFinalizer {
		r.logger.Info("Deleting NetworkConnectivityTest causes a no-op as there is no finalizer.", LogKey, networkConnectivityTest.Name)
		return reconcile.Result{}, nil
	}

	r.logger.Info("Starting the deletion of the network connectivity test ", LogKey, networkConnectivityTest.Name)
	r.recorder.Event(networkConnectivityTest, v1alpha1.EventTypeNormal, v1alpha1.EventTypeDeletion, "Deleting the network connectivity test")

	// DO stuff here

	if err := apimachinery.DeleteFinalizer(ctx, r.client, FinalizerName, networkConnectivityTest); err != nil {
		r.logger.Error(err, "Error removing finalizer from the NetworkMonitor resource", LogKey, networkConnectivityTest.Name)
		return apimachinery.ReconcileErr(err)
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileNetworkConnectivityTest) reconcile(ctx context.Context, networkConnectivityTest *v1alpha1.NetworkConnectivityTest) (reconcile.Result, error) {
	// Finalizer is needed to make sure to clean up installed debugging tools after reconciliation
	if err := apimachinery.EnsureFinalizer(ctx, r.client, FinalizerName, networkConnectivityTest); err != nil {
		return apimachinery.ReconcileErr(err)
	}
	switch networkConnectivityTest.Spec.Layer {
	case "3":
		return r.reconcileLayerThree(ctx, networkConnectivityTest)
	case "4":
		return r.reconcileLayerFour(ctx, networkConnectivityTest)
	}
	return reconcile.Result{
		RequeueAfter: reconcilePeriod,
	}, nil
}


func (r *ReconcileNetworkConnectivityTest) Ping(ctx context.Context, source v1alpha1.NetworkSourceEndpoint, host string) (*PingOutput, error) {
	var stdOut, stdErr bytes.Buffer
	execOpts := utils.ExecOptions{
		Namespace: source.Namespace,
		Name: source.Name,
		Command: fmt.Sprintf("ping -c3 %s", host),
		Container: source.Container,
		StandardCmdOpts: utils.StandardCmdOpts{
			StdErr: &stdErr,
			StdOut: &stdOut,
		},
	}

	err := utils.PodExec(ctx, r.config,execOpts)
	if err != nil {
		return &PingOutput{state: v1alpha1.FailedPing}, err
	}

	ping := &utils.Ping{}
	utils.ParsePingOutput(stdOut.Bytes(), ping)

	return &PingOutput{
		state: v1alpha1.SuccessPing,
		min: ping.Min(),
		max: ping.Max(),
		avg: ping.Average(),
	}, nil
}


func (r *ReconcileNetworkConnectivityTest) reconcileLayerThree(ctx context.Context, networkConnectivityTest *v1alpha1.NetworkConnectivityTest) (reconcile.Result, error)  {
	var (
		status = &v1alpha1.PingStatus{
			TypeMeta: StatusTypeMeta,
			PingIPEndpoints: []v1alpha1.PingIPEndpoint{},
		}
	)
	for _, destination :=range networkConnectivityTest.Spec.Destinations {
		pingOut, err := r.Ping(ctx, networkConnectivityTest.Spec.Source, destination.IP)
		if err != nil {
			status.PingIPEndpoints = append(status.PingIPEndpoints, v1alpha1.PingIPEndpoint{
				IP: destination.IP,
				PingResult: v1alpha1.PingResult{
					State: v1alpha1.FailedPing,
				},
			})
			fmt.Println(status.PingIPEndpoints)
		} else {
			status.PingIPEndpoints = append(status.PingIPEndpoints, v1alpha1.PingIPEndpoint{
				IP: destination.IP,
				PingResult: v1alpha1.PingResult{
					State: v1alpha1.SuccessPing,
					Average: pingOut.avg,
					Max: pingOut.max,
					Min: pingOut.min,
				},
			})
		}
		fmt.Println(pingOut)
		fmt.Println(status.PingIPEndpoints)
	}

	err := apimachinery.TryUpdateStatus(ctx, retry.DefaultBackoff, r.client, networkConnectivityTest, func() error {
		networkConnectivityTest.Status.TestStatus = &runtime.RawExtension{Object: status}
		return nil
	})
	if err != nil {
		return apimachinery.ReconcileErr(err)
	}

	return reconcile.Result{
		RequeueAfter: reconcilePeriod,
	}, nil
}

func (r *ReconcileNetworkConnectivityTest) reconcileLayerFour(ctx context.Context, networkConnectivityTest *v1alpha1.NetworkConnectivityTest) (reconcile.Result, error) {
	//reader, err := utils.PodExecByLabel(
	//	ctx,
	//	r.config,
	//	r.client,
	//	apiserverLabels,
	//	APIServer,
	//	fmt.Sprintf("apt-get update && apt-get -y install netcat && nc -z -w5 %s %s", host, port),
	//	shootTestOperations.ShootSeedNamespace())
	return reconcile.Result{
		RequeueAfter: reconcilePeriod,
	}, nil
}

func StringSet(s string) bool {
	return len(s) == 0
}
