package main

import (
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkdelay/cmd/app"
	"github.com/networkmachinery/networkmachinery-operators/pkg/utils"

	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func main() {
	log.SetLogger(log.ZapLogger(false))
	cmd := app.NewNetworkDelayTestCmd(utils.SetupSignalHandlerContext())

	if err := cmd.Execute(); err != nil {
		utils.LogErrAndExit(err, "error executing the main controller command")
	}
}
