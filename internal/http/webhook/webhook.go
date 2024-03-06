package webhook

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/henrywhitaker3/flow"
	"github.com/labstack/echo/v4"
)

type WebhookRoute struct {
	store *Store
}

func NewWebhookRoute(store *Store) *WebhookRoute {
	return &WebhookRoute{
		store: store,
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

		for name, cb := range cbs {
			fmt.Printf("triggering action %s\n", name)
			go func(f StoreSubscriber) {
				if err := f(context.Background()); err != nil {
					if errors.Is(err, flow.ErrThrottled) {
						fmt.Println("throttled")
						return
					}
					fmt.Println(err)
				}
			}(cb)
		}

		return c.NoContent(http.StatusOK)
	}
}
