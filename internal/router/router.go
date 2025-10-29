// Package router пакет с инициализацией роутера для маршутизации
package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/sinfirst/Test-Task-For-Effective-Mobile/internal/app"
	"github.com/sinfirst/Test-Task-For-Effective-Mobile/internal/middleware/logging"
)

// NewRouter описание всех эндпоинтов
func NewRouter(a *app.App) *chi.Mux {
	router := chi.NewRouter()

	router.Use(logging.WithLogging)
	router.Get("/api/read", a.ReadSub)
	router.Get("/api/list", a.ListSub)
	router.Post("/api/create", a.CreateSub)
	router.Put("/api/update", a.UpdateSub)
	router.Delete("/api/delete", a.DeleteSub)

	router.Post("/api/calc", a.CalcSumSub)
	return router
}
