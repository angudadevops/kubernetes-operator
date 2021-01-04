package controller

import (
	"github.com/martin-helmich/kubernetes-operator-example/pkg/controller/sd"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, sd.Add)
}
