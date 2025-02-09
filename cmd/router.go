package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Srujankm12/SRproject/internal/handlers"
	"github.com/Srujankm12/SRproject/internal/middlewares"
	"github.com/Srujankm12/SRproject/repository"
	"github.com/gorilla/mux"
)

func registerRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()
	router.Use(middlewares.CorsMiddleware)

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
	router.HandleFunc("/Users/bunny/Desktop/Finalyear/test", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileName := vars["filename"]
		filePath := fmt.Sprintf("/Users/bunny/Desktop/Finalyear/test%s", fileName)

		http.ServeFile(w, r, filePath)
	}).Methods("GET")
	return router
}
