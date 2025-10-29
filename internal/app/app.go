package app

import (
	"net/http"

	"github.com/sinfirst/Test-Task-For-Effective-Mobile/internal/config"
	"go.uber.org/zap"
)

// Storage интрефейс для работы с базой данных
type Storage interface {
}

// App структура для хранения переменных для хендлеров
type App struct {
	storage Storage
	config  config.Config
	logger  zap.SugaredLogger
}

// NewApp конструктор для App
func NewApp(storage Storage, config config.Config, logger zap.SugaredLogger) *App {
	return &App{storage: storage, config: config, logger: logger}
}

func (a *App) CreateSub(w http.ResponseWriter, r *http.Request) {

}
func (a *App) ReadSub(w http.ResponseWriter, r *http.Request) {

}
func (a *App) UpdateSub(w http.ResponseWriter, r *http.Request) {

}
func (a *App) DeleteSub(w http.ResponseWriter, r *http.Request) {

}
func (a *App) ListSub(w http.ResponseWriter, r *http.Request) {

}
func (a *App) CalcSumSub(w http.ResponseWriter, r *http.Request) {

}
