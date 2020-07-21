package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	sidecarv1alpha1 "github.com/bvwells/sidecar-operator/api/v1alpha1"
)

const (
	sidecarContainerName     = "sidecar-container"
	sidecarOperatorFinalizer = "finalizer.bvwells.github.com"
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
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile
			// request.
			// Owned objects are automatically garbage collected. For additional
			// cleanup logic use finalizers. Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Check if the SidecarOperator instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	if sidecarOperator.DeletionTimestamp != nil {
		if contains(sidecarOperator.GetFinalizers(), sidecarOperatorFinalizer) {
			// Run finalization logic for sidecarOperatorFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeSidecarOperator(ctx, logger, sidecarOperator); err != nil {
				return ctrl.Result{}, err
			}

			// Remove sidecarOperatorFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			controllerutil.RemoveFinalizer(sidecarOperator, sidecarOperatorFinalizer)
			err := r.Update(ctx, sidecarOperator)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer for this CR
	if !contains(sidecarOperator.GetFinalizers(), sidecarOperatorFinalizer) {
		if err := r.addFinalizer(ctx, logger, sidecarOperator); err != nil {
			return ctrl.Result{}, err
		}
	}

	image := sidecarOperator.Spec.Image

	logger.Info(fmt.Sprintf("deploying sidecar '%s'", image))

	deployments := &appsv1.DeploymentList{}
	if err := r.List(ctx, deployments); err != nil {
		return ctrl.Result{}, err
	}

	for _, deployment := range deployments.Items {
		injectSidecar := true
		for i, container := range deployment.Spec.Template.Spec.Containers {
			// Deployment has existing sidecar
			if container.Name == sidecarContainerName {
				injectSidecar = false
				// Sidecar image does not match
				if container.Image != image {
					logger.Info(fmt.Sprintf("upgrading deployment '%s' to use image '%s'", deployment.Name, image))
					deployment.Spec.Template.Spec.Containers[i] = newSidecarContainer(image)
				} else {
					continue
				}
			}
		}
		if injectSidecar {
			logger.Info(fmt.Sprintf("adding sidecar to deployment '%s'", deployment.Name))
			container := newSidecarContainer(image)
			deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, container)
		}

		err = r.Update(ctx, &deployment, &client.UpdateOptions{})
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	// TODO - watch for new deployments.

	return ctrl.Result{}, nil
}

func (r *SidecarOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&sidecarv1alpha1.SidecarOperator{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}

func newSidecarContainer(image string) corev1.Container {
	return corev1.Container{
		Name:    sidecarContainerName,
		Image:   image,
		Command: []string{"sleep", "3600"},
	}
}

func (r *SidecarOperatorReconciler) finalizeSidecarOperator(ctx context.Context,
	logger logr.Logger, s *sidecarv1alpha1.SidecarOperator) error {
	deployments := &appsv1.DeploymentList{}
	if err := r.List(ctx, deployments); err != nil {
		return err
	}

	for _, deployment := range deployments.Items {
		deleteSidecar := false
		n := 0
		for _, container := range deployment.Spec.Template.Spec.Containers {
			if container.Name != sidecarContainerName {
				deployment.Spec.Template.Spec.Containers[n] = container
				n++
			} else {
				deleteSidecar = true
			}
		}
		deployment.Spec.Template.Spec.Containers = deployment.Spec.Template.Spec.Containers[:n]

		if deleteSidecar {
			logger.Info(fmt.Sprintf("removing sidecar from deployment '%s'", deployment.Name))
			err := r.Update(ctx, &deployment, &client.UpdateOptions{})
			if err != nil {
				return err
			}
		}
	}

	logger.Info("Successfully finalized sidecar operator")

	return nil
}

func (r *SidecarOperatorReconciler) addFinalizer(ctx context.Context,
	logger logr.Logger, s *sidecarv1alpha1.SidecarOperator) error {
	logger.Info("Adding Finalizer for the sidecar operator")

	controllerutil.AddFinalizer(s, sidecarOperatorFinalizer)

	// Update CR
	err := r.Update(ctx, s)
	if err != nil {
		logger.Error(err, "Failed to update SidecarOperator with finalizer")
		return err
	}
	return nil
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
