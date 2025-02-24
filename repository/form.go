package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Srujankm12/SRproject/internal/models"
	"github.com/Srujankm12/SRproject/pkg/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type FormDataRepo struct {
	db *sql.DB
}

func NewFormDataRepo(db *sql.DB) *FormDataRepo {
	return &FormDataRepo{
		db: db,
	}
}
func (fdr *FormDataRepo) SubmitFormData(r *http.Request) error {
	// Set a max file size limit for the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		return errors.New("failed to parse multipart form: " + err.Error())
	}

	// Create error and file data channels
	erChan := make(chan error, 2)
	dataChan1 := make(chan []byte)
	dataChan2 := make(chan []byte)

	// Goroutine to process file1 (optional)
	go func() {
		file1, _, err := r.FormFile("file1")
		if err != nil {
			if err == http.ErrMissingFile {
				dataChan1 <- nil // No file uploaded
			} else {
				erChan <- err
			}
			close(dataChan1)
			return
		}
		defer file1.Close()
		file1Data, err := io.ReadAll(file1)
		if err != nil {
			erChan <- err
			close(dataChan1)
			return
		}
		dataChan1 <- file1Data
		close(dataChan1)
	}()

	// Goroutine to process file2 (optional)
	go func() {
		file2, _, err := r.FormFile("file2")
		if err != nil {
			if err == http.ErrMissingFile {
				dataChan2 <- nil // No file uploaded
			} else {
				erChan <- err
			}
			close(dataChan2)
			return
		}
		defer file2.Close()
		file2Data, err := io.ReadAll(file2)
		if err != nil {
			erChan <- err
			close(dataChan2)
			return
		}
		dataChan2 <- file2Data
		close(dataChan2)
	}()

	// Receive data or errors
	var datas1, datas2 []byte
	var errOccurred error

	select {
	case errOccurred = <-erChan:
		if errOccurred != nil {
			return fmt.Errorf("error occurred: %v", errOccurred)
		}
	case datas1 = <-dataChan1:
	case datas2 = <-dataChan2:
	}

	// Unmarshal JSON data
	var data models.FormData
	err = json.Unmarshal([]byte(r.FormValue("json_data")), &data)
	if err != nil {
		return errors.New("failed to unmarshal JSON data: " + err.Error())
	}

	// Ensure employee ID is set
	data.EmployeeID = uuid.NewString()

	// Store files only if they exist
	query := database.NewQuery(fdr.db)
	err = query.StoreFile(data.UserID, data.EmployeeID, "file1", "file2", datas1, datas2)
	if err != nil {
		log.Println("Error storing files:", err)
		return err
	}

	// Store form data in the database
	err = query.StoreFormData(data)
	if err != nil {
		log.Println("Error storing form data:", err)
		return err
	}

	return nil
}

func (fdr *FormDataRepo) FetchFormData(r *http.Request) ([]models.FormData, error) {
	// Extract user ID from URL params
	userID := mux.Vars(r)["id"]
	log.Println("Extracted userID:", userID) // Debugging

	if userID == "" {
		log.Println("Error: Employee ID is missing")
		return nil, errors.New("employee ID is missing")
	}

	// Ensure database connection is valid
	if fdr.db == nil {
		log.Println("Error: Database connection is nil")
		return nil, errors.New("database connection error")
	}

	// Create a query instance to interact with the database
	query := database.NewQuery(fdr.db)

	// Fetch form data **filtered by userID**
	formData, err := query.FetchFormData(userID)
	if err != nil {
		log.Println("Error fetching form data:", err)
		return nil, err
	}

	log.Println("Successfully fetched form data for userID:", userID) // Debugging

	return formData, nil
}
