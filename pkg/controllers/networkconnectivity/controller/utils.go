package controller

import (
	"bytes"
	"fmt"
	"github.com/networkmachinery/networkmachinery-operators/pkg/apis/networkmachinery/v1alpha1"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"
	"k8s.io/client-go/rest"
	"context"
)

func Ping(ctx context.Context, config *rest.Config, source v1alpha1.NetworkSourceEndpoint, host string) (*PingOutput, error) {
	var stdOut, stdErr bytes.Buffer
	execOpts := utils.ExecOptions{
		Namespace: source.Namespace,
		Name: source.Name,
		Command: fmt.Sprintf("ping -c3 %s", host),
		Container: source.Container,
		StandardCmdOpts: utils.StandardCmdOpts{
			StdErr: &stdErr,
			StdOut: &stdOut,
		},
	}

	err := utils.PodExec(ctx, config,execOpts)
	if err != nil {
		return &PingOutput{state: v1alpha1.FailedPing}, err
	}

	ping := &utils.Ping{}
	utils.ParsePingOutput(stdOut.Bytes(), ping)

	return &PingOutput{
		state: v1alpha1.SuccessPing,
		min: ping.Min(),
		max: ping.Max(),
		avg: ping.Average(),
	}, nil
}
