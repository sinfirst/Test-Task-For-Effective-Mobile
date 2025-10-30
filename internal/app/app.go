package app

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sinfirst/Test-Task-For-Effective-Mobile/internal/config"
	"github.com/sinfirst/Test-Task-For-Effective-Mobile/internal/models"
)

// Storage интрефейс для работы с базой данных
type Storage interface {
	CreateInDB(ctx context.Context, sub models.SubJSON) (int, error)
	ReadFromDB(ctx context.Context, id string) (models.SubJSON, error)
	UpdateInDB(ctx context.Context, sub models.SubJSON) error
	DeleteFromDB(ctx context.Context, id string) error
	ListFromDB(ctx context.Context, id string) ([]models.SubJSON, error)
	CostSumSubFromDB(ctx context.Context, req models.SubJSON) (int, error)
}

// App структура для хранения переменных для хендлеров
type App struct {
	storage Storage
	config  config.Config
}

// NewApp конструктор для App
func NewApp(storage Storage, config config.Config) *App {
	return &App{storage: storage, config: config}
}

func (a *App) CreateSub(w http.ResponseWriter, r *http.Request) {
	var req models.SubJSON

	json.NewDecoder(r.Body).Decode(&req)

	id, err := a.storage.CreateInDB(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.SubRespJSON{ID: id})
}
func (a *App) ReadSub(w http.ResponseWriter, r *http.Request) {
	var appErr models.AppError
	id := chi.URLParam(r, "id")
	resp, err := a.storage.ReadFromDB(r.Context(), id)

	if errors.As(err, &appErr); appErr == models.ErrNotFound {
		http.Error(w, "id not found", http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (a *App) UpdateSub(w http.ResponseWriter, r *http.Request) {
	var req models.SubJSON
	var appErr models.AppError

	json.NewDecoder(r.Body).Decode(&req)

	err := a.storage.UpdateInDB(r.Context(), req)
	if errors.As(err, &appErr); appErr == models.ErrNotFound {
		http.Error(w, "id not found", http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *App) DeleteSub(w http.ResponseWriter, r *http.Request) {
	var req models.SubJSON
	var appErr models.AppError

	json.NewDecoder(r.Body).Decode(&req)

	err := a.storage.DeleteFromDB(r.Context(), req.ID)
	if errors.As(err, &appErr); appErr == models.ErrNotFound {
		http.Error(w, "id not found", http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *App) ListSub(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "user_id")
	resp, err := a.storage.ListFromDB(r.Context(), id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
func (a *App) CalcSumSub(w http.ResponseWriter, r *http.Request) {
	var req models.SubJSON

	json.NewDecoder(r.Body).Decode(&req)

	sum, err := a.storage.CostSumSubFromDB(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.CostSumSubRespJSON{Sum: sum})
}
