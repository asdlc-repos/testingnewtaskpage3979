package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type LeaveStatus string

const (
	StatusPending  LeaveStatus = "pending"
	StatusApproved LeaveStatus = "approved"
	StatusRejected LeaveStatus = "rejected"
)

type LeaveRequest struct {
	Id             string      `json:"id"`
	UserId         string      `json:"userId"`
	StartDate      string      `json:"startDate"`
	EndDate        string      `json:"endDate"`
	Reason         string      `json:"reason"`
	Status         LeaveStatus `json:"status"`
	ManagerId      string      `json:"managerId"`
	ManagerComment string      `json:"managerComment,omitempty"`
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`
}

type Store struct {
	mu       sync.RWMutex
	requests map[string]*LeaveRequest
}

var store = &Store{
	requests: make(map[string]*LeaveRequest),
}

var userServiceURL string

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

type User struct {
	Id        string `json:"id"`
	ManagerId string `json:"managerId"`
	Name      string `json:"name"`
}

func getUserFromService(ctx context.Context, userId string) (*User, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s", userServiceURL, userId)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user service returned %d", resp.StatusCode)
	}
	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func getDirectReports(ctx context.Context, managerId string) ([]string, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s/direct-reports", userServiceURL, managerId)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user service returned %d", resp.StatusCode)
	}
	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}
	ids := make([]string, len(users))
	for i, u := range users {
		ids[i] = u.Id
	}
	return ids, nil
}

func deductLeaveBalance(ctx context.Context, userId string, days int) error {
	url := fmt.Sprintf("%s/api/v1/users/%s/leave-balance/deduct", userServiceURL, userId)
	body := fmt.Sprintf(`{"days":%d}`, days)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("user service returned %d when deducting balance", resp.StatusCode)
	}
	return nil
}

func calculateDays(startDate, endDate string) (int, error) {
	layout := "2006-01-02"
	start, err := time.Parse(layout, startDate)
	if err != nil {
		return 0, fmt.Errorf("invalid startDate format, expected YYYY-MM-DD")
	}
	end, err := time.Parse(layout, endDate)
	if err != nil {
		return 0, fmt.Errorf("invalid endDate format, expected YYYY-MM-DD")
	}
	days := int(end.Sub(start).Hours()/24) + 1
	if days <= 0 {
		return 0, fmt.Errorf("endDate must be on or after startDate")
	}
	return days, nil
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("→ %s %s", r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
		log.Printf("← %s %s (%v)", r.Method, r.URL.RequestURI(), time.Since(start))
	})
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func errorResponse(w http.ResponseWriter, status int, msg string) {
	jsonResponse(w, status, map[string]string{"error": msg})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func listRequestsHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	managerId := r.URL.Query().Get("managerId")
	status := r.URL.Query().Get("status")

	var employeeIds []string
	if managerId != "" {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		ids, err := getDirectReports(ctx, managerId)
		if err != nil {
			log.Printf("Warning: could not fetch direct reports for %s: %v", managerId, err)
		} else {
			employeeIds = ids
		}
	}

	store.mu.RLock()
	defer store.mu.RUnlock()

	result := make([]*LeaveRequest, 0)
	for _, req := range store.requests {
		if userId != "" && req.UserId != userId {
			continue
		}
		if managerId != "" {
			matchesDirect := false
			for _, id := range employeeIds {
				if req.UserId == id {
					matchesDirect = true
					break
				}
			}
			if !matchesDirect && req.ManagerId != managerId {
				continue
			}
		}
		if status != "" && string(req.Status) != status {
			continue
		}
		result = append(result, req)
	}

	jsonResponse(w, http.StatusOK, result)
}

func createRequestHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserId    string `json:"userId"`
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
		Reason    string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.UserId == "" || input.StartDate == "" || input.EndDate == "" || input.Reason == "" {
		errorResponse(w, http.StatusBadRequest, "userId, startDate, endDate, and reason are required")
		return
	}

	if _, err := calculateDays(input.StartDate, input.EndDate); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	user, err := getUserFromService(ctx, input.UserId)
	if err != nil {
		log.Printf("Error fetching user %s: %v", input.UserId, err)
		errorResponse(w, http.StatusInternalServerError, "failed to verify user")
		return
	}
	if user == nil {
		errorResponse(w, http.StatusBadRequest, "user not found")
		return
	}

	now := time.Now()
	req := &LeaveRequest{
		Id:        uuid.New().String(),
		UserId:    input.UserId,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		Reason:    input.Reason,
		Status:    StatusPending,
		ManagerId: user.ManagerId,
		CreatedAt: now,
		UpdatedAt: now,
	}

	store.mu.Lock()
	store.requests[req.Id] = req
	store.mu.Unlock()

	jsonResponse(w, http.StatusCreated, req)
}

func approveRequestHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	id := strings.TrimPrefix(path, "/api/v1/requests/")
	id = strings.TrimSuffix(id, "/approve")

	managerId := r.URL.Query().Get("managerId")
	if managerId == "" {
		errorResponse(w, http.StatusBadRequest, "managerId query parameter is required")
		return
	}

	var input struct {
		Comment string `json:"comment"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	store.mu.Lock()
	defer store.mu.Unlock()

	req, ok := store.requests[id]
	if !ok {
		errorResponse(w, http.StatusNotFound, "leave request not found")
		return
	}

	if req.ManagerId != managerId {
		errorResponse(w, http.StatusBadRequest, "manager not authorized for this leave request")
		return
	}

	if req.Status != StatusPending {
		errorResponse(w, http.StatusBadRequest, "only pending requests can be approved")
		return
	}

	days, err := calculateDays(req.StartDate, req.EndDate)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "failed to calculate leave days")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := deductLeaveBalance(ctx, req.UserId, days); err != nil {
		log.Printf("Error deducting leave balance for user %s: %v", req.UserId, err)
		errorResponse(w, http.StatusInternalServerError, "failed to deduct leave balance")
		return
	}

	req.Status = StatusApproved
	req.ManagerComment = input.Comment
	req.UpdatedAt = time.Now()

	jsonResponse(w, http.StatusOK, req)
}

func rejectRequestHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	id := strings.TrimPrefix(path, "/api/v1/requests/")
	id = strings.TrimSuffix(id, "/reject")

	managerId := r.URL.Query().Get("managerId")
	if managerId == "" {
		errorResponse(w, http.StatusBadRequest, "managerId query parameter is required")
		return
	}

	var input struct {
		Comment string `json:"comment"`
	}
	json.NewDecoder(r.Body).Decode(&input)

	store.mu.Lock()
	defer store.mu.Unlock()

	req, ok := store.requests[id]
	if !ok {
		errorResponse(w, http.StatusNotFound, "leave request not found")
		return
	}

	if req.ManagerId != managerId {
		errorResponse(w, http.StatusBadRequest, "manager not authorized for this leave request")
		return
	}

	if req.Status != StatusPending {
		errorResponse(w, http.StatusBadRequest, "only pending requests can be rejected")
		return
	}

	req.Status = StatusRejected
	req.ManagerComment = input.Comment
	req.UpdatedAt = time.Now()

	jsonResponse(w, http.StatusOK, req)
}

func seedData() {
	now := time.Now()
	requests := []*LeaveRequest{
		{
			Id:        uuid.New().String(),
			UserId:    "user-001",
			StartDate: "2026-05-10",
			EndDate:   "2026-05-12",
			Reason:    "Family vacation",
			Status:    StatusPending,
			ManagerId: "mgr-001",
			CreatedAt: now.Add(-24 * time.Hour),
			UpdatedAt: now.Add(-24 * time.Hour),
		},
		{
			Id:             uuid.New().String(),
			UserId:         "user-002",
			StartDate:      "2026-04-15",
			EndDate:        "2026-04-17",
			Reason:         "Medical appointment",
			Status:         StatusApproved,
			ManagerId:      "mgr-001",
			ManagerComment: "Approved, take care",
			CreatedAt:      now.Add(-96 * time.Hour),
			UpdatedAt:      now.Add(-72 * time.Hour),
		},
		{
			Id:             uuid.New().String(),
			UserId:         "user-003",
			StartDate:      "2026-05-20",
			EndDate:        "2026-05-21",
			Reason:         "Personal matters",
			Status:         StatusRejected,
			ManagerId:      "mgr-002",
			ManagerComment: "Team is short-staffed this week",
			CreatedAt:      now.Add(-48 * time.Hour),
			UpdatedAt:      now.Add(-24 * time.Hour),
		},
	}
	for _, r := range requests {
		store.requests[r.Id] = r
	}
	log.Printf("Seeded %d leave requests", len(requests))
}

func main() {
	userServiceURL = getEnv("USER_SERVICE_URL", "http://user-service:9091")
	port := getEnv("PORT", "9090")

	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)

	mux.HandleFunc("/api/v1/requests", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listRequestsHandler(w, r)
		case http.MethodPost:
			createRequestHandler(w, r)
		default:
			errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	})

	mux.HandleFunc("/api/v1/requests/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/approve") && r.Method == http.MethodPut:
			approveRequestHandler(w, r)
		case strings.HasSuffix(path, "/reject") && r.Method == http.MethodPut:
			rejectRequestHandler(w, r)
		default:
			errorResponse(w, http.StatusNotFound, "not found")
		}
	})

	seedData()

	handler := corsMiddleware(loggingMiddleware(mux))
	log.Printf("Starting leave-service on port %s (user-service: %s)", port, userServiceURL)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
