package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Srujankm12/SRproject/internal/models"
)

type Query struct {
	db *sql.DB
}

func NewQuery(db *sql.DB) *Query {
	return &Query{
		db,
	}
}

func (q *Query) CreateTables() error {
	tx, err := q.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
    		user_id VARCHAR(100) PRIMARY KEY,
    		email VARCHAR(255) NOT NULL,
    		password VARCHAR(100) NOT NULL
		)`,
		`
		CREATE TABLE IF NOT EXISTS documents (
    		emp_id VARCHAR(100) PRIMARY KEY,
    		user_id VARCHAR(100),
    		file_name_one VARCHAR(255) NOT NULL,
    		file_data_one BYTEA NOT NULL,
			file_name_two VARCHAR(255) NOT NULL,
    		file_data_two BYTEA NOT NULL,
    		FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		)
		`,
		`	CREATE TABLE IF NOT EXISTS formdata (
    			user_id VARCHAR(100) NOT NULL,
    			emp_id VARCHAR(100) NOT NULL,
    			report_date DATE NOT NULL,
    			employee_name VARCHAR(255) NOT NULL,
    			premises VARCHAR(255) NOT NULL,
    			site_location VARCHAR(255) NOT NULL,
    			client_name VARCHAR(255) NOT NULL,
    			scope_of_work TEXT,
    			work_details TEXT,
    			joint_visits VARCHAR(255),
    			support_needed VARCHAR(255),
    			status_of_work VARCHAR(255),
    			priority_of_work VARCHAR(255),
    			next_action_plan TEXT,
    			type_of_work VARCHAR(255),
    			closing_time VARCHAR(200),
    			contact_person_name VARCHAR(255),
    			contact_emailid VARCHAR(255),
    			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
				FOREIGN KEY (emp_id) REFERENCES documents(emp_id) ON DELETE CASCADE
				)`,
		`CREATE TABLE IF NOT EXISTS admin (
    		admin_id VARCHAR(100) PRIMARY KEY,
    		admin_email VARCHAR(255) NOT NULL,  
    		admin_password VARCHAR(100) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sales_reports (
			user_id VARCHAR(100) NOT NULL,
			emp_id VARCHAR(100) UNIQUE DEFAULT NULL,
			work TEXT NOT NULL,
			todays_work_plan TEXT NOT NULL,
			login_time TIMESTAMP NOT NULL DEFAULT NOW(),
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			report_date DATE NOT NULL DEFAULT CURRENT_DATE, 
			PRIMARY KEY (user_id, report_date)	



	)`,
		// `CREATE UNIQUE INDEX unique_user_login_per_day ON sales_reports (user_id, DATE(login_time))`,

		`CREATE TABLE IF NOT EXISTS logout_summaries (
			user_id VARCHAR(100) NOT NULL,
			emp_id VARCHAR(100) NOT NULL,
			total_no_of_visits INT DEFAULT 0,
			total_no_of_cold_calls INT DEFAULT 0,
			total_no_of_follow_ups INT DEFAULT 0,
			total_enquiry_generated INT DEFAULT 0,
			total_enquiry_value NUMERIC(10,2) DEFAULT 0,
			total_order_lost INT DEFAULT 0,
			total_order_lost_value NUMERIC(10,2) DEFAULT 0,
			total_order_won INT DEFAULT 0,
			total_order_won_value NUMERIC(10,2) DEFAULT 0,
			customer_follow_up_name TEXT,
			notes TEXT,
			tomorrow_goals TEXT,
			how_was_today TEXT,
			work_location TEXT,
			logout_time TIMESTAMP NOT NULL DEFAULT NOW(),
			report_date DATE NOT NULL DEFAULT CURRENT_DATE,
			PRIMARY KEY (user_id, report_date),
			FOREIGN KEY (user_id, report_date) 
			REFERENCES sales_reports(user_id, report_date) 
			ON DELETE CASCADE

		)`,
	}

	for _, query := range queries {
		_, err = tx.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}
func (q *Query) InsertSalesReport(userID, empID, work, todaysWorkPlan string) error {
	if userID == "" || empID == "" {
		return fmt.Errorf("user_id or emp_id is empty, cannot insert sales report")
	}

	_, err := q.db.Exec(`
		INSERT INTO sales_reports (user_id, emp_id, work, todays_work_plan, login_time, created_at, report_date)
		VALUES ($1, $2, $3, $4, NOW(), NOW(), CURRENT_DATE)
		ON CONFLICT (user_id, report_date) 
		DO UPDATE SET 
			work = EXCLUDED.work, 
			todays_work_plan = EXCLUDED.todays_work_plan,
			emp_id = COALESCE(sales_reports.emp_id, EXCLUDED.emp_id)
	`, userID, empID, work, todaysWorkPlan)

	if err != nil {
		return fmt.Errorf("failed to insert/update sales report: %v", err)
	}
	return nil
}

func (q *Query) GetSalesReport(userID string) ([]models.SalesReport, error) {
	rows, err := q.db.Query(`
		SELECT user_id, emp_id, work, todays_work_plan, login_time, created_at, report_date 
		FROM sales_reports 
		WHERE user_id = $1 
		AND report_date = CURRENT_DATE
	`, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch sales reports: %v", err)
	}
	defer rows.Close()

	var reports []models.SalesReport
	for rows.Next() {
		var report models.SalesReport
		err := rows.Scan(&report.UserID, &report.EmployeeID, &report.Work, &report.TodaysWorkPlan, &report.LoginTime, &report.CreatedAt, &report.ReportDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sales report: %v", err)
		}
		reports = append(reports, report)
	}
	return reports, nil
}

func (q *Query) InsertLogoutSummary(
	userID, empID, customerFollowUpName, notes, tomorrowGoals, howWasToday, workLocation string,
	totalNoOfVisits, totalNoOfColdCalls, totalNoOfFollowUps, totalEnquiryGenerated, totalOrderLost, totalOrderWon int,
	totalEnquiryValue, totalOrderLostValue, totalOrderWonValue float64) error {

	// Step 1: Check if the user has already logged out today
	var exists bool
	err := q.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM logout_summaries 
			WHERE user_id = $1 AND report_date = CURRENT_DATE
		)
	`, userID).Scan(&exists)

	if err != nil {
		return fmt.Errorf("error checking logout existence: %v", err)
	}

	// Step 2: Prevent duplicate insertions
	if exists {
		return fmt.Errorf("user has already logged out today")
	}

	// Step 3: Insert new logout record if none exists
	_, err = q.db.Exec(`
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
func (q *Query) GetLogoutSummary(userID string) ([]models.LogoutSummary, error) {
	log.Printf("Fetching logout summary for user_id: %s", userID)

	rows, err := q.db.Query(`
	SELECT user_id, emp_id, total_no_of_visits, total_no_of_cold_calls, total_no_of_follow_ups,
	       total_enquiry_generated, total_enquiry_value, total_order_lost, total_order_lost_value,
	       total_order_won, total_order_won_value, customer_follow_up_name, notes, tomorrow_goals,
	       how_was_today, work_location, logout_time, report_date
	FROM logout_summaries
	WHERE user_id = $1
	AND DATE(report_date) = CURRENT_DATE
	ORDER BY logout_time DESC
	LIMIT 1
`, userID)

	if err != nil {
		log.Printf("Query failed: %v", err) // Log query failure
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
			log.Printf("Row scan failed: %v", err) // Log scan failure
			return nil, fmt.Errorf("failed to scan logout summary: %v", err)
		}
		summaries = append(summaries, summary)
	}

	if len(summaries) == 0 {
		log.Printf("No logout record found for user_id: %s", userID)
		return nil, fmt.Errorf("no logout history found for user_id %s today", userID)
	}

	log.Printf("Fetched logout summary successfully for user_id: %s", userID)
	return summaries, nil
}

func (q *Query) Register(userid, email, password string) error {
	if _, err := q.db.Exec("INSERT INTO users(user_id,email,password) VALUES($1,$2,$3)", userid, email, password); err != nil {
		return err
	}
	return nil
}
func (q *Query) Login(email string) (models.UserModel, error) {
	var user models.UserModel
	if err := q.db.QueryRow("SELECT user_id,email,password FROM users WHERE email = $1", email).Scan(&user.UserID, &user.Email, &user.Password); err != nil {
		return user, err
	}
	return user, nil
}
func (q *Query) AdminRegister(adminId, adminEmail, adminPassword string) error {
	if _, err := q.db.Exec("INSERT INTO admin(admin_id,admin_email,admin_password) VALUES($1,$2,$3)", adminId, adminEmail, adminPassword); err != nil {
		return err
	}
	return nil
}
func (q *Query) AdminLogin(adminEmail string) (models.Adminmodel, error) {
	var admin models.Adminmodel
	if err := q.db.QueryRow("SELECT admin_id,admin_email,admin_password FROM admin WHERE admin_email = $1", adminEmail).Scan(&admin.AdminID, &admin.AdminEmail, &admin.AdminPassword); err != nil {
		return admin, err
	}
	return admin, nil
}
func (q *Query) RetriveAdminPassowrd(adminEmail string) (models.Adminmodel, error) {
	var admin models.Adminmodel
	if err := q.db.QueryRow("SELECT admin_id,admin_email,admin_password FROM admin WHERE admin_email = $1", adminEmail).Scan(&admin.AdminID, &admin.AdminEmail, &admin.AdminPassword); err != nil {
		return admin, err
	}
	return admin, nil
}

func (q *Query) RetrivePassword(email string) (models.UserModel, error) {
	var user models.UserModel
	if err := q.db.QueryRow("SELECT user_id,password FROM users WHERE email = $1", email).Scan(&user.UserID, &user.Password); err != nil {
		return user, err
	}
	return user, nil
}

func (q *Query) StoreFile(userId, empId, fileNameOne, fileNameTwo string, fileDataOne, fileDataTwo []byte) error {
	_, err := q.db.Exec("INSERT INTO documents (user_id , emp_id , file_name_one , file_data_one , file_name_two , file_data_two) VALUES ($1, $2, $3,$4,$5,$6)", userId, empId, fileNameOne, fileDataOne, fileNameTwo, fileDataTwo)
	if err != nil {
		return err
	}
	return nil
}

func (q *Query) StoreFormData(data models.FormData) error {
	_, err := q.db.Exec(`
		INSERT INTO formdata (
			user_id, emp_id, report_date, employee_name, premises, 
			site_location, client_name, scope_of_work, work_details, joint_visits, 
			support_needed, status_of_work, priority_of_work, next_action_plan, 
		 type_of_work, closing_time, contact_person_name, contact_emailid
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		)`,
		data.UserID,
		data.EmployeeID,
		data.ReportDate,
		data.EmployeeName,
		data.Premises,
		data.SiteLocation,
		data.ClientName,
		data.ScopeOfWork,
		data.WorkDetails,
		data.JointVisits,
		data.SupportNeeded,
		data.StatusOfWork,
		data.PriorityOfWork,
		data.NextActionPlan,

		data.TypeOfWork,
		data.ClosingTime,
		data.ContactPersonName,
		data.ContactEmailID,
	)
	if err != nil {
		return err
	}
	return nil
}
func (q *Query) FetchFormData(userID string) ([]models.FormData, error) {
	var formDatas []models.FormData

	// Debugging: Ensure `userID` is being used
	log.Println("Executing query for userID:", userID)

	rows, err := q.db.Query(`
		SELECT 
			user_id, emp_id, report_date, employee_name, premises, site_location, client_name,
			scope_of_work, work_details, joint_visits, support_needed, status_of_work, 
			priority_of_work, next_action_plan, type_of_work, closing_time, 
			contact_person_name, contact_emailid
		FROM formdata
		WHERE user_id = $1`, userID) // ✅ Filtering by user_id
	if err != nil {
		log.Println("Database query error:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var formData models.FormData
		err := rows.Scan(
			&formData.UserID, &formData.EmployeeID, &formData.ReportDate, &formData.EmployeeName,
			&formData.Premises, &formData.SiteLocation, &formData.ClientName, &formData.ScopeOfWork,
			&formData.WorkDetails, &formData.JointVisits, &formData.SupportNeeded, &formData.StatusOfWork,
			&formData.PriorityOfWork, &formData.NextActionPlan, &formData.TypeOfWork,
			&formData.ClosingTime, &formData.ContactPersonName, &formData.ContactEmailID,
		)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		formDatas = append(formDatas, formData)
	}

	if err = rows.Err(); err != nil {
		log.Println("Row iteration error:", err)
		return nil, err
	}

	log.Println("Fetched data for userID:", userID, formDatas)
	return formDatas, nil
}
func (q *Query) AdminFetchFormData() ([]models.FormData, error) {
	var formDatas []models.FormData

	rows, err := q.db.Query(`
		SELECT 
			user_id, emp_id, report_date, employee_name, premises, site_location, client_name,
			scope_of_work, work_details, joint_visits, support_needed, status_of_work, 
			priority_of_work, next_action_plan, type_of_work, closing_time, 
			contact_person_name, contact_emailid
		FROM formdata
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var formData models.FormData
		err := rows.Scan(
			&formData.UserID, &formData.EmployeeID, &formData.ReportDate, &formData.EmployeeName,
			&formData.Premises, &formData.SiteLocation, &formData.ClientName, &formData.ScopeOfWork,
			&formData.WorkDetails, &formData.JointVisits, &formData.SupportNeeded, &formData.StatusOfWork,
			&formData.PriorityOfWork, &formData.NextActionPlan, &formData.TypeOfWork,
			&formData.ClosingTime, &formData.ContactPersonName, &formData.ContactEmailID,
		)
		if err != nil {
			return nil, err
		}
		formDatas = append(formDatas, formData)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return formDatas, nil
}
func (q *Query) DeleteEmployee(empId string) error {
	tx, err := q.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = tx.Exec("DELETE FROM formdata WHERE emp_id = $1", empId)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM documents WHERE emp_id = $1", empId)
	if err != nil {
		return err
	}

	return nil
}

