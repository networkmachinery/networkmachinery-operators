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

func Qperf(ctx context.Context, config *rest.Config, source networkmachineryv1alpha1.NetworkSourceEndpoint, destination networkmachineryv1alpha1.NetworkDestinationEndpoint, protocol, podIP string) (*QperfOutput, error) {
	var prot networkmachineryv1alpha1.Protocol
	if protocol == "tcp" {
		prot = networkmachineryv1alpha1.TCP
	}
	if protocol == "udp" {
		prot = networkmachineryv1alpha1.UDP
	}
	var (
		stdOut, stdErr bytes.Buffer
		execOptsServer = executor.PodExecOptions{
			Namespace: destination.Namespace,
			Name:      destination.Name,
			Command:   "",
			Container: "net-server",
			StandardCmdOpts: executor.StandardCmdOpts{
				StdErr: &stdErr,
				StdOut: &stdOut,
			},
		}
		execOpts = executor.PodExecOptions{
			Namespace: source.Namespace,
			Name:      source.Name,
			Command:   fmt.Sprintf("qperf -vvs %s %s", podIP, prot),
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
		image := "docktofuture/ubuntuqperfserver:latest"
		serverContainerName := "net-server"
		// run server container
		if err := apimachinery.CreateOrUpdateEphemeralContainer(config, destination.Namespace, destination.Name, serverContainerName, image); err != nil {
			return nil, err
		}
		execOptsServer.Container = serverContainerName

		if err := apimachinery.EphemeralContainerInStatus(ctx, config, destination.Namespace, destination.Name); err != nil {
			return nil, err
		}
		// run client container
		debugContainerName := "net-server"
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
		return &QperfOutput{state: networkmachineryv1alpha1.FailedQperfCmd}, err
	}

	qperf := &utils.Qperf{}
	utils.ParseQperfOutput(stdOut.Bytes(), qperf, prot)

	var output QperfOutput
	if protocol == "tcp" {
		output = QperfOutput{
			state:    networkmachineryv1alpha1.SuccessQperfCmd,
			tcpdelay: qperf.TCPDelay(),
		}
	}
	if protocol == "udp" {
		output = QperfOutput{
			state:    networkmachineryv1alpha1.SuccessQperfCmd,
			udpdelay: qperf.UDPDelay(),
		}
	}
	return &output, nil
}
