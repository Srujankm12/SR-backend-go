package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Srujankm12/SRproject/repository"
	"github.com/gorilla/mux"
)

type AdminFHandler struct {
	Repo *repository.AdminRepository
}

func NewAdminFHandler(repo *repository.AdminRepository) *AdminFHandler {
	return &AdminFHandler{
		Repo: repo,
	}
}

// HandleAdminFetchFormData fetches all form data for the admin
func (h *AdminFHandler) HandleAdminFetchFormData(w http.ResponseWriter, r *http.Request) {
	formData, err := h.Repo.FetchAllFormData()
	if err != nil {
		http.Error(w, "Failed to fetch form data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(formData)
}
func (h *AdminFHandler) HandleDeleteEmployee(w http.ResponseWriter, r *http.Request) {
	// Get the emp_id from the URL path parameter
	vars := mux.Vars(r)
	empID := vars["id"]

	// Ensure emp_id is provided
	if empID == "" {
		http.Error(w, "Employee ID is required", http.StatusBadRequest)
		return
	}

	// Call the repository method to delete the employee
	err := h.Repo.DeleteEmployee(empID)
	if err != nil {
		http.Error(w, "Failed to delete employee", http.StatusInternalServerError)
		return
	}

	// Optionally return a message with a 200 OK response

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Employee successfully deleted",
	})
}
