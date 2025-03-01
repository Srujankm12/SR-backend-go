package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Srujankm12/SRproject/repository"
	"github.com/gorilla/mux"
)

// SalesHandler handles sales-related HTTP requests.
type SalesHandler struct {
	Repo *repository.SalesRepository
}

func NewSalesHandler(repo *repository.SalesRepository) *SalesHandler {
	return &SalesHandler{Repo: repo}
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func sendJSONError(w http.ResponseWriter, statusCode int, message string) {
	sendJSONResponse(w, statusCode, map[string]string{"error": message})
}

func (h *SalesHandler) CreateSalesReport(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID         string `json:"user_id"`
		Work           string `json:"work"`
		TodaysWorkPlan string `json:"todays_work_plan"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendJSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if request.UserID == "" || request.Work == "" || request.TodaysWorkPlan == "" {
		sendJSONError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	empID, err := h.Repo.InsertSalesReport(request.UserID, request.Work, request.TodaysWorkPlan)
	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create sales report: %v", err))
		return
	}

	sendJSONResponse(w, http.StatusCreated, map[string]string{
		"message": "Sales report created successfully",
		"emp_id":  empID,
	})
}

func (h *SalesHandler) GetSalesReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"] // Extract user_id from URL path

	if userID == "" {
		sendJSONError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	log.Println("Fetching report for user ID:", userID) // Debugging log

	report, err := h.Repo.FetchSalesReport(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendJSONError(w, http.StatusNotFound, "No sales report found for today")
		} else {
			sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch sales report: %v", err))
		}
		return
	}

	sendJSONResponse(w, http.StatusOK, report)
}
func (h *SalesHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID                string  `json:"user_id"`
		TotalNoOfVisits       int     `json:"total_no_of_visits"`
		TotalNoOfColdCalls    int     `json:"total_no_of_cold_calls"`
		TotalNoOfFollowUps    int     `json:"total_no_of_follow_ups"`
		TotalEnquiryGenerated int     `json:"total_enquiry_generated"`
		TotalEnquiryValue     float64 `json:"total_enquiry_value"`
		TotalOrderLost        int     `json:"total_order_lost"`
		TotalOrderLostValue   float64 `json:"total_order_lost_value"`
		TotalOrderWon         int     `json:"total_order_won"`
		TotalOrderWonValue    float64 `json:"total_order_won_value"`
		CustomerFollowUpName  string  `json:"customer_follow_up_name"`
		Notes                 string  `json:"notes"`
		TomorrowGoals         string  `json:"tomorrow_goals"`
		HowWasToday           string  `json:"how_was_today"`
		WorkLocation          string  `json:"work_location"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		sendJSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Ensure required fields are provided
	if request.UserID == "" {
		sendJSONError(w, http.StatusBadRequest, "User ID is required")
		return
	}

	// Fetch emp_id for the user from today's sales report
	empID, err := h.Repo.GetEmpIDByUserID(request.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendJSONError(w, http.StatusNotFound, "Employee not found for this user ID")
		} else {
			sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch emp_id: %v", err))
		}
		return
	}

	// Check if user has already logged out today
	exists, err := h.Repo.CheckLogoutExists(request.UserID)
	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to check logout status: %v", err))
		return
	}

	if exists {
		sendJSONError(w, http.StatusConflict, "User has already logged out today")
		return
	}

	// Insert the logout summary using the fetched empID
	err = h.Repo.InsertLogoutSummary(
		request.UserID, empID, request.CustomerFollowUpName, request.Notes, request.TomorrowGoals,
		request.HowWasToday, request.WorkLocation, request.TotalNoOfVisits, request.TotalNoOfColdCalls,
		request.TotalNoOfFollowUps, request.TotalEnquiryGenerated, request.TotalOrderLost, request.TotalOrderWon,
		request.TotalEnquiryValue, request.TotalOrderLostValue, request.TotalOrderWonValue,
	)
	if err != nil {
		sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to insert logout summary: %v", err))
		return
	}

	sendJSONResponse(w, http.StatusCreated, map[string]string{
		"message": "Logout summary recorded successfully",
		"emp_id":  empID,
	})
}

func (h *SalesHandler) GetLogoutSummary(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"] // Extract user_id from URL
	if userID == "" {
		sendJSONError(w, http.StatusBadRequest, "user_id is required")
		return
	}

	// Fetch logout history using user_id directly
	summary, err := h.Repo.GetLogoutSummary(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendJSONError(w, http.StatusNotFound, fmt.Sprintf("No logout history found for user_id %s", userID))
		} else {
			sendJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to fetch logout summary: %v", err))
		}
		return
	}

	sendJSONResponse(w, http.StatusOK, summary)
}
