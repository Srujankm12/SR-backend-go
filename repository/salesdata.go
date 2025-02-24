package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

type SalesRepository struct {
	db *sql.DB
}

func NewSalesRepository(db *sql.DB) *SalesRepository {
	return &SalesRepository{db: db}
}

func (r *SalesRepository) GenerateEmpID(userID, todaysWorkPlan string) (string, error) {
	if userID == "" {
		log.Println("ERROR: Empty userID received in GenerateEmpID")
		return "", fmt.Errorf("invalid userID: empty value")
	}

	userID = strings.TrimSpace(userID)

	var exists bool
	err := r.db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE user_id = $1)", userID).Scan(&exists)
	if err != nil {
		log.Println("Database error checking user existence:", err)
		return "", fmt.Errorf("failed to check user existence: %w", err)
	}
	if !exists {
		log.Println("ERROR: userID does not exist in users table")
		return "", fmt.Errorf("userID does not exist")
	}

	var empID string
	log.Println("Checking for existing empID for user:", userID)

	err = r.db.QueryRow(`
		SELECT emp_id FROM login_logout_report 
		WHERE user_id = $1 AND DATE(login_time) = CURRENT_DATE
	`, userID).Scan(&empID)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No existing empID found, creating new one")

			if todaysWorkPlan == "" {
				return "", fmt.Errorf("todaysWorkPlan cannot be empty")
			}

			newEmpID := uuid.New().String()

			_, insertErr := r.db.Exec(`
				INSERT INTO login_logout_report (user_id, emp_id, login_time, work_location, todays_work_plan)
				VALUES ($1, $2, NOW(), 'Office', $3)
			`, userID, newEmpID, todaysWorkPlan)

			if insertErr != nil {
				log.Println("Database insert error:", insertErr)
				return "", fmt.Errorf("failed to insert empID: %w", insertErr)
			}
			return newEmpID, nil
		}
		return "", fmt.Errorf("failed to check empID: %w", err)
	}

	return empID, nil
}

func (r *SalesRepository) UserLogin(userID, workLocation, todaysWorkPlan string) (string, error) {
	userID = strings.TrimSpace(userID)

	if todaysWorkPlan == "" {
		return "", fmt.Errorf("todaysWorkPlan cannot be empty")
	}

	empID, err := r.GenerateEmpID(userID, todaysWorkPlan)
	if err != nil {
		return "", err
	}

	_, err = r.db.Exec(`
		UPDATE login_logout_report 
		SET todays_work_plan = $1, work_location = $2
		WHERE user_id = $3 AND emp_id = $4 AND DATE(login_time) = CURRENT_DATE
	`, todaysWorkPlan, workLocation, userID, empID)

	if err != nil {
		return "", fmt.Errorf("failed to log in: %w", err)
	}

	return empID, nil
}

func (r *SalesRepository) UserLogout(userID, empID string, totalVisits, coldCalls, followUps, enquiries, ordersLost, ordersWon int, notes, tomorrowGoals, howWasToday string) error {
	userID = strings.TrimSpace(userID)
	empID = strings.TrimSpace(empID)

	_, err := r.db.Exec(`
		UPDATE login_logout_report 
		SET logout_time = NOW(),
		    total_no_of_visits = $1,
		    total_no_of_cold_calls = $2,
		    total_no_of_customer_follow_up = $3,
		    total_enquiry_generated = $4,
		    total_order_lost = $5,
		    total_order_won = $6,
		    notes = $7,
		    tomorrow_goals = $8,
		    how_was_today = $9
		WHERE user_id = $10 AND emp_id = $11 AND DATE(login_time) = CURRENT_DATE
	`, totalVisits, coldCalls, followUps, enquiries, ordersLost, ordersWon, notes, tomorrowGoals, howWasToday, userID, empID)

	if err != nil {
		return fmt.Errorf("failed to log out: %w", err)
	}
	return nil
}

func (r *SalesRepository) SiteCheckIn(userID, empID, companyName, purpose string) error {
	userID = strings.TrimSpace(userID)
	empID = strings.TrimSpace(empID)

	_, err := r.db.Exec(`
		INSERT INTO site_checkin_checkout_report (user_id, emp_id, checkin_time, company_name, purpose)
		VALUES ($1, $2, NOW(), $3, $4)
	`, userID, empID, companyName, purpose)

	if err != nil {
		return fmt.Errorf("failed to check in: %w", err)
	}
	return nil
}

func (r *SalesRepository) SiteCheckOut(userID, empID, engineerName, companySalesStage, visitOn, timelineForNextActionPlan, challenges, resultOfVisit, notes string, visitRating int) error {
	userID = strings.TrimSpace(userID)
	empID = strings.TrimSpace(empID)

	_, err := r.db.Exec(`
		UPDATE site_checkin_checkout_report 
		SET checkout_time = NOW(),
		    engineer_name = $1,
		    company_sales_stage = $2,
		    visit_on = $3,
		    timeline_for_next_action_plan = $4,
		    challenges = $5,
		    visit_rating = $6,
		    result_of_visit = $7,
		    notes = $8
		WHERE user_id = $9 AND emp_id = $10 AND checkout_time IS NULL
	`, engineerName, companySalesStage, visitOn, timelineForNextActionPlan, challenges, visitRating, resultOfVisit, notes, userID, empID)

	if err != nil {
		return fmt.Errorf("failed to check out: %w", err)
	}
	return nil
}
