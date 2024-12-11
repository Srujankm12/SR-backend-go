package handlers

import (
	"net/http"

	"github.com/Srujankm12/SRproject/internal/models"
	"github.com/Srujankm12/SRproject/pkg/utils"
)

type AuthController struct {
	authRepo models.AuthInterface
}

func NewAuthController(authRepo models.AuthInterface) *AuthController {
	return &AuthController{
		authRepo,
	}
}
func (ac *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	err := ac.authRepo.Register(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, map[string]string{"message": "success"})
}
func (ac *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	res, err := ac.authRepo.Login(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Encode(w, map[string]string{"message": res})
}
