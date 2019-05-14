package main

import (
	"github.com/gardener/gardener-extensions/pkg/controller"
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkconnectivity/cmd/app"

	controllercmd "github.com/gardener/gardener-extensions/pkg/controller/cmd"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func main() {
	log.SetLogger(log.ZapLogger(false))
	cmd := app.NewNetworkConnectivityTestCmd(controller.SetupSignalHandlerContext())

	if err := cmd.Execute(); err != nil {
		controllercmd.LogErrAndExit(err, "error executing the main controller command")
	}
}
