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
	err := r.ParseMultipartForm(10 << 20)
	var datas1 , datas2 []byte
	if err != nil {
		return errors.New("failed to parse multipart form: " + err.Error())
	}
	erChan := make(chan error)
	dataChan1  := make(chan  []byte)
	dataChan2  := make(chan  []byte)
	go func() {
		file1, _, err := r.FormFile("file1")
		if err != nil {
			erChan <- err
		}
		defer file1.Close()
		file1Data, err := io.ReadAll(file1)
		if err != nil {
			erChan <- err
		}
		dataChan1 <- file1Data
		close(dataChan1)
	}()
	go func() {
		file2, _, err := r.FormFile("file2")
		if err != nil {
			erChan <- err
		}
		defer file2.Close()
		file2Data, err := io.ReadAll(file2)
		if err != nil {
			erChan <- err
		}
		dataChan2 <- file2Data
		close(dataChan2)
	}()
	if _ , ok := <-erChan; !ok {
		return err
	}
	if data1 , ok := <- dataChan1; !ok {
		return fmt.Errorf("Unable to validate data")
	}else{
		datas1 = data1
	}
	if data2 , ok := <- dataChan2; !ok {
		return fmt.Errorf("Unable to validate data")
	}else{
		datas2 = data2
	}
	var data models.FormData
	err = json.Unmarshal([]byte(r.FormValue("json_data")), &data)
	if err != nil {
		return errors.New("failed to unmarshal JSON data: " + err.Error())
	}

	query := database.NewQuery(fdr.db)
	data.EmployeeID = uuid.NewString()

	err = query.StoreFile(data.UserID , data.EmployeeID , "file1" , "file2" , datas1 , datas2)
	if err != nil {
		log.Println("Error storing file1:", err)
		return err
	}

	err = query.StoreFormData(data)
	if err != nil {
		log.Println("Error storing form data:", err)
		return err
	}

	return nil
}

// func (fdr *FormDataRepo) FetchFormData(r *http.Request) (*map[string]string, []byte, []byte, error) {
// 	// Parse the request parameters to get the employee ID (empID)
// 	err := r.ParseForm()
// 	if err != nil {
// 		return nil, nil, nil, errors.New("failed to parse form: " + err.Error())
// 	}

// 	userID := r.FormValue("user_id")
// 	if userID == "" {
// 		return nil, nil, nil, errors.New("employee ID is missing")
// 	}

// 	// Create a query instance to interact with the database
// 	query := database.NewQuery(fdr.db)

// 	// Fetch form data
// 	formData, err := query.FetchFormData(userID)
// 	if err != nil {
// 		log.Println("Error fetching form data:", err)
// 		return nil, nil, nil, err
// 	}

// 	// Fetch files associated with the empID
// 	file1Data, err := query.FetchFile(userID, "file1")
// 	if err != nil {
// 		log.Println("Error fetching file1:", err)
// 		return nil, nil, nil, err
// 	}

// 	file2Data, err := query.FetchFile(userID, "file2")
// 	if err != nil {
// 		log.Println("Error fetching file2:", err)
// 		return nil, nil, nil, err
// 	}

// 	// Return the form data, along with file contents
// 	return &formData, file1Data, file2Data, nil
// }
