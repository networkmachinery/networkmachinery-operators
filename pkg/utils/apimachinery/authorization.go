package apimachinery

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	authorizationv1 "k8s.io/api/authorization/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Can checks if current user has permissions
func Can(ctx context.Context, client client.Client, resourceAttr *authorizationv1.ResourceAttributes) error {
	selfSubjectAccessReview := &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: resourceAttr,
		},
	}

	err := client.Create(ctx, selfSubjectAccessReview)
	if err != nil {
		return errors.Wrap(err, "Failed to create SelfSubjectAccessReview")
	}

	if !selfSubjectAccessReview.Status.Allowed {
		denyReason := fmt.Sprintf("Current user has no permission to %s %s/%s subresource in namespace:%s. Detail:", resourceAttr.Verb, resourceAttr.Resource, resourceAttr.Subresource, resourceAttr.Namespace)
		if len(selfSubjectAccessReview.Status.Reason) > 0 {
			denyReason = fmt.Sprintf("%s %v, ", denyReason, selfSubjectAccessReview.Status.Reason)
		}
		if len(selfSubjectAccessReview.Status.EvaluationError) > 0 {
			denyReason = fmt.Sprintf("%s %v", denyReason, selfSubjectAccessReview.Status.EvaluationError)
		}
		return fmt.Errorf(denyReason)
	}
	return nil
}
