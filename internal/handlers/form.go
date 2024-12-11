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
func (fc *FormController) SubmitFormController(w http.ResponseWriter, r *http.Request) {
	err := fc.formRepo.ReportForm(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, map[string]string{"message": "success"})

}
