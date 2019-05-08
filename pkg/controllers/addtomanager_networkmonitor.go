package controllers

import (
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkmonitor/controller"
)

func init() {
	AddToManagerFuncs = append(AddToManagerFuncs, controller.Add)
}
