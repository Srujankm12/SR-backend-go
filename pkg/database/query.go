package database

import (
	"database/sql"

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
		`
		CREATE TABLE IF NOT EXISTS formdata (
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
    			result TEXT,
    			type_of_work VARCHAR(255),
    			closing_time VARCHAR(200),
    			contact_person_name VARCHAR(255),
    			contact_emailid VARCHAR(255),
    			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
				FOREIGN KEY (emp_id) REFERENCES documents(emp_id) ON DELETE CASCADE
	)
		`,
	}

	for _ , j := range queries{
		_ ,err = tx.Exec(j)
		if err != nil {
			return err
		}
	}
	return nil
}

func (q *Query) Register(userid, email, password string) error {
	if _, err := q.db.Exec("INSERT INTO users(user_id,email,password) VALUES($1,$2,$3)", userid, email, password); err != nil {
		return err
	}
	return nil
}

func (q *Query) RetrivePassword(email string) (models.UserModel, error) {
	var user models.UserModel
	if err := q.db.QueryRow("SELECT user_id,password FROM users WHERE email = $1", email).Scan(&user.UserID, &user.Password); err != nil {
		return user, err
	}
	return user, nil
}

// StoreFile stores file data for the given employee ID and file name
func (q *Query) StoreFile(userId , empId  , fileNameOne , fileNameTwo string , fileDataOne , fileDataTwo []byte) error {
	_, err := q.db.Exec("INSERT INTO documents (user_id , emp_id , file_name_one , file_data_one , file_name_two , file_data_two) VALUES ($1, $2, $3,$4,$5,$6)",userId , empId , fileNameOne, fileDataOne ,  fileNameTwo , fileDataTwo)
	if err != nil {
		return err
	}
	return nil
}

// StoreFormData stores the form data associated with the given employee ID
func (q *Query) StoreFormData(data models.FormData) error {
	_, err := q.db.Exec(`
		INSERT INTO formdata (
			user_id, emp_id, report_date, employee_name, premises, 
			site_location, client_name, scope_of_work, work_details, joint_visits, 
			support_needed, status_of_work, priority_of_work, next_action_plan, 
			result, type_of_work, closing_time, contact_person_name, contact_emailid
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
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
		data.Result,
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

func (q *Query) FetchFormData(userId string) ([]models.FormData, error) {
	// Prepare a map to store the form data
	var formData models.FormData
	var formDatas []models.FormData

	// Query the database to fetch form data
	row, err := q.db.Query(`
		SELECT 
			user_id , emp_id , report_date, employee_name, premises, site_location, client_name,
			scope_of_work, work_details, joint_visits, support_needed, status_of_work, 
			priority_of_work, next_action_plan, result, type_of_work, closing_time, 
			contact_person_name, contact_emailid , type_of_work , closing_time , contact_emailid
		FROM formdata
		WHERE user_id = $1
	`, userId)

	defer func() {
		row.Close()
	} ()

	if err != nil {
		return nil, err
	}

	for row.Next() {
		if err := row.Scan(&formData.UserID , &formData.EmployeeID , &formData.ReportDate , &formData.EmployeeName , &formData.Premises , 
			&formData.SiteLocation , &formData.ClientName , &formData.ScopeOfWork,&formData.WorkDetails,&formData.JointVisits,
			&formData.SupportNeeded , &formData.StatusOfWork , &formData.PriorityOfWork,
			&formData.NextActionPlan , &formData.Result , &formData.TypeOfWork , &formData.ClosingTime,
			&formData.ContactPersonName , &formData.ContactEmailID , &formData.TypeOfWork , &formData.ClosingTime , &formData.ContactEmailID); err != nil {
			return nil,err
		}
		formDatas = append(formDatas, formData)
	}
	if row.Err() != nil {
		return nil,err
	}

	return formDatas, nil
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
