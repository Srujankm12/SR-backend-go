package repository

import (
	"database/sql"
	"fmt"

	"github.com/Srujankm12/SRproject/internal/models"
	"github.com/xuri/excelize/v2"
)

type ExcelDownload struct {
	db *sql.DB
}

func NewExcelDownload(db *sql.DB) *ExcelDownload {
	return &ExcelDownload{
		db,
	}
}

func (e *ExcelDownload) FetchExcel() ([]models.DownloadExcel, error) {
	var data []models.DownloadExcel

	if e.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	rows, err := e.db.Query("SELECT * FROM formdata")
	if err != nil {
		fmt.Println("Database query error:", err)
		return nil, fmt.Errorf("Error executing database query: %v", err)

	}
	defer rows.Close()
	for rows.Next() {
		var record models.DownloadExcel
		if err := rows.Scan(&record.UserID, &record.EmployeeID, &record.ReportDate, &record.EmployeeName, &record.Premises, &record.SiteLocation, &record.ClientName, &record.ScopeOfWork, &record.WorkDetails, &record.JointVisits, &record.SupportNeeded, &record.StatusOfWork, &record.PriorityOfWork, &record.NextActionPlan, &record.TypeOfWork, &record.ClosingTime, &record.ContactPersonName, &record.ContactEmailID); err != nil {
			return nil, err
		}
		data = append(data, record)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("no data found")
	}
	return data, nil
}

func (e *ExcelDownload) CreateTechnialExcel() (*excelize.File, error) {
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing the file:", err)
		}
	}()
	data, err := e.FetchExcel()
	if err != nil {
		return nil, err
	}

	sheetName := "technical"
	index, err := file.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	file.SetActiveSheet(index)
	file.DeleteSheet("Sheet1")

	headers := []string{"UserID", "EmployeeID", "ReportDate", "EmployeeName", "Premises", "SiteLocation", "ClientName", "ScopeOfWork", "WorkDetails", "JointVisits", "SupportNeeded", "StatusOfWork", "PriorityOfWork", "NextActionPlan", "TypeOfWork", "ClosingTime", "ContactPersonName", "ContactEmailID"}
	for col, header := range headers {
		cell, err := excelize.CoordinatesToCellName(col+1, 1)
		if err != nil {
			return nil, err
		}
		file.SetCellValue(sheetName, cell, header)
	}

	for i, record := range data {
		row := i + 2
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), record.UserID)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), record.EmployeeID)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), record.ReportDate)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), record.EmployeeName)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", row), record.Premises)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", row), record.SiteLocation)
		file.SetCellValue(sheetName, fmt.Sprintf("G%d", row), record.ClientName)
		file.SetCellValue(sheetName, fmt.Sprintf("H%d", row), record.ScopeOfWork)
		file.SetCellValue(sheetName, fmt.Sprintf("I%d", row), record.WorkDetails)
		file.SetCellValue(sheetName, fmt.Sprintf("J%d", row), record.JointVisits)
		file.SetCellValue(sheetName, fmt.Sprintf("K%d", row), record.SupportNeeded)
		file.SetCellValue(sheetName, fmt.Sprintf("L%d", row), record.StatusOfWork)
		file.SetCellValue(sheetName, fmt.Sprintf("M%d", row), record.PriorityOfWork)
		file.SetCellValue(sheetName, fmt.Sprintf("N%d", row), record.NextActionPlan)
		file.SetCellValue(sheetName, fmt.Sprintf("O%d", row), record.TypeOfWork)
		file.SetCellValue(sheetName, fmt.Sprintf("P%d", row), record.ClosingTime)
		file.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), record.ContactPersonName)
		file.SetCellValue(sheetName, fmt.Sprintf("R%d", row), record.ContactEmailID)

	}
	if len(data) == 0 {
		return nil, fmt.Errorf("no data found")
	}
	return file, nil
}
