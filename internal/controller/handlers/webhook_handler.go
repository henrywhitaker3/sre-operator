package handlers

import (
	"context"

	"github.com/google/uuid"
	configurationv1alpha1 "github.com/henrywhitaker3/sre-operator/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type WebhookHandler struct {
	ctx    context.Context
	client client.Client
	req    ctrl.Request
	obj    *configurationv1alpha1.Webhook

	// The generated id
	id string
}

func NewWebhookHandler(ctx context.Context, client client.Client, req ctrl.Request) *WebhookHandler {
	return &WebhookHandler{
		ctx:    ctx,
		client: client,
		req:    req,
	}
}

func (h *WebhookHandler) Get() error {
	d := &configurationv1alpha1.Webhook{}
	if err := h.client.Get(h.ctx, h.req.NamespacedName, d); err != nil {
		return err
	}
	h.obj = d
	return nil
}

func (h *WebhookHandler) CreateOrUpdate() error {
	h.id = uuid.NewString()
	// TODO: add webhook the echo configuration
	if !controllerutil.ContainsFinalizer(h.obj, fn) {
		controllerutil.AddFinalizer(h.obj, fn)
		if err := h.client.Update(h.ctx, h.obj); err != nil {
			return err
		}
	}
	return nil
}

func (h *WebhookHandler) Delete() error {
	// TODO: remove echo config
	return nil
}

func (h *WebhookHandler) DeletionTimestampIsZero() bool {
	return h.obj.DeletionTimestamp.IsZero()
}

func (h *WebhookHandler) GetFinalizers() []string {
	return h.obj.Finalizers
}

func (h *WebhookHandler) RemoveFinalizer(fn string) error {
	controllerutil.RemoveFinalizer(h.obj, fn)
	return h.client.Update(h.ctx, h.obj)
}

func (h *WebhookHandler) SuccessStatus() error {
	h.obj.Status = configurationv1alpha1.WebhookStatus{
		Valid: true,
		ID:    h.id,
	}
	return h.client.SubResource("status").Update(h.ctx, h.obj)
}

func (h *WebhookHandler) ErrorStatus(err error) error {
	h.obj.Status = configurationv1alpha1.WebhookStatus{
		Valid: false,
		ID:    h.id,
	}
	return h.client.SubResource("status").Update(h.ctx, h.obj)
}

var _ Handler = &WebhookHandler{}
