package utils

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

// SetupSignalHandlerContext sets up a context from signals.SetupSignalHandler stop channel.
func SetupSignalHandlerContext() context.Context {
	return ContextFromStopChannel(signals.SetupSignalHandler())
}
