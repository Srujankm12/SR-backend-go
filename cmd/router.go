package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/Srujankm12/SRproject/internal/handlers"
	"github.com/Srujankm12/SRproject/internal/middlewares"
	"github.com/Srujankm12/SRproject/repository"
	"github.com/gorilla/mux"
)

func registerRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()
	router.Use(middlewares.CorsMiddleware)

	salescon := handlers.NewSalesHandler(repository.NewSalesRepository(db))
	router.HandleFunc("/sales/login", salescon.LoginHandler).Methods("POST")
	router.HandleFunc("/sales/logout", salescon.LogoutHandler).Methods("POST")
	router.HandleFunc("/sales/checkin", salescon.CheckInHandler).Methods("POST")
	router.HandleFunc("/sales/checkout", salescon.CheckOutHandler).Methods("POST")

	authCont := handlers.NewAuthController(repository.NewAuth(db))
	router.HandleFunc("/register", authCont.Register).Methods("POST")
	router.HandleFunc("/login", authCont.Login).Methods("POST")

	formCont := handlers.NewFormController(repository.NewFormDataRepo(db))
	router.HandleFunc("/submit", formCont.SubmitFormController).Methods("POST")
	router.HandleFunc("/getdata/{id}", formCont.FetchFormDataController).Methods("GET")

	adminCont := handlers.NewAdminHandler(repository.NewAdmin(db))
	router.HandleFunc("/admin/register", adminCont.AdminRegister).Methods("POST")
	router.HandleFunc("/admin/login", adminCont.AdminLogin).Methods("POST")

	adminf := handlers.NewAdminFHandler(repository.NewAdminRepository(db))
	router.HandleFunc("/admin/fetch", adminf.HandleAdminFetchFormData).Methods("GET")
	router.HandleFunc("/admin/delete/{id}", adminf.HandleDeleteEmployee).Methods("GET")

	excelcon := handlers.NewTechnicalFormExcelHandler(repository.NewExcelDownload(db))
	router.HandleFunc("/excel", excelcon.HandleDownloadExcel).Methods("GET")
	tempDir := "/tmp"
	if os.Getenv("OS") == "Windows_NT" {
		tempDir = os.Getenv("TEMP")
	}

	router.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir(tempDir)))).Methods("GET")

	return router
}
