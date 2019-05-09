package apimachinery

import (
	apimachinery "github.com/networkmachinery/networkmachinery-operators/pkg/utils/apimachinery/error"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcileErr returns a reconcile.Result or an error, depending on whether the error is a
// RequeueAfterError or not.
func ReconcileErr(err error) (reconcile.Result, error) {
	if requeueAfter, ok := err.(*apimachinery.RequeueAfterError); ok {
		return reconcile.Result{Requeue: true, RequeueAfter: requeueAfter.RequeueAfter}, nil
	}
	return reconcile.Result{}, err
}
