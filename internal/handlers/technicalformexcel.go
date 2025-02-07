package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Srujankm12/SRproject/repository"
	"github.com/xuri/excelize/v2"
)

type TechnicalFormExcelHandler struct {
	Repo *repository.ExcelDownload
}

func NewTechnicalFormExcelHandler(repo *repository.ExcelDownload) *TechnicalFormExcelHandler {
	return &TechnicalFormExcelHandler{
		Repo: repo,
	}
}

func (h *TechnicalFormExcelHandler) HandleDownloadExcel(w http.ResponseWriter, r *http.Request) {
	excelData, err := h.Repo.FetchExcel()
	if err != nil {
		http.Error(w, "Failed to fetch Excel data", http.StatusInternalServerError)
		return
	}

	file := excelize.NewFile()
	sheetName := "technical"
	file.NewSheet(sheetName)

	headers := []string{"UserID", "EmployeeID", "ReportDate", "EmployeeName", "Premises", "SiteLocation", "ClientName", "ScopeOfWork", "WorkDetails", "JointVisits", "SupportNeeded", "StatusOfWork", "PriorityOfWork", "NextActionPlan", "Result", "TypeOfWork", "ClosingTime", "ContactPersonName", "ContactEmailID"}

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		file.SetCellValue(sheetName, cell, header)
	}

	for i, record := range excelData {
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
		file.SetCellValue(sheetName, fmt.Sprintf("O%d", row), record.Result)
		file.SetCellValue(sheetName, fmt.Sprintf("P%d", row), record.TypeOfWork)
		file.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), record.ClosingTime)
		file.SetCellValue(sheetName, fmt.Sprintf("R%d", row), record.ContactPersonName)
		file.SetCellValue(sheetName, fmt.Sprintf("S%d", row), record.ContactEmailID)
	}

	tempDir := "/Users/bunny/Desktop/Finalyear/test/"
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		http.Error(w, fmt.Sprintf("Error creating temporary directory: %v", err), http.StatusInternalServerError)
		return
	}

	filePath := tempDir + "technical.xlsx"
	if err := file.SaveAs(filePath); err != nil {
		http.Error(w, fmt.Sprintf("Error saving file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("File saved. You can download it from: %s\n", filePath)))
}
