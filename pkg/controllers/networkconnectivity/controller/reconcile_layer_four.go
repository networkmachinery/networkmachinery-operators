package controller

import (
	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"context"
)

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
