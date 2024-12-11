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

func (q *Query) StoreFile(emp_id string, filename string, filedata []byte) error {
	_, err := q.db.Exec("INSERT INTO documents (emp_id,file_name,file_data)VALUES($1,$2)", filename, filedata)
	if err != nil {
		return err
	}
	return nil
}

func (q *Query) StoreFormData(empid string, data map[string]string) error {
	_, err := q.db.Exec(`
    INSERT INTO formdata(
        emp_id, 
        report_date, 
        employee_name, 
        premises, 
        site_location, 
        client_name, 
        scope_of_work, 
        work_details, 
        joint_visits, 
        support_needed, 
        status_of_work, 
        priority_of_work, 
        next_action_plan, 
        result, 
        type_of_work, 
        closing_time, 
        contact_person_name
    ) 
    VALUES(
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
    )`,
		empid,
		data["report_date"],
		data["employee_name"],
		data["premises"],
		data["site_location"],
		data["client_name"],
		data["scope_of_work"],
		data["work_details"],
		data["joint_visits"],
		data["support_needed"],
		data["status_of_work"],
		data["priority_of_work"],
		data["next_action_plan"],
		data["result"],
		data["type_of_work"],
		data["closing_time"],
		data["contact_person_name"],
	)
	if err != nil {
		return err
	}
	return nil
}
