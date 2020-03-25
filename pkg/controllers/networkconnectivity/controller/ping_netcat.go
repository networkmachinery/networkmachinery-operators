package controller

import (
	"bytes"
	"context"
	"fmt"

	"github.com/networkmachinery/networkmachinery-operators/pkg/utils/apimachinery"

	networkmachineryv1alpha1 "github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils/executor"
	"k8s.io/client-go/rest"
)

func Ping(ctx context.Context, config *rest.Config, source networkmachineryv1alpha1.NetworkSourceEndpoint, host string) (*PingOutput, error) {
	var (
		stdOut, stdErr bytes.Buffer
		execOpts       = executor.PodExecOptions{
			Namespace: source.Namespace,
			Name:      source.Name,
			Command:   fmt.Sprintf("ping -c3 %s", host),
			Container: source.Container,
			StandardCmdOpts: executor.StandardCmdOpts{
				StdErr: &stdErr,
				StdOut: &stdOut,
			},
		}
	)

	useEphemeralContainers, err := utils.ShouldUseEphemeralContainers(config)
	if err != nil {
		return nil, err
	}

	if useEphemeralContainers {
		debugContainerName := "net-debug"
		image := "nicolaka/netshoot"
		if err := apimachinery.CreateOrUpdateEphemeralContainer(config, source.Namespace, source.Name, debugContainerName, image); err != nil {
			return nil, err
		}
		execOpts.Container = debugContainerName

		if err := apimachinery.EphemeralContainerInStatus(ctx, config, source.Namespace, source.Name); err != nil {
			return nil, err
		}
	}

	err = utils.PodExec(ctx, config, execOpts)
	if err != nil {
		return &PingOutput{state: networkmachineryv1alpha1.FailedPing}, err
	}

	ping := &utils.Ping{}
	utils.ParsePingOutput(stdOut.Bytes(), ping)

	return &PingOutput{
		state: networkmachineryv1alpha1.SuccessPing,
		min:   ping.Min(),
		max:   ping.Max(),
		avg:   ping.Average(),
	}, nil
}

func NetCat(ctx context.Context, config *rest.Config, source networkmachineryv1alpha1.NetworkSourceEndpoint, host, port string) (*NetcatOutput, error) {
	var (
		stdOut, stdErr bytes.Buffer
		execOpts       = executor.PodExecOptions{
			Namespace: source.Namespace,
			Name:      source.Name,
			Command:   fmt.Sprintf("nc -z -v %s %s 2>&1", host, port),
			Container: source.Container,
			StandardCmdOpts: executor.StandardCmdOpts{
				StdErr: &stdErr,
				StdOut: &stdOut,
			},
		}
	)

	useEphemeralContainers, err := utils.ShouldUseEphemeralContainers(config)
	if err != nil {
		return nil, err
	}

	if useEphemeralContainers {
		debugContainerName := "net-debug"
		image := "nicolaka/netshoot"
		if err := apimachinery.CreateOrUpdateEphemeralContainer(config, source.Namespace, source.Name, debugContainerName, image); err != nil {
			return nil, err
		}
		execOpts.Container = debugContainerName

		if err := apimachinery.EphemeralContainerInStatus(ctx, config, source.Namespace, source.Name); err != nil {
			return nil, err
		}
	}

	err = utils.PodExec(ctx, config, execOpts)
	if err != nil {
		return &NetcatOutput{state: networkmachineryv1alpha1.Refused}, err
	}

	netcat := &utils.Netcat{}
	utils.ParseNetcatOutput(stdOut.String(), netcat)

	return &NetcatOutput{
		state: netcat.State(),
	}, nil
}
