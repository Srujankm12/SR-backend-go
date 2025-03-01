package models

import "time"

type SalesReport struct {
	UserID         string    `json:"user_id"`
	EmployeeID     string    `json:"emp_id"`
	Work           string    `json:"work"`
	TodaysWorkPlan string    `json:"todays_work_plan"`
	LoginTime      time.Time `json:"login_time"`
	CreatedAt      time.Time `json:"created_at"`
	ReportDate     string    `json:"report_date"`
}

type LogoutSummary struct {
	UserID                string    `json:"user_id"`
	EmpID                 string    `json:"emp_id"`
	TotalNoOfVisits       int       `json:"total_no_of_visits"`
	TotalNoOfColdCalls    int       `json:"total_no_of_cold_calls"`
	TotalNoOfFollowUps    int       `json:"total_no_of_follow_ups"`
	TotalEnquiryGenerated int       `json:"total_enquiry_generated"`
	TotalEnquiryValue     float64   `json:"total_enquiry_value"`
	TotalOrderLost        int       `json:"total_order_lost"`
	TotalOrderLostValue   float64   `json:"total_order_lost_value"`
	TotalOrderWon         int       `json:"total_order_won"`
	TotalOrderWonValue    float64   `json:"total_order_won_value"`
	CustomerFollowUpName  string    `json:"customer_follow_up_name"`
	Notes                 string    `json:"notes"`
	TomorrowGoals         string    `json:"tomorrow_goals"`
	HowWasToday           string    `json:"how_was_today"`
	WorkLocation          string    `json:"work_location"`
	LogoutTime            time.Time `json:"logout_time"`
	ReportDate            string    `json:"report_date"`
}
