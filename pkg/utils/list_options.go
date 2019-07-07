package utils

import "sigs.k8s.io/controller-runtime/pkg/client"

func ToListOptionFunc(fromListOptions *client.ListOptions) client.ListOptionFunc {
	return func(toListOptions *client.ListOptions) {
		*toListOptions = *fromListOptions
	}
}
