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
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	configurationv1alpha1 "github.com/henrywhitaker3/sre-operator/api/v1alpha1"
	"github.com/henrywhitaker3/sre-operator/internal/app"
	"github.com/henrywhitaker3/sre-operator/internal/controller/handlers"
)

// WebhookReconciler reconciles a Webhook object
type WebhookReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	app    *app.App
}

//+kubebuilder:rbac:groups=sre.henrywhitaker.com,resources=webhooks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=sre.henrywhitaker.com,resources=webhooks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=sre.henrywhitaker.com,resources=webhooks/finalizers,verbs=update

func (r *WebhookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	h := handlers.NewWebhookHandler(ctx, r.Client, req, r.app.HookStore)
	res, err := handlers.RunHandler(l, h)
	if err != nil {
		l.Error(err, "webhook handler error")
	}
	return res, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebhookReconciler) SetupWithManager(mgr ctrl.Manager, app *app.App) error {
	r.app = app
	return ctrl.NewControllerManagedBy(mgr).
		For(&configurationv1alpha1.Webhook{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}
