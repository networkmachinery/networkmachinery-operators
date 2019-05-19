package controllers

import (
	"github.com/networkmachinery/networkmachinery-operators/pkg/controllers/networktrafficshaper/controller"
)

func init() {
	AddToManagerFuncs = append(AddToManagerFuncs, controller.Add)
}
