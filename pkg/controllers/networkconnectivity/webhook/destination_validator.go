/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhook

import (
	"context"
	"net/http"

	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/validate-destination-v1alpha1-networkconnectivitytest,mutating=false,failurePolicy=fail,groups="networkmachinery.io",resources=networkconnectivitytests,verbs=create;update,versions=v1alpha1,name=networkconnectivitytest.networkmachinery.io

// DestinationValidator validates the destination
type DestinationValidator struct {
	client  client.Client
	decoder *admission.Decoder
}

// InjectClient injects the client.
func (v *DestinationValidator) InjectClient(c client.Client) error {
	v.client = c
	return nil
}

// InjectDecoder injects the decoder.
func (v *DestinationValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}

// DestinationValidator makes sure that endpoint specs conform with the layer number set (e.g., a layer 3 endpoint can not have a port set)
func (v *DestinationValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	networkConnectivityTest := &v1alpha1.NetworkConnectivityTest{}

	err := v.decoder.Decode(req, networkConnectivityTest)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	allowed, reason, err := v.validateDestinationsFn(ctx, networkConnectivityTest)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}
	return admission.ValidationResponse(allowed, reason)
}

func (v *DestinationValidator) validateDestinationsFn(ctx context.Context, nct *v1alpha1.NetworkConnectivityTest) (bool, string, error) {
	for _, destination := range nct.Spec.Destinations {
		switch destination.Kind {
		case v1alpha1.IP:
			if len(destination.IP) == 0 {
				return false, "A destination IP endpoint needs to have an IP", nil
			}
		case v1alpha1.Pod, v1alpha1.Service:
			if len(destination.Name) == 0 || len(destination.Namespace) == 0 {
				return false, "Endpoint needs the namespace and name specified", nil
			}
		}
	}
	return true, "", nil
}
