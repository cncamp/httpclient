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
	"net/http"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/cncamp/httpclient/api/v1alpha1"
)

// HTTPClientReconciler reconciles a HTTPClient object
type HTTPClientReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.cncamp.io,resources=httpclients,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.cncamp.io,resources=httpclients/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.cncamp.io,resources=httpclients/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HTTPClient object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *HTTPClientReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	var httpclientholder appsv1alpha1.HTTPClient
	// your logic here
	if err := r.Get(ctx, req.NamespacedName, &httpclientholder); err != nil {
		log.Error(err, "Failed to read httpClient")
	}
	log.Info("Calling server,", "server address", httpclientholder.Spec.ServerAddress)

	if httpclientholder.Spec.ServerAddress != "" {
		totalCall := httpclientholder.Status.Failed + httpclientholder.Status.Succeeded
		if totalCall >= httpclientholder.Spec.Count {
			return ctrl.Result{}, nil
		}
		for i := 0; i < httpclientholder.Spec.Count; i++ {
			log.Info("Calling server,", "times", i)
			_, err := http.Get(httpclientholder.Spec.ServerAddress)
			if err != nil {
				log.Info("HTTP get failed with error: ", "error", err)
				httpclientholder.Status.Failed = httpclientholder.Status.Failed + 1
			} else {
				log.Info("HTTP get succeeded")
				httpclientholder.Status.Succeeded = httpclientholder.Status.Succeeded + 1
			}
		}
	}
	err := r.Client.Status().Update(context.TODO(), &httpclientholder, &client.UpdateOptions{})
	if err != nil {
		log.Error(err, "Failed to update httpclient")
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HTTPClientReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.HTTPClient{}).
		Complete(r)
}
