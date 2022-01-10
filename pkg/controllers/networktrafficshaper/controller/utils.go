package controller

import (
	"bytes"
	"context"
	"fmt"

	"github.com/networkmachinery/networkmachinery-operators/pkg/utils/apimachinery"

	"github.com/networkmachinery/networkmachinery-operators/pkg/utils/executor"

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
	execOpts := executor.PodExecOptions{
		Namespace: namespace,
		Name:      name,
		Command:   command,
		Container: "", //TODO Fixme, get the right container value
		StandardCmdOpts: executor.StandardCmdOpts{
			StdErr: &stdErr,
			StdOut: &stdOut,
		},
	}

	useEphemeralContainers, err := utils.ShouldUseEphemeralContainers(config)
	if err != nil {
		return err
	}

	if useEphemeralContainers {
		debugContainerName := "net-debug"
		image := "nicolaka/netshoot"
		if err := apimachinery.CreateOrUpdateEphemeralContainer(config, namespace, name, debugContainerName, image); err != nil {
			return err
		}
		execOpts.Container = debugContainerName

		if err := apimachinery.EphemeralContainerInStatus(ctx, config, namespace, name); err != nil {
			return err
		}
	}

	// TODO: handle error if tc config already Exists
	if err = utils.PodExec(ctx, config, execOpts); err != nil {
		return err
	}

	return nil
}
