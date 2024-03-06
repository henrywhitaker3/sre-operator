package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/henrywhitaker3/flow"
	v1alpha1 "github.com/henrywhitaker3/sre-operator/api/v1alpha1"
	"github.com/henrywhitaker3/sre-operator/internal/http/webhook"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

var (
	ErrUnknwonHook   = errors.New("unknown hook")
	ErrUnknwonAction = errors.New("unknown action")
)

type RolloutHandler struct {
	ctx    context.Context
	client client.Client
	req    ctrl.Request
	obj    *v1alpha1.Rollout
	store  *webhook.Store
}

func NewRolloutHandler(ctx context.Context, client client.Client, req ctrl.Request, store *webhook.Store) *RolloutHandler {
	return &RolloutHandler{
		ctx:    ctx,
		client: client,
		req:    req,
		store:  store,
	}
}

func (h *RolloutHandler) Get() error {
	d := &v1alpha1.Rollout{}
	if err := h.client.Get(h.ctx, h.req.NamespacedName, d); err != nil {
		return err
	}
	h.obj = d
	return nil
}

func (h *RolloutHandler) CreateOrUpdate() (error, bool) {
	if h.obj.Spec.Hook == "" {
		return ErrUnknwonHook, false
	}

	_, ok := h.store.Get(h.obj.Spec.Hook)
	if !ok {
		return ErrUnknwonHook, true
	}

	var subs webhook.StoreSubscriber
	switch h.obj.Spec.Action {
	case "restart":
		subs = h.buildRestartFunc()
	default:
		return ErrUnknwonAction, false
	}

	if h.obj.Spec.Throttle != "" {
		if _, err := time.ParseDuration(h.obj.Spec.Throttle); err != nil {
			return err, false
		}
		subs = h.throttle(subs)
	}

	h.store.StoreFunc(h.obj.Spec.Hook, h.obj.Name, subs)
	return nil, true
}

func (h *RolloutHandler) buildRestartFunc() webhook.StoreSubscriber {
	return func(ctx context.Context) error {
		get := func(t client.Object) error {
			if err := h.client.Get(ctx, types.NamespacedName{
				Namespace: h.obj.Spec.Target.Namespace,
				Name:      h.obj.Spec.Target.Name,
			}, t); err != nil {
				return err
			}
			return nil
		}

		var target client.Object
		switch h.obj.Spec.Target.Kind {
		case "deployment":
			t := &appsv1.Deployment{}
			if err := get(t); err != nil {
				return err
			}
			if t.Spec.Template.Annotations == nil {
				t.Spec.Template.Annotations = make(map[string]string)
			}
			t.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
			target = t
		case "daemonset":
			t := &appsv1.DaemonSet{}
			if err := get(t); err != nil {
				return err
			}
			if t.Spec.Template.Annotations == nil {
				t.Spec.Template.Annotations = make(map[string]string)
			}
			t.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
			target = t
		case "statefulset":
			t := &appsv1.StatefulSet{}
			if err := get(t); err != nil {
				return err
			}
			if t.Spec.Template.Annotations == nil {
				t.Spec.Template.Annotations = make(map[string]string)
			}
			t.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
			target = t
		}

		return h.client.Patch(ctx, target, client.Merge)
	}
}

func (h *RolloutHandler) throttle(f webhook.StoreSubscriber) webhook.StoreSubscriber {
	dur, _ := time.ParseDuration(h.obj.Spec.Throttle)
	throttle := flow.Throttle[struct{}](func(ctx context.Context) (struct{}, error) {
		return struct{}{}, f(ctx)
	}, dur)
	return func(ctx context.Context) error {
		_, err := throttle(ctx)
		return err
	}
}

func (h *RolloutHandler) Delete() error {
	return h.store.DropFunc(h.obj.Spec.Hook, h.obj.Name)
}

func (h *RolloutHandler) DeletionTimestampIsZero() bool {
	return h.obj.DeletionTimestamp.IsZero()
}

func (h *RolloutHandler) AddFinalizer(fn string) error {
	controllerutil.AddFinalizer(h.obj, fn)
	return h.client.Update(h.ctx, h.obj)
}

func (h *RolloutHandler) RemoveFinalizer(fn string) error {
	controllerutil.RemoveFinalizer(h.obj, fn)
	return h.client.Update(h.ctx, h.obj)
}

func (h *RolloutHandler) SuccessStatus() error {
	h.obj.Status = v1alpha1.RolloutStatus{
		Registered: true,
	}
	return h.client.SubResource("status").Update(h.ctx, h.obj)
}

func (h *RolloutHandler) ErrorStatus(err error) error {
	h.obj.Status = v1alpha1.RolloutStatus{
		Registered: false,
		Error:      err.Error(),
	}
	return h.client.SubResource("status").Update(h.ctx, h.obj)
}

var _ Handler = &RolloutHandler{}
