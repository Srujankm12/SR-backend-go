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

	// Create channels for errors and file data
	erChan := make(chan error, 2) // Buffered to handle multiple errors
	dataChan1 := make(chan []byte)
	dataChan2 := make(chan []byte)

	// Goroutine to process file1
	go func() {
		file1, _, err := r.FormFile("file1")
		if err != nil {
			erChan <- err
			close(dataChan1) // Ensure channels are closed on error
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
		close(dataChan1) // Close after sending data
	}()

	// Goroutine to process file2
	go func() {
		file2, _, err := r.FormFile("file2")
		if err != nil {
			erChan <- err
			close(dataChan2) // Ensure channels are closed on error
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
		close(dataChan2) // Close after sending data
	}()

	// Wait for both file data or error
	var datas1, datas2 []byte
	var errOccurred error

	// Wait for data or error from channels
	select {
	case errOccurred = <-erChan:
		if errOccurred != nil {
			return fmt.Errorf("error occurred: %v", errOccurred)
		}
	case datas1 = <-dataChan1:
	case datas2 = <-dataChan2:
	}

	// Now we unmarshal the JSON data
	var data models.FormData
	err = json.Unmarshal([]byte(r.FormValue("json_data")), &data)
	if err != nil {
		return errors.New("failed to unmarshal JSON data: " + err.Error())
	}

	// Ensure employee ID is set
	data.EmployeeID = uuid.NewString()

	// Store files in the database
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
	// Parse the request parameters to get the employee ID (empID)
	// err := r.ParseForm()
	// if err != nil {
	// 	return nil, errors.New("failed to parse form: " + err.Error())
	// }

	// userID := r.FormValue("user_id")
	var userID = mux.Vars(r)["id"]
	if userID == "" {
		return nil, errors.New("employee ID is missing")
	}

	// Create a query instance to interact with the database
	query := database.NewQuery(fdr.db)

	// Fetch form data
	formData, err := query.FetchFormData()
	if err != nil {
		log.Println("Error fetching form data:", err)
		return nil, err
	}
	// Return the form data, along with file contents
	return formData, nil
}
