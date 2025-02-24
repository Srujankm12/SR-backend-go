// repository/admin_repository.go

package repository

import (
	"database/sql"

	"github.com/Srujankm12/SRproject/internal/models"
)

type AdminRepository struct {
	db *sql.DB
}

// Constructor function for AdminRepository
func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{
		db: db,
	}
}

// Fetch all submitted form data for admin
func (repo *AdminRepository) FetchAllFormData() ([]models.AdminFetchFormData, error) {
	var formDatas []models.AdminFetchFormData

	rows, err := repo.db.Query(`
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
		var formData models.AdminFetchFormData
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

// DeleteEmployee deletes the employee from formdata and documents based on emp_id
func (repo *AdminRepository) DeleteEmployee(empId string) error {
	tx, err := repo.db.Begin()
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
