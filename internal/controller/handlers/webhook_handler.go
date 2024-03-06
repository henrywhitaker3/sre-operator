package handlers

import (
	"context"
	"errors"

	configurationv1alpha1 "github.com/henrywhitaker3/sre-operator/api/v1alpha1"
	"github.com/henrywhitaker3/sre-operator/internal/metrics"
	"github.com/henrywhitaker3/sre-operator/internal/store"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var (
	ErrMissingID = errors.New("id is missing from spec field")
)

type WebhookHandler struct {
	ctx    context.Context
	client client.Client
	req    ctrl.Request
	obj    *configurationv1alpha1.Webhook

	store   *store.Store
	metrics *metrics.Metrics

	id string
}

func NewWebhookHandler(ctx context.Context, client client.Client, req ctrl.Request, store *store.Store, metrics *metrics.Metrics) *WebhookHandler {
	return &WebhookHandler{
		ctx:     ctx,
		client:  client,
		req:     req,
		store:   store,
		metrics: metrics,
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

func (h *WebhookHandler) CreateOrUpdate() (error, bool) {
	if h.obj.Spec.ID == "" {
		return ErrMissingID, false
	}
	h.id = h.obj.Spec.ID

	if _, err := h.store.Get(store.WEBHOOK, h.id); err != nil {
		if errors.Is(err, store.ErrUnknownTrigger) {
			h.metrics.WebhooksRegistered.Inc()
		} else {
			return err, false
		}
	}

	return h.store.Register(store.WEBHOOK, h.id), true
}

func (h *WebhookHandler) Delete() error {
	return h.store.Drop(store.WEBHOOK, h.id)
}

func (h *WebhookHandler) DeletionTimestampIsZero() bool {
	return h.obj.DeletionTimestamp.IsZero()
}

func (h *WebhookHandler) AddFinalizer(fn string) error {
	controllerutil.AddFinalizer(h.obj, fn)
	return h.client.Update(h.ctx, h.obj)
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
		Error: err.Error(),
	}
	return h.client.SubResource("status").Update(h.ctx, h.obj)
}

var _ Handler = &WebhookHandler{}
