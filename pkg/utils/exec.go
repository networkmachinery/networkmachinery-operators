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
	"bytes"
	"fmt"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"strings"

	"context"
)

// NewPodExecutor returns a podExecutor
func NewPodExecutor(config *rest.Config) PodExecutor {
	return &podExecutor{
		config: config,
	}
}

// PodExecutor is the pod executor interface
type PodExecutor interface {
	Execute(ctx context.Context, options ExecOptions) error
}

type podExecutor struct {
	config *rest.Config
}

type ExecOptions struct {
	Namespace string
	Name string
	Container string
	Command string

	StandardCmdOpts
}

type StandardCmdOpts struct {
	StdOut, StdErr  *bytes.Buffer
}


// Execute executes a command on a pod
func (p *podExecutor) Execute(ctx context.Context, options ExecOptions) error {
	client, err := corev1client.NewForConfig(p.config)
	if err != nil {
		return err
	}

	request := client.RESTClient().
		Post().
		Resource("pods").
		Name(options.Name).
		Namespace(options.Namespace).
		SubResource("exec").
		Param("container", options.Container).
		Param("command", "/bin/sh").
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "false").
		Context(ctx)

	executor, err := remotecommand.NewSPDYExecutor(p.config, http.MethodPost, request.URL())
	if err != nil {
		return fmt.Errorf("failed to initialized the command exector: %v", err)
	}

	return  executor.Stream(remotecommand.StreamOptions{
		Stdin:  strings.NewReader(options.Command),
		Stdout: options.StdOut,
		Stderr: options.StdErr,
		Tty:    false,
	})
}
