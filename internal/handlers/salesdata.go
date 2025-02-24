package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Srujankm12/SRproject/internal/models"
	"github.com/Srujankm12/SRproject/repository"
)

// SalesHandler manages sales reports and logout summaries.
type SalesHandler struct {
	Repo *repository.SalesRepository
}

// NewSalesHandler initializes the handler with the repository.
func NewSalesHandler(db *sql.DB) *SalesHandler {
	return &SalesHandler{
		Repo: repository.NewSalesRepository(db),
	}
}

// HandleCheckIn handles sales report entry (check-in).
func (h *SalesHandler) HandleCheckIn(w http.ResponseWriter, r *http.Request) {
	var report models.SalesReport

	// Parse request body
	err := json.NewDecoder(r.Body).Decode(&report)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Set timestamps
	report.CreatedAt = time.Now()

	// Check if the user is already logged in
	isLoggedIn, err := h.Repo.CheckIfUserLoggedIn(report.UserID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if isLoggedIn {
		http.Error(w, "User already checked in", http.StatusConflict)
		return
	}

	// Insert check-in report
	err = h.Repo.InsertSalesReport(report)
	if err != nil {
		http.Error(w, "Failed to check in", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Checked in successfully"})
}

// HandleCheckOut handles user logout summary (check-out).
func (h *SalesHandler) HandleCheckOut(w http.ResponseWriter, r *http.Request) {
	var summary models.LogoutSummary

	// Parse request body
	err := json.NewDecoder(r.Body).Decode(&summary)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Set logout timestamp
	summary.LogoutTime = time.Now()

	// Insert logout summary
	err = h.Repo.InsertLogoutSummary(summary)
	if err != nil {
		http.Error(w, "Failed to log out", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// HandleGetSalesReport retrieves a user's latest sales report.
func (h *SalesHandler) HandleGetSalesReport(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	report, err := h.Repo.GetUserSalesReport(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve sales report", http.StatusInternalServerError)
		return
	}

	if report == nil {
		http.Error(w, "No sales report found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(report)
}

// HandleGetLogoutSummary retrieves a user's latest logout summary.
func (h *SalesHandler) HandleGetLogoutSummary(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	summary, err := h.Repo.GetUserLogoutSummary(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve logout summary", http.StatusInternalServerError)
		return
	}

	if summary == nil {
		http.Error(w, "No logout summary found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(summary)
}
