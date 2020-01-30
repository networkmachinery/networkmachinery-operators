package controllers

import (
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networkdelay/controller"
)

func init() {
	AddToManagerFuncs = append(AddToManagerFuncs, controller.Add)
}
