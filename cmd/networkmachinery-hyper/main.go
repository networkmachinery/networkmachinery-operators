package main

import (
	"os"

	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"

	"github.com/networkmachinery/networkmachinery-operators/cmd/networkmachinery-hyper/app"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

var (
	Log = log.Log
)

// LogErrAndExit logs the given error with msg and keysAndValues and calls `os.Exit(1)`.
func LogErrAndExit(err error, msg string, keysAndValues ...interface{}) {
	Log.Error(err, msg, keysAndValues...)
	os.Exit(1)
}

func main() {
	log.SetLogger(log.ZapLogger(false))
	cmd := app.NewHyperCommand(utils.ContextFromStopChannel(signals.SetupSignalHandler()))

	if err := cmd.Execute(); err != nil {
		LogErrAndExit(err, "error executing the main command")
	}
}
