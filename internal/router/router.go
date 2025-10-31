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
	router.Get("/{id}", a.ReadSub)
	router.Get("/list/{user_id}", a.ListSub)
	router.Post("/", a.CreateSub)
	router.Put("/", a.UpdateSub)
	router.Delete("/", a.DeleteSub)

	router.Post("/calc", a.CalcSumSub)
	return router
}
