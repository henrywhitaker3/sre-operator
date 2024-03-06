package app

import (
	"github.com/henrywhitaker3/sre-operator/internal/http"
	"github.com/henrywhitaker3/sre-operator/internal/http/webhook"
)

type App struct {
	Http      *http.Http
	HookStore *webhook.Store
}

func NewApp(addr string) *App {
	store := webhook.NewStore()
	return &App{
		Http:      http.NewHttp(addr, store),
		HookStore: store,
	}
}
