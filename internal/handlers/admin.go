package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Srujankm12/SRproject/repository"
)

type AdminHandler struct {
	repo *repository.Admin
}

func NewAdminHandler(repo *repository.Admin) *AdminHandler {
	return &AdminHandler{repo}
}

func (h *AdminHandler) AdminRegister(w http.ResponseWriter, r *http.Request) {

	err := h.repo.AdminRegisterM(r)
	if err != nil {

		http.Error(w, fmt.Sprintf("Admin registration failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Admin registered successfully"))
}

func (h *AdminHandler) AdminLogin(w http.ResponseWriter, r *http.Request) {

	adminID, err := h.repo.AdminLogin(r)
	if err != nil {

		http.Error(w, fmt.Sprintf("Admin login failed: %v", err), http.StatusUnauthorized)
		return
	}

	response := map[string]string{"admin_id": adminID}

	json.NewEncoder(w).Encode(response)
}
