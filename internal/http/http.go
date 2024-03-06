package http

import (
	"context"

	"github.com/henrywhitaker3/sre-operator/internal/http/webhook"
	"github.com/henrywhitaker3/sre-operator/internal/metrics"
	"github.com/henrywhitaker3/sre-operator/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Http struct {
	echo *echo.Echo
	addr string
}

type Route interface {
	Path() string
	Handler() echo.HandlerFunc
}

func NewHttp(addr string, store *store.Store, metrics *metrics.Metrics) *Http {
	h := &Http{
		echo: echo.New(),
		addr: addr,
	}

	h.echo.HideBanner = true
	h.echo.Use(middleware.Logger())
	h.POST(webhook.NewWebhookRoute(store, metrics))

	return h
}

func (h *Http) GET(r Route) {
	h.echo.GET(r.Path(), r.Handler())
}

func (h *Http) POST(r Route) {
	h.echo.POST(r.Path(), r.Handler())
}

func (h *Http) Start(ctx context.Context) error {
	go h.echo.Start(h.addr)
	<-ctx.Done()
	return h.echo.Close()
}
