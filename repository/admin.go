package repository

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Srujankm12/SRproject/internal/models"
	"github.com/Srujankm12/SRproject/pkg/database"
	"github.com/Srujankm12/SRproject/pkg/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	db *sql.DB
}

func NewAdmin(db *sql.DB) *Admin {
	return &Admin{
		db,
	}
}

func (a *Admin) AdminRegisterM(r *http.Request) error {
	var newAdmin models.Adminmodel
	if err := utils.Decode(r, &newAdmin); err != nil {
		return err
	}
	if newAdmin.AdminPassword != newAdmin.AdminPassword {
		return fmt.Errorf("passwords do not match")
	}

	newAdmin.AdminID = uuid.NewString()

	hash, err := bcrypt.GenerateFromPassword([]byte(newAdmin.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := database.NewQuery(a.db)
	err = query.AdminRegister(newAdmin.AdminID, newAdmin.AdminEmail, string(hash))
	if err != nil {
		return err
	}
	return nil
}

func (a *Admin) AdminLogin(r *http.Request) (string, error) {
	var reqDet models.Adminmodel
	if err := utils.Decode(r, &reqDet); err != nil {
		return "", err
	}
	log.Println(reqDet)

	query := database.NewQuery(a.db)
	admin, err := query.RetriveAdminPassowrd(reqDet.AdminEmail)
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.AdminPassword), []byte(reqDet.AdminPassword))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return admin.AdminID, nil
}
