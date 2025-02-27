package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Srujankm12/SRproject/repository"
	"github.com/gorilla/mux"
)

type SalesHandler struct {
	Repo *repository.SalesRepository
}

// NewSalesHandler initializes a new SalesHandler
func NewSalesHandler(repo *repository.SalesRepository) *SalesHandler {
	return &SalesHandler{Repo: repo}
}

// Create a new sales report
func (h *SalesHandler) CreateSalesReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request struct {
		UserID         string `json:"user_id"` // Ideally, user_id should come from authentication middleware
		Work           string `json:"work"`
		TodaysWorkPlan string `json:"todays_work_plan"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Insert sales report and generate emp_id
	empID, err := h.Repo.InsertSalesReport(request.UserID, request.Work, request.TodaysWorkPlan)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Sales report created successfully",
		"emp_id":  empID,
	})
}

// Handle fetching a sales report
func (h *SalesHandler) GetSalesReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"] // Extract user ID from URL

	fmt.Println("Received request for user ID:", userID) // Log for debugging

	report, err := h.Repo.FetchSalesReport(userID)
	if err != nil {
		fmt.Println("Sales report not found for user ID:", userID) // Additional log
		http.Error(w, `{"error": "Sales report not found"}`, http.StatusNotFound)
		return
	}

	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// // Handle updating a sales report
// func (h *SalesHandler) UpdateSalesReport(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	var request struct {
// 		UserID         string `json:"user_id"`
// 		Work           string `json:"work"`
// 		TodaysWorkPlan string `json:"todays_work_plan"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
// 		return
// 	}

// 	err := h.Repo.UpdateSalesReport(request.UserID, request.Work, request.TodaysWorkPlan)
// 	if err != nil {
// 		http.Error(w, `{"error": "Failed to update sales report"}`, http.StatusInternalServerError)
// 		return
// 	}

//		w.WriteHeader(http.StatusOK)
//		json.NewEncoder(w).Encode(map[string]string{"message": "Sales report updated successfully"})
//	}
func (h *SalesHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Define request payload structure
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

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// Ensure user_id is provided
	if request.UserID == "" {
		http.Error(w, `{"error": "User ID is required"}`, http.StatusBadRequest)
		return
	}

	// Fetch emp_id for the user from today's sales report
	empID, err := h.Repo.GetEmpIDByUserID(request.UserID)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Call repository method to insert the logout summary
	err = h.Repo.InsertLogoutSummary(
		request.UserID, empID, request.CustomerFollowUpName, request.Notes, request.TomorrowGoals,
		request.HowWasToday, request.WorkLocation, request.TotalNoOfVisits, request.TotalNoOfColdCalls,
		request.TotalNoOfFollowUps, request.TotalEnquiryGenerated, request.TotalOrderLost, request.TotalOrderWon,
		request.TotalOrderLostValue, request.TotalOrderWonValue, request.TotalEnquiryValue,
	)

	// Handle errors
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout summary recorded successfully",
		"emp_id":  empID,
	})
}
