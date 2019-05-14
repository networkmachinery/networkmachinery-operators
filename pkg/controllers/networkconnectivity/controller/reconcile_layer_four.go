package controller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils/apimachinery"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var NetcatStatusTypeMeta = metav1.TypeMeta{
	APIVersion: v1alpha1.SchemeGroupVersion.String(),
	Kind:       "NetcatStatus",
}

type NetcatOutput struct {
	state v1alpha1.NetcatResultState
}

func (r *ReconcileNetworkConnectivityTest) IPNetcat(ctx context.Context, status *v1alpha1.NetcatStatus, source *v1alpha1.NetworkSourceEndpoint, destination *v1alpha1.NetworkDestinationEndpoint) {
	netcatOut, err := NetCat(ctx, r.config, *source, destination.IP, destination.Port)
	if err != nil || netcatOut.state == v1alpha1.Refused {
		status.NetcatIPEndpoints = append(status.NetcatIPEndpoints, v1alpha1.NetcatIPEndpoint{
			IP:   destination.IP,
			Port: destination.Port,
			NetcatResult: v1alpha1.NetcatResult{
				State: v1alpha1.Refused,
			},
		})
	} else {
		status.NetcatIPEndpoints = append(status.NetcatIPEndpoints, v1alpha1.NetcatIPEndpoint{
			IP:   destination.IP,
			Port: destination.Port,
			NetcatResult: v1alpha1.NetcatResult{
				State: v1alpha1.Succeeded,
			},
		})
	}
}

func (r *ReconcileNetworkConnectivityTest) PodNetcat(ctx context.Context, status *v1alpha1.NetcatStatus, source *v1alpha1.NetworkSourceEndpoint, destination *v1alpha1.NetworkDestinationEndpoint) error {
	destinationPod := &corev1.Pod{}
	err := r.client.Get(ctx, client.ObjectKey{Namespace: destination.Namespace, Name: destination.Name}, destinationPod)
	if err != nil {
		return err
	}

	podIP := destinationPod.Status.PodIP
	if len(podIP) == 0 {
		return fmt.Errorf("could not find pod IP to netcat")
	}

	netcatOut, err := NetCat(ctx, r.config, *source, podIP, destination.Port)
	if err != nil || netcatOut.state == v1alpha1.Refused {
		status.NetcatPodEndpoints = append(status.NetcatPodEndpoints, v1alpha1.NetcatPodEndpoint{
			PodParams: v1alpha1.Params{
				Namespace: destination.Namespace,
				Name:      destination.Name,
				IP:        podIP,
				Port:      destination.Port,
			},
			NetcatResult: v1alpha1.NetcatResult{
				State: v1alpha1.Succeeded,
			},
		})
	} else {
		status.NetcatPodEndpoints = append(status.NetcatPodEndpoints, v1alpha1.NetcatPodEndpoint{
			PodParams: v1alpha1.Params{
				Namespace: destination.Namespace,
				Name:      destination.Name,
				IP:        podIP,
				Port:      destination.Port,
			},
			NetcatResult: v1alpha1.NetcatResult{
				State: v1alpha1.Succeeded,
			},
		})
	}
	return nil
}

func (r *ReconcileNetworkConnectivityTest) ServiceNetcat(ctx context.Context, status *v1alpha1.NetcatStatus, source *v1alpha1.NetworkSourceEndpoint, destination *v1alpha1.NetworkDestinationEndpoint) error {
	objectKey := client.ObjectKey{Namespace: destination.Namespace, Name: destination.Name}

	service := &corev1.Service{}
	err := r.client.Get(ctx, objectKey, service)
	if err != nil {
		return err
	}

	endpoints := &corev1.Endpoints{}
	err = r.client.Get(ctx, objectKey, endpoints)
	if err != nil {
		return err
	}

	var netcatIPEndpoints []v1alpha1.NetcatIPEndpoint
	// TODO: handle multiple subsets and separeate error from state
	for _, endpoint := range endpoints.Subsets[0].Addresses {
		endPointPort := strconv.Itoa(int(endpoints.Subsets[0].Ports[0].Port))
		netcatOut, err := NetCat(ctx, r.config, *source, endpoint.IP, endPointPort)
		if err != nil || netcatOut.state == v1alpha1.Refused {
			netcatIPEndpoints = append(netcatIPEndpoints, v1alpha1.NetcatIPEndpoint{
				IP:   endpoint.IP,
				Port: endPointPort,
				NetcatResult: v1alpha1.NetcatResult{
					State: v1alpha1.Refused,
				},
			})
		} else {
			netcatIPEndpoints = append(netcatIPEndpoints, v1alpha1.NetcatIPEndpoint{
				IP:   endpoint.IP,
				Port: endPointPort,
				NetcatResult: v1alpha1.NetcatResult{
					State: v1alpha1.Succeeded,
				},
			})
		}
	}

	status.NetcatServiceEndpoints = append(status.NetcatServiceEndpoints, v1alpha1.NetcatServiceEndpoint{
		ServiceParams: v1alpha1.Params{
			IP:        service.Spec.ClusterIP,
			Port:      strconv.Itoa(int(service.Spec.Ports[0].Port)),
			Name:      destination.Name,
			Namespace: destination.Namespace,
		},
		ServiceResults: netcatIPEndpoints,
	})
	return nil
}

func (r *ReconcileNetworkConnectivityTest) reconcileLayerFour(ctx context.Context, networkConnectivityTest *v1alpha1.NetworkConnectivityTest) (reconcile.Result, error) {
	var (
		status = &v1alpha1.NetcatStatus{
			TypeMeta: NetcatStatusTypeMeta,
		}
	)
	for _, destination := range networkConnectivityTest.Spec.Destinations {
		switch destination.Kind {
		case v1alpha1.IP:
			r.IPNetcat(ctx, status, &networkConnectivityTest.Spec.Source, &destination)
		case v1alpha1.Pod:
			err := r.PodNetcat(ctx, status, &networkConnectivityTest.Spec.Source, &destination)
			if err != nil {
				return apimachinery.ReconcileErr(err)
			}
		case v1alpha1.Service:
			err := r.ServiceNetcat(ctx, status, &networkConnectivityTest.Spec.Source, &destination)
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