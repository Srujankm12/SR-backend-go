package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/Srujankm12/SRproject/pkg/database"
	"github.com/google/uuid"
)

type FormDataRepo struct {
	db *sql.DB
}

func NewFormDataRepo(db *sql.DB) *FormDataRepo {
	return &FormDataRepo{
		db: db,
	}
}

func (fdr *FormDataRepo) ReportForm(r *http.Request) error {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		return errors.New("failed to parse multipart form: " + err.Error())
	}

	file1, _, err := r.FormFile("file1")
	if err != nil {
		return errors.New("failed to retrieve file1: " + err.Error())
	}
	defer file1.Close()

	file2, _, err := r.FormFile("file2")
	if err != nil {
		return errors.New("failed to retrieve file2: " + err.Error())
	}
	defer file2.Close()

	file1Data, err := io.ReadAll(file1)
	if err != nil {
		return errors.New("failed to read file1: " + err.Error())
	}

	file2Data, err := io.ReadAll(file2)
	if err != nil {
		return errors.New("failed to read file2: " + err.Error())
	}

	var data map[string]string
	err = json.Unmarshal([]byte(r.FormValue("json_data")), &data)
	if err != nil {
		return errors.New("failed to unmarshal JSON data: " + err.Error())
	}

	query := database.NewQuery(fdr.db)
	empID := uuid.NewString()

	err = query.StoreFile(empID, "file1", file1Data)
	if err != nil {
		log.Println("Error storing file1:", err)
		return err
	}

	err = query.StoreFile(empID, "file2", file2Data)
	if err != nil {
		log.Println("Error storing file2:", err)
		return err
	}

	err = query.StoreFormData(empID, data)
	if err != nil {
		log.Println("Error storing form data:", err)
		return err
	}

	return nil
}
