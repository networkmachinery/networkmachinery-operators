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
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/go-logr/logr"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	waitFn func(interval time.Duration, condition wait.ConditionFunc, stopCh <-chan struct{}) error
	readFn func(obj runtime.Object) bool
)

// CreateOrUpdate creates or updates the object. Optionally, it executes a transformation function before the
// request is made.
func CreateOrUpdate(ctx context.Context, c client.Client, obj runtime.Object, transform func() error) error {
	key, err := client.ObjectKeyFromObject(obj)
	if err != nil {
		return err
	}

	if err := c.Get(ctx, key, obj); err != nil {
		if apierrors.IsNotFound(err) {
			if transform != nil && transform() != nil {
				return err
			}
			return c.Create(ctx, obj)
		}
		return err
	}

	if transform != nil && transform() != nil {
		return err
	}
	return c.Update(ctx, obj)
}

func CreateIfNotFound(ctx context.Context, c client.Client, logger logr.Logger, waitFn waitFn, readyFn readFn, obj runtime.Object) error {
	key, err := client.ObjectKeyFromObject(obj)
	if err != nil {
		return err
	}

	if err := c.Get(ctx, key, obj); err != nil {
		if apierrors.IsNotFound(err) {
			if err := c.Create(ctx, obj); err != nil {
				return err
			}
			err := waitFn(5*time.Second, func() (bool, error) {
				logger.Info("Waiting for agent obj to be ready...!!")
				err := c.Get(ctx, key, obj)
				if err != nil {
					return false, err
				}
				return readyFn(obj), nil
			}, ctx.Done())
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}
