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

type Auth struct {
	db *sql.DB
}

func NewAuth(db *sql.DB) *Auth {
	return &Auth{
		db,
	}
}

func (a *Auth) Register(r *http.Request) error {
	var newUser models.UserModel
	if err := utils.Decode(r, &newUser); err != nil {
		return nil
	}
	if newUser.Password != newUser.ConfirmPassword {
		return fmt.Errorf("please enter valid password")
	}
	newUser.UserID = uuid.NewString()
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := database.NewQuery(a.db)
	err = query.Register(newUser.UserID, newUser.Email, string(hash))
	if err != nil {
		return err
	}

	return nil
}

func (a *Auth) Login(r *http.Request) (string, error) {
	var reqDet models.UserModel
	if err := utils.Decode(r, &reqDet); err != nil {
		return "", nil
	}
	log.Println(reqDet)
	query := database.NewQuery(a.db)
	user, err := query.RetrivePassword(reqDet.Email)
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Println(user.Password)
	log.Println(reqDet.Password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password) , []byte(reqDet.Password))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return user.UserID, nil
}
