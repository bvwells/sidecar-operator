package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	sidecarv1alpha1 "github.com/bvwells/sidecar-operator/api/v1alpha1"
)

// SidecarOperatorReconciler reconciles a SidecarOperator object
type SidecarOperatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=sidecar.bvwells.github.com,resources=sidecaroperators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=sidecar.bvwells.github.com,resources=sidecaroperators/status,verbs=get;update;patch

func (r *SidecarOperatorReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("sidecaroperator", req.NamespacedName)

	sidecarOperator := &sidecarv1alpha1.SidecarOperator{}
	err := r.Get(ctx, req.NamespacedName, sidecarOperator)
	if err != nil {
		return ctrl.Result{}, err
	}

	logger.Info(fmt.Sprintf("deploying sidecar '%s'", sidecarOperator.Spec.Image))

	return ctrl.Result{}, nil
}

func (r *SidecarOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&sidecarv1alpha1.SidecarOperator{}).
		Complete(r)
}
