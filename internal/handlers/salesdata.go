package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Srujankm12/SRproject/repository"
)

type SalesHandler struct {
	repo *repository.SalesRepository
}

func NewSalesHandler(repo *repository.SalesRepository) *SalesHandler {
	return &SalesHandler{repo: repo}
}

// 🟢 User Login
func (h *SalesHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	workLocation := r.FormValue("work_location")
	todaysWorkPlan := r.FormValue("todays_work_plan")

	empID, err := h.repo.UserLogin(userID, workLocation, todaysWorkPlan)
	if err != nil {
		http.Error(w, "Login failed", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"emp_id":  empID,
	})
}

func (h *SalesHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	empID := r.FormValue("emp_id")

	totalVisits := parseInt(r.FormValue("total_no_of_visits"))
	coldCalls := parseInt(r.FormValue("total_no_of_cold_calls"))
	followUps := parseInt(r.FormValue("total_no_of_customer_follow_up"))
	enquiries := parseInt(r.FormValue("total_enquiry_generated"))
	ordersLost := parseInt(r.FormValue("total_order_lost"))
	ordersWon := parseInt(r.FormValue("total_order_won"))
	notes := r.FormValue("notes")
	tomorrowGoals := r.FormValue("tomorrow_goals")
	howWasToday := r.FormValue("how_was_today")

	err := h.repo.UserLogout(userID, empID, totalVisits, coldCalls, followUps, enquiries, ordersLost, ordersWon, notes, tomorrowGoals, howWasToday)
	if err != nil {
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}

// 🟢 Site Check-In
func (h *SalesHandler) CheckInHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	empID := r.FormValue("emp_id")
	companyName := r.FormValue("company_name")
	purpose := r.FormValue("purpose")

	err := h.repo.SiteCheckIn(userID, empID, companyName, purpose)
	if err != nil {
		http.Error(w, "Check-in failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Check-in successful",
	})
}

// 🟢 Site Check-Out
func (h *SalesHandler) CheckOutHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	empID := r.FormValue("emp_id")
	engineerName := r.FormValue("engineer_name")
	companySalesStage := r.FormValue("company_sales_stage")
	visitOn := r.FormValue("visit_on")
	timelineForNextActionPlan := r.FormValue("timeline_for_next_action_plan")
	challenges := r.FormValue("challenges")
	visitRating := parseInt(r.FormValue("visit_rating"))
	resultOfVisit := r.FormValue("result_of_visit")
	notes := r.FormValue("notes")

	err := h.repo.SiteCheckOut(userID, empID, engineerName, companySalesStage, visitOn, timelineForNextActionPlan, challenges, resultOfVisit, notes, visitRating)
	if err != nil {
		http.Error(w, "Check-out failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Check-out successful",
	})
}

func parseInt(value string) int {
	num, _ := strconv.Atoi(value)
	return num
}
