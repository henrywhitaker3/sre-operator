package handlers

import (
	"context"

	v1alpha1 "github.com/henrywhitaker3/sre-operator/api/v1alpha1"
	"github.com/henrywhitaker3/sre-operator/internal/store"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type ScriptHandler struct {
	ctx    context.Context
	client client.Client
	req    ctrl.Request
	obj    *v1alpha1.Script
	store  *store.Store
}

func NewScriptHandler(ctx context.Context, client client.Client, req ctrl.Request, store *store.Store) *ScriptHandler {
	return &ScriptHandler{
		ctx:    ctx,
		client: client,
		req:    req,
		store:  store,
	}
}

func (h *ScriptHandler) Get() error {
	d := &v1alpha1.Script{}
	if err := h.client.Get(h.ctx, h.req.NamespacedName, d); err != nil {
		return err
	}
	h.obj = d
	return nil
}

func (h *ScriptHandler) CreateOrUpdate() (error, bool) {
	// if len(h.obj.Spec.Triggers) == 0 {
	// 	return ErrNoTriggers, false
	// }

	// var subs store.StoreSubscriber

	// // TODO: implement it pls

	// for _, t := range h.obj.Spec.Triggers {
	// 	_, ok := h.store.Get(t)
	// 	if !ok {
	// 		return ErrUnknwonHook, true
	// 	}
	// 	h.store.StoreFunc(t, h.obj.Name, subs)
	// }
	return nil, true
}

func (h *ScriptHandler) Delete() error {
	// for _, t := range h.obj.Spec.Triggers {
	// 	_, ok := h.store.Get(t)
	// 	if !ok {
	// 		continue
	// 	}
	// 	h.store.DropFunc(t, h.obj.Name)
	// }
	return nil
}

func (h *ScriptHandler) DeletionTimestampIsZero() bool {
	return h.obj.DeletionTimestamp.IsZero()
}

func (h *ScriptHandler) AddFinalizer(fn string) error {
	controllerutil.AddFinalizer(h.obj, fn)
	return h.client.Update(h.ctx, h.obj)
}

func (h *ScriptHandler) RemoveFinalizer(fn string) error {
	controllerutil.RemoveFinalizer(h.obj, fn)
	return h.client.Update(h.ctx, h.obj)
}

func (h *ScriptHandler) SuccessStatus() error {
	h.obj.Status = v1alpha1.ScriptStatus{
		Registered: true,
	}
	return h.client.SubResource("status").Update(h.ctx, h.obj)
}

func (h *ScriptHandler) ErrorStatus(err error) error {
	h.obj.Status = v1alpha1.ScriptStatus{
		Registered: false,
		Error:      err.Error(),
	}
	return h.client.SubResource("status").Update(h.ctx, h.obj)
}

var _ Handler = &ScriptHandler{}
