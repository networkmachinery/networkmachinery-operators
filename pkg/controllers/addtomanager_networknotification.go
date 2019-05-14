package controllers

import (
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkcontrol/controller"
)

func init() {
	AddToManagerFuncs = append(AddToManagerFuncs, controller.Add)
}
