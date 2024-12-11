package models

import "net/http"

type UserModel struct {
	UserID          string `json:"user_id"`
	Email           string `json:"email"`
	ConfirmPassword string `json:"confirm_password"`
	Password        string `json:"password"`
}

type AuthInterface interface {
	Register(*http.Request) error
	Login(*http.Request) (string, error)
}
