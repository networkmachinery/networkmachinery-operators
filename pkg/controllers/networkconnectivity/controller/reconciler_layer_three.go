package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils/apimachinery"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	reconcilePeriod = 1 * time.Second

	PingStatusTypeMeta = metav1.TypeMeta{
		APIVersion: v1alpha1.SchemeGroupVersion.String(),
		Kind:       "PingStatus",
	}
)

type PingOutput struct {
	state         v1alpha1.PingResultState
	min, avg, max string
}

func (r *ReconcileNetworkConnectivityTest) IPPing(ctx context.Context, status *v1alpha1.PingStatus, source *v1alpha1.NetworkSourceEndpoint, destination string) {
	pingOut, err := Ping(ctx, r.config, *source, destination)
	if err != nil {
		status.PingIPEndpoints = append(status.PingIPEndpoints, v1alpha1.PingIPEndpoint{
			IP: destination,
			PingResult: v1alpha1.PingResult{
				State: v1alpha1.FailedPing,
			},
		})
	} else {
		status.PingIPEndpoints = append(status.PingIPEndpoints, v1alpha1.PingIPEndpoint{
			IP: destination,
			PingResult: v1alpha1.PingResult{
				State:   v1alpha1.SuccessPing,
				Average: pingOut.avg,
				Max:     pingOut.max,
				Min:     pingOut.min,
			},
		})
	}
}

func (r *ReconcileNetworkConnectivityTest) PodPing(ctx context.Context, status *v1alpha1.PingStatus, source *v1alpha1.NetworkSourceEndpoint, destination *v1alpha1.NetworkDestinationEndpoint) error {
	destinationPod := &corev1.Pod{}
	err := r.client.Get(ctx, client.ObjectKey{Namespace: destination.Namespace, Name: destination.Name}, destinationPod)
	if err != nil {
		status.PingPodEndpoints = append(status.PingPodEndpoints, v1alpha1.PingPodEndpoint{
			PodParams: v1alpha1.Params{
				Namespace: destination.Namespace,
				Name:      destination.Name,
			},
			PingResult: v1alpha1.PingResult{
				State: v1alpha1.FailedPing,
			},
		})
		return err
	}

	podIP := destinationPod.Status.PodIP
	if len(podIP) == 0 {
		return fmt.Errorf("could not find pod IP to ping")
	}

	pingOut, err := Ping(ctx, r.config, *source, podIP)
	if err != nil {
		status.PingPodEndpoints = append(status.PingPodEndpoints, v1alpha1.PingPodEndpoint{
			PodParams: v1alpha1.Params{
				Namespace: destination.Namespace,
				Name:      destination.Name,
				IP:        podIP,
			},
			PingResult: v1alpha1.PingResult{
				State: v1alpha1.FailedPing,
			},
		})
	} else {
		status.PingPodEndpoints = append(status.PingPodEndpoints, v1alpha1.PingPodEndpoint{
			PodParams: v1alpha1.Params{
				Namespace: destination.Namespace,
				Name:      destination.Name,
				IP:        podIP,
			},
			PingResult: v1alpha1.PingResult{
				State:   v1alpha1.SuccessPing,
				Average: pingOut.avg,
				Max:     pingOut.max,
				Min:     pingOut.min,
			},
		})
	}
	return nil
}

func (r *ReconcileNetworkConnectivityTest) ServicePing(ctx context.Context, status *v1alpha1.PingStatus, source *v1alpha1.NetworkSourceEndpoint, destination *v1alpha1.NetworkDestinationEndpoint) error {
	endpoints := &corev1.Endpoints{}
	err := r.client.Get(ctx, client.ObjectKey{Namespace: destination.Namespace, Name: destination.Name}, endpoints)
	if err != nil {
		return err
	}

	var pingIPEndpoints []v1alpha1.PingIPEndpoint
	// TODO: handle multiple subsets
	for _, endpoint := range endpoints.Subsets[0].Addresses {
		pingOut, err := Ping(ctx, r.config, *source, endpoint.IP)
		if err != nil {
			pingIPEndpoints = append(pingIPEndpoints, v1alpha1.PingIPEndpoint{
				IP: endpoint.IP,
				PingResult: v1alpha1.PingResult{
					State: v1alpha1.FailedPing,
				},
			})
		} else {
			pingIPEndpoints = append(pingIPEndpoints, v1alpha1.PingIPEndpoint{
				IP: endpoint.IP,
				PingResult: v1alpha1.PingResult{
					State:   v1alpha1.SuccessPing,
					Average: pingOut.avg,
					Max:     pingOut.max,
					Min:     pingOut.min,
				},
			})
		}
	}

	service := &corev1.Service{}
	err = r.client.Get(ctx, client.ObjectKey{Namespace: destination.Namespace, Name: destination.Name}, service)
	if err != nil {
		return err
	}

	status.PingServiceEndpoint = append(status.PingServiceEndpoint, v1alpha1.PingServiceEndpoint{
		ServiceParams: v1alpha1.Params{
			IP:        service.Spec.ClusterIP,
			Name:      destination.Name,
			Namespace: destination.Namespace,
		},
		ServiceResults: pingIPEndpoints,
	})

	return nil
}

func (r *ReconcileNetworkConnectivityTest) reconcileLayerThree(ctx context.Context, networkConnectivityTest *v1alpha1.NetworkConnectivityTest) (reconcile.Result, error) {
	var (
		status = &v1alpha1.PingStatus{
			TypeMeta: PingStatusTypeMeta,
		}
	)
	for _, destination := range networkConnectivityTest.Spec.Destinations {
		switch destination.Kind {
		case v1alpha1.IP:
			r.IPPing(ctx, status, &networkConnectivityTest.Spec.Source, destination.IP)
		case v1alpha1.Pod:
			err := r.PodPing(ctx, status, &networkConnectivityTest.Spec.Source, &destination)
			if err != nil {
				return apimachinery.ReconcileErr(err)
			}
		case v1alpha1.Service:
			err := r.ServicePing(ctx, status, &networkConnectivityTest.Spec.Source, &destination)
			if err != nil {
				return apimachinery.ReconcileErr(err)
			}
		}
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
