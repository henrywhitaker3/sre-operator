package webhook

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/henrywhitaker3/flow"
	"github.com/henrywhitaker3/sre-operator/internal/metrics"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

type WebhookRoute struct {
	store   *Store
	metrics *metrics.Metrics
}

func NewWebhookRoute(store *Store, metrics *metrics.Metrics) *WebhookRoute {
	return &WebhookRoute{
		store:   store,
		metrics: metrics,
	}
}

func (w *WebhookRoute) Path() string {
	return "/webhook/:id"
}

func (w *WebhookRoute) Handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		hook := c.Param("id")
		cbs, ok := w.store.Get(hook)
		if !ok {
			return c.NoContent(http.StatusNotFound)
		}

		w.metrics.WebhooksCalled.With(prometheus.Labels{"id": hook}).Inc()

		for name, cb := range cbs {
			fmt.Printf("triggering action %s\n", name)
			go func(name string, f StoreSubscriber) {
				status := "success"
				if err := f(context.Background()); err != nil {
					if errors.Is(err, flow.ErrThrottled) {
						status = "throttled"
					} else {
						fmt.Println(err)
					}
				}

				w.metrics.ActionsRun.With(prometheus.Labels{
					"action":  name,
					"trigger": hook,
					"status":  status,
				}).Inc()
			}(name, cb)
		}

		return c.NoContent(http.StatusOK)
	}
}
