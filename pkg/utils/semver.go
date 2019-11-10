package utils

import (
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/pkg/errors"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
)

var Kubernetes116 = semver.New("1.16.0")

func ShouldUseEphemeralContainers(config *rest.Config) (bool, error) {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return false, errors.Wrap(err, "failed to initialize discovery client to get the cluster version, reason: ")
	}

	serverVersion, err := discoveryClient.ServerVersion()
	if err != nil {
		return false, errors.Wrap(err, "failed to get cluster server version, reason: ")
	}

	normalizedVersion := strings.Trim(serverVersion.GitVersion, "v")
	return !semver.New(normalizedVersion).LessThan(*Kubernetes116), nil
}
