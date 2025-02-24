package models

import "time"

type User struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SalesReport struct {
	EmpID          string    `json:"emp_id"`
	UserID         string    `json:"user_id"`
	Work           string    `json:"work"`
	TodaysWorkPlan string    `json:"todays_work_plan"`
	CreatedAt      time.Time `json:"created_at"`
}

type SiteVisit struct {
	VisitID     string    `json:"visit_id"`
	EmpID       string    `json:"emp_id"`
	UserID      string    `json:"user_id"`
	CompanyName string    `json:"company_name"`
	Purpose     string    `json:"purpose"`
	CheckInTime time.Time `json:"checkin_time"`
}

type SiteCheckout struct {
	CheckoutID                string    `json:"checkout_id"`
	VisitID                   string    `json:"visit_id"`
	EmpID                     string    `json:"emp_id"`
	UserID                    string    `json:"user_id"`
	EngineerName              string    `json:"engineer_name"`
	CompanySalesStage         string    `json:"company_sales_stage"`
	VisitOn                   time.Time `json:"visit_on"`
	TimelineForNextActionPlan string    `json:"timeline_for_next_action_plan"`
	Challenges                string    `json:"challenges"`
	VisitRating               int       `json:"visit_rating"`
	ResultOfVisit             string    `json:"result_of_visit"`
	Notes                     string    `json:"notes"`
	CheckoutTime              time.Time `json:"checkout_time"`
}

type LogoutSummary struct {
	LogoutID              string    `json:"logout_id"`
	UserID                string    `json:"user_id"`
	TotalNoOfVisits       int       `json:"total_no_of_visits"`
	TotalNoOfColdCalls    int       `json:"total_no_of_cold_calls"`
	TotalNoOfFollowUps    int       `json:"total_no_of_follow_ups"`
	TotalEnquiryGenerated int       `json:"total_enquiry_generated"`
	TotalEnquiryValue     int       `json:"total_enquiry_value"`
	TotalOrderLost        int       `json:"total_order_lost"`
	TotalOrderLostValue   int       `json:"total_order_lost_value"`
	TotalOrderWon         int       `json:"total_order_won"`
	TotalOrderWonValue    int       `json:"total_order_won_value"`
	CustomerFollowUpName  string    `json:"customer_follow_up_name"`
	Notes                 string    `json:"notes"`
	TomorrowGoals         string    `json:"tomorrow_goals"`
	HowWasToday           string    `json:"how_was_today"`
	WorkLocation          string    `json:"work_location"`
	LogoutTime            time.Time `json:"logout_time"`
}
