package controllers

import (
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkconnectivity/controller"
)

func init() {
	AddToManagerFuncs = append(AddToManagerFuncs, controller.Add)
}
