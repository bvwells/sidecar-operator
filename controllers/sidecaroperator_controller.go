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
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	sidecarv1alpha1 "github.com/bvwells/sidecar-operator/api/v1alpha1"
)

const sidecarContainerName = "sidecar-container"

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
	// TODO - set controller reference for newly created resources.

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
