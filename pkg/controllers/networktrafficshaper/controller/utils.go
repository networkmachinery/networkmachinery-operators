package controller

import (
	"bytes"
	"context"
	"fmt"

	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	"k8s.io/client-go/rest"
)

func shapeTraffic(ctx context.Context, config *rest.Config, namespace, name, device, value, shapeType string) error {
	command := fmt.Sprintf("tc qdisc add dev %s root netem %s %s", device, shapeType, value)
	return shape(ctx, config, namespace, name, command)
}
func undoShape(ctx context.Context, config *rest.Config, namespace, name, device string) error {
	command := fmt.Sprintf("tc qdisc del dev %s root", device)
	return shape(ctx, config, namespace, name, command)
}

func shape(ctx context.Context, config *rest.Config, namespace, name, command string) error {
	var stdOut, stdErr bytes.Buffer
	execOpts := utils.ExecOptions{
		Namespace: namespace,
		Name:      name,
		Command:   command,
		Container: "", //TODO Fixme, get the right container value
		StandardCmdOpts: utils.StandardCmdOpts{
			StdErr: &stdErr,
			StdOut: &stdOut,
		},
	}

	// TODO: handle error if tc config already Exists
	_ = utils.PodExec(ctx, config, execOpts)
	return nil
}
