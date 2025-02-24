package handlers

import (
	"log"
	"net/http"

	"github.com/Srujankm12/SRproject/internal/models"
	"github.com/Srujankm12/SRproject/pkg/utils"
)

type FormController struct {
	formRepo models.FormInterface
}

func NewFormController(formRepo models.FormInterface) *FormController {
	return &FormController{
		formRepo,
	}
}

func (fc *FormController) SubmitFormController(w http.ResponseWriter, r *http.Request) {
	err := fc.formRepo.SubmitFormData(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, map[string]string{"message": "success"})
}

func (fc *FormController) FetchFormDataController(w http.ResponseWriter, r *http.Request) {
	log.Println("FetchFormDataController called") // <-- Debug log

	formData, err := fc.formRepo.FetchFormData(r)
	if err != nil {
		log.Println("Error in FetchFormData:", err) // <-- Log error
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": err.Error()})
		return
	}

	log.Println("Fetched form data successfully") // <-- Confirm data was fetched
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, formData)
}
