/*
Copyright 2024.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1alpha1 "github.com/henrywhitaker3/sre-operator/api/v1alpha1"
	"github.com/henrywhitaker3/sre-operator/internal/app"
	"github.com/henrywhitaker3/sre-operator/internal/controller/handlers"
)

// RolloutReconciler reconciles a Rollout object
type RolloutReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	app    *app.App
}

//+kubebuilder:rbac:groups=sre.henrywhitaker.com,resources=rollouts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=sre.henrywhitaker.com,resources=rollouts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=sre.henrywhitaker.com,resources=rollouts/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps;extensions,resources=deployments;pods;statefulsets;daemonsets,verbs=get;patch

func (r *RolloutReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	h := handlers.NewRolloutHandler(ctx, r.Client, req, r.app.HookStore)
	return handlers.RunHandler(l, h)
}

// SetupWithManager sets up the controller with the Manager.
func (r *RolloutReconciler) SetupWithManager(mgr ctrl.Manager, app *app.App) error {
	r.app = app
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Rollout{}).
		Complete(r)
}