func (q *Query) FetchExcel() ([]models.DownloadExcel, error) {
	var data []models.DownloadExcel
	rows, err := q.db.Query("SELECT * FROM formdata")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var record models.DownloadExcel
		if err := rows.Scan(&record.UserID, &record.EmployeeID, &record.ReportDate, &record.EmployeeName, &record.Premises, &record.SiteLocation, &record.ClientName, &record.ScopeOfWork, &record.WorkDetails, &record.JointVisits, &record.SupportNeeded, &record.StatusOfWork, &record.PriorityOfWork, &record.NextActionPlan, &record.TypeOfWork, &record.ClosingTime, &record.ContactPersonName, &record.ContactEmailID); err != nil {
			return nil, err
		}
		data = append(data, record)
	}
	if len(data) == 0 {
		return nil, errors.New("no data found")
	}
	return data, nil
}

// func (q *Query) FetchFile(, filename string) ([]byte, error) {
// 	var fileData []byte

// 	// Query the database to fetch the file data based on emp_id and file_name
// 	row := q.db.QueryRow(`
// 		SELECT file_data
// 		FROM documents
// 		WHERE emp_id = $1 AND file_name = $2
// 	`, empid, filename)

// 	// Scan the result into the fileData variable
// 	err := row.Scan(&fileData)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, errors.New("file not found") // No file found for this emp_id and file_name
// 		}
// 		return nil, err // Other errors
// 	}

// 	return fileData, nil
// }
