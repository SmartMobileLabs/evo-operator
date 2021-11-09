/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	app "smartmobilelabs.com/evo/evo-operator/api/v1alpha1"
)

var log = logf.Log.WithName("controller_sml-evo")

// SmlEvoReconciler reconciles a SmlEvo object
type SmlEvoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ops.dac.nokia.com,resources=licenceexpireds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.dac.nokia.com.example.com,resources=smlevoes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.dac.nokia.com.example.com,resources=smlevoes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.dac.nokia.com.example.com,resources=smlevoes/finalizers,verbs=update
//+kubebuilder:rbac:groups="";apps,resources=secrets;configmaps;statefulsets;pods;services;endpoints;persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SmlEvo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *SmlEvoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	reqLogger.Info("Reconciling SmlEvo")

	// Fetch the Consul instance
	instance := &app.SmlEvo{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	return r.handleCrChange(instance, req.Namespace)
}

// predicate filter added to extend for future extensions
func smlEventPredicate() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(event.CreateEvent) bool {
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Ignore updates to CR status in which case metadata.Generation does not change
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			// Evaluates to false if the object has been confirmed deleted.
			return !e.DeleteStateUnknown
		},
		GenericFunc: func(event.GenericEvent) bool {
			return true
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *SmlEvoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&app.SmlEvo{}).
		//		WithEventFilter(smlEventPredicate()).
		Complete(r)
}
