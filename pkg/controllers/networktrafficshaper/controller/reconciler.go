package controller

import (
	"context"
	"time"

	"k8s.io/client-go/rest"

	"k8s.io/apimachinery/pkg/labels"

	"github.com/networkmachinery/networkmachinery-operators/pkg/utils/apimachinery"

	"github.com/go-logr/logr"
	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// FinalizerName is the controlplane controller finalizer.
	FinalizerName = "networkmachinery.io/networkmonitor"
	LogKey        = "network-traffic-shaper"
)

// ReconcileMachineDeployment reconciles a MachineDeployment object.
type ReconcileNetworkTrafficShaper struct {
	logger logr.Logger
	client client.Client
	config *rest.Config

	ctx      context.Context
	scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// InjectConfig implements inject.Config.
func (r *ReconcileNetworkTrafficShaper) InjectConfig(config *rest.Config) error {
	r.config = config
	return nil
}
func (r *ReconcileNetworkTrafficShaper) InjectClient(client client.Client) error {
	r.client = client
	return nil
}

func (r *ReconcileNetworkTrafficShaper) InjectStopChannel(stopCh <-chan struct{}) error {
	r.ctx = utils.ContextFromStopChannel(stopCh)
	return nil
}

func (r *ReconcileNetworkTrafficShaper) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	networkTrafficShaper := &v1alpha1.NetworkTrafficShaper{}
	if err := r.client.Get(r.ctx, request.NamespacedName, networkTrafficShaper); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return apimachinery.ReconcileErr(err)
	}

	if networkTrafficShaper.DeletionTimestamp != nil {
		return r.delete(r.ctx, networkTrafficShaper)
	}

	r.logger.Info("Reconciling NetworkTrafficShaper", "Name", networkTrafficShaper.Name)
	return r.reconcile(r.ctx, networkTrafficShaper)
}

func (r *ReconcileNetworkTrafficShaper) reconcile(ctx context.Context, networkTrafficShaper *v1alpha1.NetworkTrafficShaper) (reconcile.Result, error) {
	// Finalizer is needed to make sure to clean up installed debugging tools after reconciliation
	if err := apimachinery.EnsureFinalizer(ctx, r.client, FinalizerName, networkTrafficShaper); err != nil {
		return apimachinery.ReconcileErr(err)
	}

	shapePodList := func(ctx context.Context, pods *corev1.PodList, device, value, shapeType string) error {
		for _, pod := range pods.Items {
			if err := shapeTraffic(ctx, r.config, pod.Namespace, pod.Name, device, value, shapeType); err != nil {
				return err
			}
		}
		return nil
	}

	for _, target := range networkTrafficShaper.Spec.Targets {
		switch target.Kind {
		case v1alpha1.Selector:
			podList, err := utils.GetPodsByLabels(ctx, r.client, labels.SelectorFromSet(labels.Set(target.SourceSelector.MatchLabels)), target.Namespace)
			if err != nil {
				return apimachinery.ReconcileErr(err)
			}
			if err := shapePodList(ctx, podList, target.ShaperConfig.Device, target.ShaperConfig.Value, string(target.ShaperConfig.Type)); err != nil {
				return apimachinery.ReconcileErr(err)
			}
		case v1alpha1.Pod:
			pod := &corev1.Pod{}
			if err := r.client.Get(ctx, client.ObjectKey{Namespace: target.Namespace, Name: target.Name}, pod); err != nil {
				return apimachinery.ReconcileErr(err)
			}
			if err := shapeTraffic(ctx, r.config, pod.Namespace, pod.Name, target.ShaperConfig.Device, target.ShaperConfig.Value, string(target.ShaperConfig.Type)); err != nil {
				return apimachinery.ReconcileErr(err)
			}
		}
	}
	return reconcile.Result{
		RequeueAfter: 10 * time.Second,
	}, nil
}

func (r *ReconcileNetworkTrafficShaper) delete(ctx context.Context, networkTrafficShaper *v1alpha1.NetworkTrafficShaper) (reconcile.Result, error) {
	hasFinalizer, err := apimachinery.HasFinalizer(networkTrafficShaper, FinalizerName)
	if err != nil {
		r.logger.Error(err, "Could not instantiate finalizer deletion")
		return apimachinery.ReconcileErr(err)
	}
	if !hasFinalizer {
		r.logger.Info("Deleting NetworkTrafficShaper causes a no-op as there is no finalizer.", LogKey, networkTrafficShaper.Name)
		return reconcile.Result{}, nil
	}

	r.logger.Info("Starting the deletion of the network connectivity test ", LogKey, networkTrafficShaper.Name)
	r.recorder.Event(networkTrafficShaper, v1alpha1.EventTypeNormal, v1alpha1.EventTypeDeletion, "Deleting the NetworkTrafficShaper")

	undoShapePodList := func(ctx context.Context, pods *corev1.PodList, device string) error {
		for _, pod := range pods.Items {
			if err := undoShape(ctx, r.config, pod.Namespace, pod.Name, device); err != nil {
				return err
			}
		}
		return nil
	}

	for _, target := range networkTrafficShaper.Spec.Targets {
		switch target.Kind {
		case v1alpha1.Selector:
			podList, err := utils.GetPodsByLabels(ctx, r.client, labels.SelectorFromSet(labels.Set(target.SourceSelector.MatchLabels)), target.Namespace)
			if err != nil {
				return apimachinery.ReconcileErr(err)
			}
			if err := undoShapePodList(ctx, podList, target.ShaperConfig.Device); err != nil {
				return apimachinery.ReconcileErr(err)
			}
		case v1alpha1.Pod:
			pod := &corev1.Pod{}
			if err := r.client.Get(ctx, client.ObjectKey{Namespace: target.Namespace, Name: target.Name}, pod); err != nil {
				return apimachinery.ReconcileErr(err)
			}
			if err := undoShape(ctx, r.config, pod.Namespace, pod.Name, target.ShaperConfig.Device); err != nil {
				return apimachinery.ReconcileErr(err)
			}
		}
	}

	if err := apimachinery.DeleteFinalizer(ctx, r.client, FinalizerName, networkTrafficShaper); err != nil {
		r.logger.Error(err, "Error removing finalizer from the NetworkTrafficShaper resource", LogKey, networkTrafficShaper.Name)
		return apimachinery.ReconcileErr(err)
	}

	return reconcile.Result{}, nil
}
