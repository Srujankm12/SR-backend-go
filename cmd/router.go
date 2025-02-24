package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/Srujankm12/SRproject/internal/handlers"
	"github.com/Srujankm12/SRproject/internal/middlewares"
	"github.com/Srujankm12/SRproject/repository"
)

func registerRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(middlewares.CorsMiddleware)

	// Sales Routes
	salesRepo := repository.NewSalesRepository(db)
	salesHandler := handlers.NewSalesHandler(salesRepo.DB)
	router.HandleFunc("/sales/register", salesHandler.HandleCheckIn).Methods(http.MethodPost)
	router.HandleFunc("/sales/checkout", salesHandler.HandleCheckOut).Methods(http.MethodPost)

	// Auth Routes
	authRepo := repository.NewAuth(db)
	authHandler := handlers.NewAuthController(authRepo)
	router.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost)
	router.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)

	// Form Data Routes
	formRepo := repository.NewFormDataRepo(db)
	formHandler := handlers.NewFormController(formRepo)
	router.HandleFunc("/submit", formHandler.SubmitFormController).Methods(http.MethodPost)
	router.HandleFunc("/getdata/{id}", formHandler.FetchFormDataController).Methods(http.MethodGet)

	// Admin Routes
	adminRepo := repository.NewAdmin(db)
	adminHandler := handlers.NewAdminHandler(adminRepo)
	router.HandleFunc("/admin/register", adminHandler.AdminRegister).Methods(http.MethodPost)
	router.HandleFunc("/admin/login", adminHandler.AdminLogin).Methods(http.MethodPost)

	// Admin Fetch & Delete
	adminFormRepo := repository.NewAdminRepository(db)
	adminFormHandler := handlers.NewAdminFHandler(adminFormRepo)
	router.HandleFunc("/admin/fetch", adminFormHandler.HandleAdminFetchFormData).Methods(http.MethodGet)
	router.HandleFunc("/admin/delete/{id}", adminFormHandler.HandleDeleteEmployee).Methods(http.MethodDelete) // Changed to DELETE

	// Excel Download
	excelRepo := repository.NewExcelDownload(db)
	excelHandler := handlers.NewTechnicalFormExcelHandler(excelRepo)
	router.HandleFunc("/excel", excelHandler.HandleDownloadExcel).Methods(http.MethodGet)

	// File Server Setup
	tempDir := "/tmp"
	if os.Getenv("OS") == "Windows_NT" {
		tempDir = os.Getenv("TEMP")
	}
	router.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir(tempDir)))).Methods(http.MethodGet)

	return router
}
