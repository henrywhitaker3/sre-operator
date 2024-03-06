package webhook

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/henrywhitaker3/flow"
	"github.com/henrywhitaker3/sre-operator/internal/metrics"
	"github.com/henrywhitaker3/sre-operator/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

type WebhookRoute struct {
	store   *store.Store
	metrics *metrics.Metrics
}

func NewWebhookRoute(store *store.Store, metrics *metrics.Metrics) *WebhookRoute {
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
		subs, err := w.store.Get(store.WEBHOOK, hook)
		if err != nil {
			fmt.Println(err)
			return c.NoContent(http.StatusNotFound)
		}

		w.metrics.WebhooksCalled.With(prometheus.Labels{"id": hook}).Inc()

		for _, s := range subs {
			fmt.Printf("triggering action %s\n", s.Name)
			go func(s store.Subscription) {
				status := "success"
				if err := s.Do(context.Background()); err != nil {
					if errors.Is(err, flow.ErrThrottled) {
						status = "throttled"
					} else {
						fmt.Println(err)
					}
				}

				w.metrics.ActionsRun.With(prometheus.Labels{
					"action":  s.Name,
					"trigger": hook,
					"status":  status,
				}).Inc()
			}(s)
		}

		return c.NoContent(http.StatusOK)
	}
}
