package repository

import (
	"database/sql"
	"time"

	"github.com/Srujankm12/SRproject/internal/models"
)

// SalesRepository handles database operations for sales reports.
type SalesRepository struct {
	DB *sql.DB
}

// NewSalesRepository creates a new repository instance.
func NewSalesRepository(db *sql.DB) *SalesRepository {
	return &SalesRepository{DB: db}
}

// CheckIfUserLoggedIn checks if a user has an active login session.
func (r *SalesRepository) CheckIfUserLoggedIn(userID string) (bool, error) {
	var loginTime time.Time
	err := r.DB.QueryRow("SELECT login_time FROM sales_reports WHERE user_id = ? ORDER BY created_at DESC LIMIT 1", userID).Scan(&loginTime)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// InsertSalesReport logs a user's sales report entry.
func (r *SalesRepository) InsertSalesReport(report models.SalesReport) error {
	_, err := r.DB.Exec(`
		INSERT INTO sales_reports (user_id, emp_id, work, todays_work_plan, login_time, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		report.UserID, report.EmployeeID, report.Work, report.TodaysWorkPlan, report.LoginTime, report.CreatedAt)
	return err
}

// InsertLogoutSummary logs a user's logout summary.
func (r *SalesRepository) InsertLogoutSummary(summary models.LogoutSummary) error {
	_, err := r.DB.Exec(`
		INSERT INTO logout_summaries (user_id, emp_id, total_no_of_visits, total_no_of_cold_calls, 
		total_no_of_follow_ups, total_enquiry_generated, total_enquiry_value, total_order_lost, 
		total_order_lost_value, total_order_won, total_order_won_value, customer_follow_up_name, 
		notes, tomorrow_goals, how_was_today, work_location, logout_time) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		summary.UserID, summary.EmployeeID, summary.TotalNoOfVisits, summary.TotalNoOfColdCalls,
		summary.TotalNoOfFollowUps, summary.TotalEnquiryGenerated, summary.TotalEnquiryValue,
		summary.TotalOrderLost, summary.TotalOrderLostValue, summary.TotalOrderWon,
		summary.TotalOrderWonValue, summary.CustomerFollowUpName, summary.Notes,
		summary.TomorrowGoals, summary.HowWasToday, summary.WorkLocation, summary.LogoutTime)
	return err
}

// GetUserSalesReport retrieves a user's latest sales report.
func (r *SalesRepository) GetUserSalesReport(userID string) (*models.SalesReport, error) {
	var report models.SalesReport
	err := r.DB.QueryRow(`
		SELECT user_id, emp_id, work, todays_work_plan, login_time, created_at
		FROM sales_reports WHERE user_id = ? ORDER BY created_at DESC LIMIT 1`,
		userID).Scan(&report.UserID, &report.EmployeeID, &report.Work, &report.TodaysWorkPlan, &report.LoginTime, &report.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &report, nil
}

// GetUserLogoutSummary retrieves a user's latest logout summary.
func (r *SalesRepository) GetUserLogoutSummary(userID string) (*models.LogoutSummary, error) {
	var summary models.LogoutSummary
	err := r.DB.QueryRow(`
		SELECT user_id, emp_id, total_no_of_visits, total_no_of_cold_calls, total_no_of_follow_ups, 
		total_enquiry_generated, total_enquiry_value, total_order_lost, total_order_lost_value, 
		total_order_won, total_order_won_value, customer_follow_up_name, notes, tomorrow_goals, 
		how_was_today, work_location, logout_time
		FROM logout_summaries WHERE user_id = ? ORDER BY logout_time DESC LIMIT 1`,
		userID).Scan(
		&summary.UserID, &summary.EmployeeID, &summary.TotalNoOfVisits, &summary.TotalNoOfColdCalls,
		&summary.TotalNoOfFollowUps, &summary.TotalEnquiryGenerated, &summary.TotalEnquiryValue,
		&summary.TotalOrderLost, &summary.TotalOrderLostValue, &summary.TotalOrderWon,
		&summary.TotalOrderWonValue, &summary.CustomerFollowUpName, &summary.Notes,
		&summary.TomorrowGoals, &summary.HowWasToday, &summary.WorkLocation, &summary.LogoutTime)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &summary, nil
}
