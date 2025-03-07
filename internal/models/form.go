package models

import "net/http"

type FormData struct {
	UserID         string `json:"user_id"`
	EmployeeID     string `json:"emp_id"`
	ReportDate     string `json:"report_date"`
	EmployeeName   string `json:"employee_name"`
	Premises       string `json:"premises"`
	SiteLocation   string `json:"site_location"`
	ClientName     string `json:"client_name"`
	ScopeOfWork    string `json:"scope_of_work"`
	WorkDetails    string `json:"work_details"`
	JointVisits    string `json:"joint_visits"`
	SupportNeeded  string `json:"support_needed"`
	StatusOfWork   string `json:"status_of_work"`
	PriorityOfWork string `json:"priority_of_work"`
	NextActionPlan string `json:"next_action_plan"`
	// Result            string `json:"result"`
	TypeOfWork        string `json:"type_of_work"`
	ClosingTime       string `json:"closing_time"`
	ContactPersonName string `json:"contact_person_name"`
	ContactEmailID    string `json:"contact_emailid"`
}

type FormInterface interface {
	SubmitFormData(r *http.Request) error
	FetchFormData(r *http.Request) ([]FormData, error)
}
