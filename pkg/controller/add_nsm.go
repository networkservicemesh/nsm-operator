package controller

import (
	"github.com/networkservicemesh/nsm-operator/pkg/controller/nsm"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, nsm.Add)
}
