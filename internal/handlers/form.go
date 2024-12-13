package handlers

import (
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

// SubmitFormController handles the form submission, calling the ReportForm method of formRepo
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

// // FetchFormDataController handles fetching form data and associated files
// func (fc *FormController) FetchFormDataController(w http.ResponseWriter, r *http.Request) {
// 	// Call the FetchFormData method of formRepo
// 	formData, file1Data, file2Data, err := fc.formRepo.FetchFormData(r)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		utils.Encode(w, map[string]string{"message": err.Error()})
// 		return
// 	}

// 	// Create a response structure
// 	response := map[string]interface{}{
// 		"formData": formData,
// 		"file1":    file1Data,
// 		"file2":    file2Data,
// 	}

// 	// Respond with the form data and files
// 	w.WriteHeader(http.StatusOK)
// 	utils.Encode(w, response)
// }

