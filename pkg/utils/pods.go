// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PodExec(ctx context.Context, config *rest.Config, options ExecOptions) error {
	return NewPodExecutor(config).Execute(ctx, options)
}

// getFirstRunningPodWithLabels fetches the first running pod with the desired set of labels <labelsMap>
func getFirstRunningPodWithLabels(ctx context.Context, labelsMap labels.Selector, namespace string, client client.Client) (*corev1.Pod, error) {
	var (
		podList *corev1.PodList
		err     error
	)
	podList, err = GetPodsByLabels(ctx, client, labelsMap, namespace)
	if err != nil {
		return nil, err
	}
	if len(podList.Items) == 0 {
		return nil, fmt.Errorf("no running pods found")
	}

	for _, pod := range podList.Items {
		if pod.Status.Phase == corev1.PodRunning {
			return &pod, nil
		}
	}

	return nil, fmt.Errorf("no running pods found")
}

func GetPodsByLabels(ctx context.Context, c client.Client, labelsMap labels.Selector, namespace string) (*corev1.PodList, error) {
	podList := &corev1.PodList{}
	err := c.List(ctx, podList, ToListOptionFunc(&client.ListOptions{
		Namespace:     namespace,
		LabelSelector: labelsMap,
	}))
	if err != nil {
		return nil, err
	}
	return podList, nil
}
