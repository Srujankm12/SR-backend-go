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

	_, err := r.db.Exec(`
        INSERT INTO sales_reports (user_id, work, todays_work_plan, login_time, created_at, report_date)
        VALUES ($1, $2, $3, NOW(), NOW(), CURRENT_DATE)
        ON CONFLICT (user_id, report_date) DO NOTHING
    `, userID, work, todaysWorkPlan)

	if err != nil {
		return "", fmt.Errorf("failed to insert sales report: %v", err)
	}

	var empID sql.NullString
	err = r.db.QueryRow(`
        SELECT emp_id FROM sales_reports WHERE user_id = $1 AND report_date = CURRENT_DATE
    `, userID).Scan(&empID)

	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to fetch emp_id: %v", err)
	}

	if !empID.Valid {

		newEmpID := uuid.New().String()
		_, err = r.db.Exec(`
            UPDATE sales_reports SET emp_id = $1 WHERE user_id = $2 AND report_date = CURRENT_DATE
        `, newEmpID, userID)

		if err != nil {
			return "", fmt.Errorf("failed to assign emp_id: %v", err)
		}
		return newEmpID, nil
	}

	return empID.String, nil
}

func (r *SalesRepository) GetSalesReport(userID string) (*models.SalesReport, error) {
	var report models.SalesReport

	err := r.db.QueryRow(`
		SELECT user_id, emp_id, work, todays_work_plan, login_time, created_at 
		FROM sales_reports 
		WHERE user_id = $1 AND report_date = CURRENT_DATE
	`, userID).Scan(
		&report.UserID, &report.EmployeeID, &report.Work, &report.TodaysWorkPlan,
		&report.LoginTime, &report.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sales report not found")
		}
		return nil, fmt.Errorf("failed to fetch sales report: %v", err)
	}

	return &report, nil
}

func (r *SalesRepository) GetEmpIDByUserID(userID string) (string, error) {
	var empID string
	query := `
		SELECT emp_id FROM sales_reports 
		WHERE user_id = $1 AND report_date = (SELECT CURRENT_DATE AT TIME ZONE 'UTC')
	`
	err := r.db.QueryRow(query, userID).Scan(&empID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no sales report found for today, cannot log out")
		}
		return "", fmt.Errorf("failed to fetch emp_id: %v", err)
	}
	return empID, nil
}

func (r *SalesRepository) UpdateSalesReport(userID, work, todaysWorkPlan string) error {
	_, err := r.db.Exec(`
		UPDATE sales_reports 
		SET work = $1, todays_work_plan = $2 
		WHERE user_id = $3 AND report_date = CURRENT_DATE
	`, work, todaysWorkPlan, userID)

	if err != nil {
		return fmt.Errorf("failed to update sales report: %v", err)
	}

	return nil
}

func (r *SalesRepository) HasUserLoggedInToday(userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM sales_reports WHERE user_id = $1 AND report_date = CURRENT_DATE
		)
	`, userID).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check login status: %v", err)
	}

	return exists, nil
}

func (r *SalesRepository) InsertLogoutSummary(
	userID, empID, customerFollowUpName, notes, tomorrowGoals, howWasToday, workLocation string,
	totalNoOfVisits, totalNoOfColdCalls, totalNoOfFollowUps, totalEnquiryGenerated, totalOrderLost, totalOrderWon int,
	totalEnquiryValue, totalOrderLostValue, totalOrderWonValue float64) error {

	if empID == "" {
		return fmt.Errorf("emp_id is required to log out")
	}

	_, err := r.db.Exec(`
       INSERT INTO logout_summaries (
    user_id, emp_id, total_no_of_visits, total_no_of_cold_calls, total_no_of_follow_ups,
    total_enquiry_generated, total_enquiry_value, total_order_lost, total_order_lost_value,
    total_order_won, total_order_won_value, customer_follow_up_name, notes, tomorrow_goals, 
    how_was_today, work_location, logout_time, report_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, NOW(), CURRENT_DATE
)

    `, userID, empID, totalNoOfVisits, totalNoOfColdCalls, totalNoOfFollowUps,
		totalEnquiryGenerated, totalEnquiryValue, totalOrderLost, totalOrderLostValue,
		totalOrderWon, totalOrderWonValue, customerFollowUpName, notes, tomorrowGoals,
		howWasToday, workLocation)

	if err != nil {
		return fmt.Errorf("failed to insert logout summary: %v", err)
	}

	return nil
}
