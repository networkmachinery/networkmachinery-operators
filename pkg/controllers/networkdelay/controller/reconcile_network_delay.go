package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils/apimachinery"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	reconcilePeriod = 2 * time.Second

	QperfStatusTypeMeta = metav1.TypeMeta{
		APIVersion: v1alpha1.SchemeGroupVersion.String(),
		Kind:       "QperfStatus",
	}
)

type QperfOutput struct {
	state              v1alpha1.QperfResultState
	udpdelay, tcpdelay string
}

func (r *ReconcileNetworkDelayTest) PodQperf(ctx context.Context, status *v1alpha1.QperfStatus, source *v1alpha1.NetworkSourceEndpoint, destination *v1alpha1.NetworkDestinationEndpoint, protocol string) error {
	destinationPod := &corev1.Pod{}
	if err := r.client.Get(ctx, client.ObjectKey{Namespace: destination.Namespace, Name: destination.Name}, destinationPod); err != nil {
		status.QperfPodEndpoints = append(status.QperfPodEndpoints, v1alpha1.QperfPodEndpoint{
			PodParams: v1alpha1.Params{
				Namespace: destination.Namespace,
				Name:      destination.Name,
			},
			QperfResult: v1alpha1.QperfResult{
				State: v1alpha1.FailedQperfCmd,
			},
		})
		return err
	}

	podIP := destinationPod.Status.PodIP
	if len(podIP) == 0 {
		return fmt.Errorf("could not find pod IP")
	}

	qperfOut, err := Qperf(ctx, r.config, *source, *destination, protocol, podIP)

	if err != nil {
		status.QperfPodEndpoints = append(status.QperfPodEndpoints, v1alpha1.QperfPodEndpoint{
			PodParams: v1alpha1.Params{
				Namespace: destination.Namespace,
				Name:      destination.Name,
				IP:        podIP,
			},
			QperfResult: v1alpha1.QperfResult{
				State: v1alpha1.FailedQperfCmd,
			},
		})
	} else {
		status.QperfPodEndpoints = append(status.QperfPodEndpoints, v1alpha1.QperfPodEndpoint{
			PodParams: v1alpha1.Params{
				Namespace: destination.Namespace,
				Name:      destination.Name,
				IP:        podIP,
			},
			QperfResult: v1alpha1.QperfResult{
				State:    v1alpha1.SuccessQperfCmd,
				TCPDelay: qperfOut.tcpdelay,
				UDPDelay: qperfOut.udpdelay,
			},
		})
	}
	return nil
}

