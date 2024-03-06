package app

import (
	"github.com/henrywhitaker3/sre-operator/internal/http"
	"github.com/henrywhitaker3/sre-operator/internal/metrics"
	"github.com/henrywhitaker3/sre-operator/internal/store"
)

type App struct {
	Http      *http.Http
	HookStore *store.Store
	Metrics   *metrics.Metrics
}

func NewApp(addr string) *App {
	store := store.NewStore()
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
