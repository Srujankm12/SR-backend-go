package repository

import (
	"database/sql"
	"fmt"

	"github.com/Srujankm12/SRproject/internal/models"
	"github.com/google/uuid"
)

type SalesRepository struct {
	db *sql.DB
}

func NewSalesRepository(db *sql.DB) *SalesRepository {
	return &SalesRepository{db: db}
}

func (r *SalesRepository) InsertSalesReport(userID, work, todaysWorkPlan string) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("user_id is required")
	}

	newEmpID := uuid.New().String()

	_, err := r.db.Exec(`
        INSERT INTO sales_reports (user_id, emp_id, work, todays_work_plan, login_time, created_at, report_date)
        VALUES ($1, $2, $3, $4, NOW(), NOW(), timezone('UTC', NOW())::DATE)
        ON CONFLICT (user_id, report_date) DO NOTHING
    `, userID, newEmpID, work, todaysWorkPlan)

	if err != nil {
		return "", fmt.Errorf("failed to insert sales report: %v", err)
	}

	// Retrieve emp_id
	var empID string
	err = r.db.QueryRow(`
        SELECT emp_id FROM sales_reports 
        WHERE user_id = $1 AND report_date = timezone('UTC', NOW())::DATE
    `, userID).Scan(&empID)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("emp_id not found after insertion")
		}
		return "", fmt.Errorf("failed to fetch emp_id: %v", err)
	}

	return empID, nil
}

func (r *SalesRepository) FetchSalesReport(userID string) (*models.SalesReport, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	var report models.SalesReport
	err := r.db.QueryRow(`
		SELECT user_id, emp_id, work, todays_work_plan, login_time, created_at, report_date 
		FROM sales_reports 
		WHERE user_id = $1 
		AND report_date = timezone('UTC', NOW())::DATE
	`, userID).Scan(
		&report.UserID, &report.EmployeeID, &report.Work,
		&report.TodaysWorkPlan, &report.LoginTime,
		&report.CreatedAt, &report.ReportDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no sales report found for today")
		}
		return nil, fmt.Errorf("failed to fetch sales report: %v", err)
	}

	return &report, nil
}

// Check if user has logged in today
func (r *SalesRepository) HasUserLoggedInToday(userID string) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("user_id is required")
	}

	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM sales_reports 
			WHERE user_id = $1 
			AND report_date = timezone('UTC', NOW())::DATE
		)
	`, userID).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check login status: %v", err)
	}

	return exists, nil
}

// Insert logout summary
func (r *SalesRepository) InsertLogoutSummary(
	userID, empID, customerFollowUpName, notes, tomorrowGoals, howWasToday, workLocation string,
	totalNoOfVisits, totalNoOfColdCalls, totalNoOfFollowUps, totalEnquiryGenerated, totalOrderLost, totalOrderWon int,
	totalEnquiryValue, totalOrderLostValue, totalOrderWonValue float64) error {

	if userID == "" || empID == "" {
		return fmt.Errorf("user_id and emp_id are required to log out")
	}

	_, err := r.db.Exec(`
        INSERT INTO logout_summaries (
            user_id, emp_id, total_no_of_visits, total_no_of_cold_calls, total_no_of_follow_ups,
            total_enquiry_generated, total_enquiry_value, total_order_lost, total_order_lost_value,
            total_order_won, total_order_won_value, customer_follow_up_name, notes, tomorrow_goals,
            how_was_today, work_location, logout_time, report_date
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, NOW(), timezone('UTC', NOW())::DATE
        )`,
		userID, empID, totalNoOfVisits, totalNoOfColdCalls, totalNoOfFollowUps,
		totalEnquiryGenerated, totalEnquiryValue, totalOrderLost, totalOrderLostValue,
		totalOrderWon, totalOrderWonValue, customerFollowUpName, notes, tomorrowGoals,
		howWasToday, workLocation,
	)

	if err != nil {
		return fmt.Errorf("failed to insert logout summary: %w", err)
	}

	return nil
}
func (r *SalesRepository) GetLogoutSummary(userID string) ([]models.LogoutSummary, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required to fetch logout summary")
	}

	rows, err := r.db.Query(`
		SELECT user_id, emp_id, total_no_of_visits, total_no_of_cold_calls, total_no_of_follow_ups,
		       total_enquiry_generated, total_enquiry_value, total_order_lost, total_order_lost_value,
		       total_order_won, total_order_won_value, customer_follow_up_name, notes, tomorrow_goals,
		       how_was_today, work_location, logout_time, report_date
		FROM logout_summaries
		WHERE user_id = $1
		ORDER BY logout_time DESC
	`, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch logout summary: %v", err)
	}

	defer rows.Close()

	var summaries []models.LogoutSummary
	for rows.Next() {
		var summary models.LogoutSummary
		err := rows.Scan(
			&summary.UserID, &summary.EmpID, &summary.TotalNoOfVisits, &summary.TotalNoOfColdCalls, &summary.TotalNoOfFollowUps,
			&summary.TotalEnquiryGenerated, &summary.TotalEnquiryValue, &summary.TotalOrderLost, &summary.TotalOrderLostValue,
			&summary.TotalOrderWon, &summary.TotalOrderWonValue, &summary.CustomerFollowUpName, &summary.Notes, &summary.TomorrowGoals,
			&summary.HowWasToday, &summary.WorkLocation, &summary.LogoutTime, &summary.ReportDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan logout summary: %v", err)
		}
		summaries = append(summaries, summary)
	}

	if len(summaries) == 0 {
		return nil, fmt.Errorf("no logout history found for user_id %s", userID)
	}

	return summaries, nil
}

func (r *SalesRepository) GetEmpIDByUserID(userID string) (string, error) {
	var empID string
	query :=
		`SELECT emp_id FROM sales_reports 
		WHERE user_id = $1 AND report_date = timezone('UTC', NOW())::DATE`

	err := r.db.QueryRow(query, userID).Scan(&empID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no sales report found for today, cannot log out")
		}
		return "", fmt.Errorf("failed to fetch emp_id: %v", err)
	}
	return empID, nil
}
