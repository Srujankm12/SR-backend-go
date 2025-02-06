package models

import "net/http"

type Adminmodel struct {
	AdminID       string `json:"admin_id"`
	AdminEmail    string `json:"admin_email"`
	AdminPassword string `json:"admin_password"`
}

type AdminInterface interface {
	AdminFetchFormData(r *http.Request) ([]AdminFetchFormData, error)
	AdminRegisterM(r *http.Request) error
	AdminLoginM(r *http.Request) (string, error)
}
