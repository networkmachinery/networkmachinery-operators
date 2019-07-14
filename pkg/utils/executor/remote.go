package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// NewDebugExecutor returns a debugExecutor
func NewDebugExecutor(config *rest.Config) DebugExecutor {
	return &debugExecutor{
		config: config,
	}
}

type DebugExecOptions struct {
	ContainerID  string
	Command      string
	CommandSlice []string
	Image        string

	TargetHost string
	TargetPort int

	StandardCmdOpts
}

type debugExecutor struct {
	config *rest.Config
}

// Execute executes a command on a pod
func (d *debugExecutor) Execute(ctx context.Context, options DebugExecOptions) error {
	uri, err := url.Parse(fmt.Sprintf("http://%s:%d", options.TargetHost, options.TargetPort))
	if err != nil {
		return err
	}
	uri.Path = fmt.Sprintf("/api/v1/debug")
	bytes, err := json.Marshal(options.CommandSlice)
	if err != nil {
		return err
	}
	params := url.Values{}
	params.Add("image", options.Image)
	params.Add("container", options.ContainerID)
	params.Add("stdin", "true")
	params.Add("stdout", "true")
	params.Add("stderr", "true")
	params.Add("command", string(bytes))
	uri.RawQuery = params.Encode()

	executor, err := remotecommand.NewSPDYExecutor(d.config, http.MethodPost, uri)
	if err != nil {
		return fmt.Errorf("failed to initialized the command exector: %v", err)
	}

	return executor.Stream(remotecommand.StreamOptions{
		Stdin:  strings.NewReader(options.Command),
		Stdout: options.StdOut,
		Stderr: options.StdErr,
		Tty:    false,
	})
}
