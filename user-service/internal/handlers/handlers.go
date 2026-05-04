package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/asdlc-repos/testingnewtaskpage3979/user-service/internal/store"
)

type Handler struct {
	store *store.Store
}

func New(s *store.Store) *Handler {
	return &Handler{store: s}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) http.Handler {
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/api/v1/users/", h.usersRouter)
	return corsMiddleware(mux)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// usersRouter dispatches to sub-handlers based on path segments.
func (h *Handler) usersRouter(w http.ResponseWriter, r *http.Request) {
	// Strip prefix "/api/v1/users/"
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
	parts := strings.Split(strings.Trim(path, "/"), "/")

	switch {
	case len(parts) == 1 && parts[0] != "":
		// GET /api/v1/users/{userId}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.getUser(w, r, parts[0])

	case len(parts) == 2 && parts[1] == "balance":
		// GET /api/v1/users/{userId}/balance
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.getBalance(w, r, parts[0])

	case len(parts) == 3 && parts[1] == "balance" && parts[2] == "deduct":
		// PUT /api/v1/users/{userId}/balance/deduct
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.deductBalance(w, r, parts[0])

	case len(parts) == 2 && parts[1] == "reports":
		// GET /api/v1/users/{managerId}/reports
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.getReports(w, r, parts[0])

	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request, userId string) {
	user, err := h.store.GetUser(userId)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *Handler) getBalance(w http.ResponseWriter, r *http.Request, userId string) {
	balance, err := h.store.GetBalance(userId)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, balance)
}

func (h *Handler) deductBalance(w http.ResponseWriter, r *http.Request, userId string) {
	var req struct {
		Days float64 `json:"days"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	balance, err := h.store.DeductBalance(userId, req.Days)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			http.Error(w, "user not found", http.StatusNotFound)
		case errors.Is(err, store.ErrInvalidDays):
			http.Error(w, "days must be greater than zero", http.StatusBadRequest)
		case errors.Is(err, store.ErrInsufficientBalance):
			http.Error(w, "insufficient leave balance", http.StatusBadRequest)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, balance)
}

func (h *Handler) getReports(w http.ResponseWriter, r *http.Request, managerId string) {
	reports, err := h.store.GetDirectReports(managerId)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			http.Error(w, "manager not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if reports == nil {
		reports = []string{}
	}
	writeJSON(w, http.StatusOK, reports)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
