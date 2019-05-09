// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package apimachinery

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
)

// HasFinalizer checks if the given object has a finalizer with the given name.
func HasFinalizer(obj runtime.Object, finalizerName string) (bool, error) {
	finalizers, _, err := finalizersAndAccessorOf(obj)
	if err != nil {
		return false, err
	}

	return finalizers.Has(finalizerName), nil
}

// EnsureFinalizer ensures that a finalizer of the given name is set on the given object.
// If the finalizer is not set, it adds it to the list of finalizers and updates the remote object.
func EnsureFinalizer(ctx context.Context, client client.Client, finalizerName string, obj runtime.Object) error {
	finalizers, accessor, err := finalizersAndAccessorOf(obj)
	if err != nil {
		return err
	}

	if finalizers.Has(finalizerName) {
		return nil
	}

	finalizers.Insert(finalizerName)
	accessor.SetFinalizers(finalizers.UnsortedList())

	return client.Update(ctx, obj)
}

func finalizersAndAccessorOf(obj runtime.Object) (sets.String, metav1.Object, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, nil, err
	}

	return sets.NewString(accessor.GetFinalizers()...), accessor, nil
}

// DeleteFinalizer ensures that the given finalizer is not present anymore in the given object.
// If it is set, it removes it and issues an update.
func DeleteFinalizer(ctx context.Context, client client.Client, finalizerName string, obj runtime.Object) error {
	finalizers, accessor, err := finalizersAndAccessorOf(obj)
	if err != nil {
		return err
	}

	if !finalizers.Has(finalizerName) {
		return nil
	}

	finalizers.Delete(finalizerName)
	accessor.SetFinalizers(finalizers.UnsortedList())

	return client.Update(ctx, obj)
}