func (r *ReconcileNetworkDelayTest) NodeQperf(ctx context.Context, status *v1alpha1.QperfStatus, source *v1alpha1.NetworkSourceEndpoint, destination *v1alpha1.NetworkDestinationEndpoint, protocol string) error {

	label := make(map[string]string)
	label["name"] = "qperf-ds"

	// check if ds is deployed
	dsPodList := &corev1.PodList{}
	if err := r.client.List(ctx, dsPodList, client.MatchingLabels(label)); err != nil {
		return err
	}
	if len(dsPodList.Items) == 0 {
		// apply ds
		ds := &v1.DaemonSet{
			TypeMeta: metav1.TypeMeta{
				Kind:       "DaemonSet",
				APIVersion: "apps/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Labels:    map[string]string{"k8s-app": "qperf-ds"},
				Name:      "qperf-ds",
				Namespace: "default",
			},
			Spec: v1.DaemonSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"name": "qperf-ds"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"name": "qperf-ds"},
					},
					Spec: corev1.PodSpec{
						DNSPolicy: "ClusterFirst",
						Containers: []corev1.Container{
							{
								Name:            "test",
								Image:           "ubuntu:latest",
								ImagePullPolicy: "IfNotPresent",
								Command:         []string{"sleep"},
								Args:            []string{"3600"},
							},
						},
					},
				},
			},
		}

		if err := r.client.Create(ctx, ds); err != nil {
			return err
		}
		return nil
	}

	// check if at least two nodes are in the cluster
	if len(dsPodList.Items) < 2 {
		status.QperfNodeEndpoints = append(status.QperfNodeEndpoints, v1alpha1.QperfNodeEndpoint{
			NodeParams: v1alpha1.NodeParams{
				Name: dsPodList.Items[0].Name,
			},
			NodeResults: v1alpha1.QperfResult{
				State: v1alpha1.FailedQperfCmd,
			},
		})
		return errors.New("At least two nodes are required")
	}
	// if name is empty take random node for measurement
	if destination.Name == "" {
		// take two random nodes
		nodeNameDestination := dsPodList.Items[1].Spec.NodeName
		var destinationAddress string
		node := &corev1.Node{}
		err := r.client.Get(ctx, client.ObjectKey{Namespace: "", Name: nodeNameDestination}, node)
		for _, address := range node.Status.Addresses {
			if address.Type == "InternalIP" {
				destinationAddress = address.Address
				break
			}
		}

		podIP := dsPodList.Items[1].Status.PodIP
		if len(podIP) == 0 {
			return fmt.Errorf("could not find pod IP")
		}

		source.Name = dsPodList.Items[0].Name
		source.Namespace = dsPodList.Items[0].Namespace
		source.Container = "net-server"
		destination.Kind = "pod"
		destination.Name = dsPodList.Items[1].Name
		destination.Namespace = dsPodList.Items[1].Namespace
		qperfOut, err := Qperf(ctx, r.config, *source, *destination, protocol, podIP)
		fmt.Printf("QPERFOUT for NODE: %v\n", qperfOut)
		if err != nil {
			status.QperfNodeEndpoints = append(status.QperfNodeEndpoints, v1alpha1.QperfNodeEndpoint{
				NodeParams: v1alpha1.NodeParams{
					Name: nodeNameDestination,
					IP:   destinationAddress,
				},
				NodeResults: v1alpha1.QperfResult{
					State: v1alpha1.FailedQperfCmd,
				},
			})
		} else {
			status.QperfNodeEndpoints = append(status.QperfNodeEndpoints, v1alpha1.QperfNodeEndpoint{
				NodeParams: v1alpha1.NodeParams{
					Name: nodeNameDestination,
					IP:   destinationAddress,
				},
				NodeResults: v1alpha1.QperfResult{
					State:    v1alpha1.SuccessQperfCmd,
					TCPDelay: qperfOut.tcpdelay,
					UDPDelay: qperfOut.udpdelay,
				},
			})
		}
		return nil

	}
	podIP := ""
	destinationName := destination.Name
	// if not take node with name and a random one
	for _, dsPod := range dsPodList.Items {
		if dsPod.Spec.NodeName == destination.Name {
			destination.Kind = "pod"
			destination.Name = dsPod.Name
			destination.Namespace = dsPod.Namespace
			podIP = dsPod.Status.PodIP
			break
		}
	}

	node := &corev1.Node{}
	var destinationAddress string
	err := r.client.Get(ctx, client.ObjectKey{Namespace: "", Name: destinationName}, node)
	for _, address := range node.Status.Addresses {
		if address.Type == "InternalIP" {
			destinationAddress = address.Address
			break
		}
	}

	for _, dsPod := range dsPodList.Items {
		if dsPod.Spec.NodeName != destination.Name {
			source.Container = "net-server"
			source.Name = dsPod.Name
			source.Namespace = dsPod.Namespace
			break
		}
	}
	qperfOut, err := Qperf(ctx, r.config, *source, *destination, protocol, podIP)
	if err != nil {
		status.QperfNodeEndpoints = append(status.QperfNodeEndpoints, v1alpha1.QperfNodeEndpoint{
			NodeParams: v1alpha1.NodeParams{
				Name: destinationName,
				IP:   destinationAddress,
			},
			NodeResults: v1alpha1.QperfResult{
				State: v1alpha1.FailedQperfCmd,
			},
		})
	} else {
		status.QperfNodeEndpoints = append(status.QperfNodeEndpoints, v1alpha1.QperfNodeEndpoint{
			NodeParams: v1alpha1.NodeParams{
				Name: destinationName,
				IP:   destinationAddress,
			},
			NodeResults: v1alpha1.QperfResult{
				State:    v1alpha1.SuccessQperfCmd,
				TCPDelay: qperfOut.tcpdelay,
				UDPDelay: qperfOut.udpdelay,
			},
		})
	}
	return nil
}

func (r *ReconcileNetworkDelayTest) reconcileNetworkDelay(ctx context.Context, networkDelayTest *v1alpha1.NetworkDelayTest) (reconcile.Result, error) {
	var (
		status = &v1alpha1.QperfStatus{
			TypeMeta: QperfStatusTypeMeta,
		}
	)

	for _, destination := range networkDelayTest.Spec.Destinations {
		switch destination.Kind {
		case v1alpha1.Pod:
			err := r.PodQperf(ctx, status, &networkDelayTest.Spec.Source, &destination, networkDelayTest.Spec.Protocol)
			if err != nil {
				return apimachinery.ReconcileErr(err)
			}
		case v1alpha1.Node:
			err := r.NodeQperf(ctx, status, &networkDelayTest.Spec.Source, &destination, networkDelayTest.Spec.Protocol)
			if err != nil {
				return apimachinery.ReconcileErr(err)
			}
		}
	}
	err := apimachinery.TryUpdateStatus(ctx, retry.DefaultBackoff, r.client, networkDelayTest, func() error {
		networkDelayTest.Status.TestStatus = &runtime.RawExtension{Object: status}
		return nil
	})
	if err != nil {
		return apimachinery.ReconcileErr(err)
	}

	return reconcile.Result{
		RequeueAfter: reconcilePeriod,
	}, nil
}
