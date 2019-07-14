package executor

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// NewPodExecutor returns a podExecutor
func NewPodExecutor(config *rest.Config) PodExecutor {
	return &podExecutor{
		config: config,
	}
}

type podExecutor struct {
	config *rest.Config
}

type PodExecOptions struct {
	Namespace string
	Name      string
	Container string
	Command   string

	StandardCmdOpts
}

// Execute executes a command on a pod
func (p *podExecutor) Execute(ctx context.Context, options PodExecOptions) error {
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
		return fmt.Errorf("failed to initialize the debug executor: %v", err)
	}

	return executor.Stream(remotecommand.StreamOptions{
		Stdin:  strings.NewReader(options.Command),
		Stdout: options.StdOut,
		Stderr: options.StdErr,
		Tty:    false,
	})
}
