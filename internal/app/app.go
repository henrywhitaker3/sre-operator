package app

import (
	"github.com/henrywhitaker3/sre-operator/internal/http"
	"github.com/henrywhitaker3/sre-operator/internal/http/webhook"
	"github.com/henrywhitaker3/sre-operator/internal/metrics"
)

type App struct {
	Http      *http.Http
	HookStore *webhook.Store
	Metrics   *metrics.Metrics
}

func NewApp(addr string) *App {
	store := webhook.NewStore()
	metrics, err := metrics.New()
	if err != nil {
		panic(err)
	}

	return &App{
		Http:      http.NewHttp(addr, store, metrics),
		HookStore: store,
		Metrics:   metrics,
	}
}
