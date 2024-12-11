package repository

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Srujankm12/SRproject/pkg/database"
	"github.com/google/uuid"
)

type FormDataRepo struct {
	db *sql.DB
}

func NewFormDataRepo(db *sql.DB) *FormDataRepo {
	return &FormDataRepo{
		db,
	}
}

func (fdr *FormDataRepo) ReportForm(r *http.Request) error {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return err
	}
	file1, _, err := r.FormFile("file1")
	if err != nil {
		return err
	}
	defer file1.Close()
	file2, _, err := r.FormFile("file2")
	if err != nil {
		return err
	}
	defer file2.Close()

	var data map[string]string
	err = json.Unmarshal([]byte(r.FormValue("json_data")), &data)
	if err != nil {
		return err
	}
	file1data, err := io.ReadAll(file1)
	if err != nil {
		return err
	}
	file2data, err := io.ReadAll(file2)
	if err != nil {
		return err
	}
	query := database.NewQuery(fdr.db)
	empid := uuid.NewString()
	err = query.StoreFile(empid, empid, file1data)
	if err != nil {
		return err
	}
	err = query.StoreFile(empid, empid, file2data)
	if err != nil {
		return err
	}
	err = query.StoreFormData(empid, data)

	if err != nil {
		return err
	}
	return nil
}
