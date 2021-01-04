package controller

import (
	"github.com/angudadevops/service-deployment-operator/pkg/controller/sd"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, sd.Add)
}
