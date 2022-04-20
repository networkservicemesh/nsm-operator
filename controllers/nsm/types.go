package controllers

import (
	"context"

	nsmv1alpha1 "github.com/networkservicemesh/nsm-operator/apis/nsm/v1alpha1"
)

type Reconciler interface {
	Reconcile(ctx context.Context, nsm *nsmv1alpha1.NSM) error
}
