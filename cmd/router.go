package main

import (
	"database/sql"

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
	router.HandleFunc("/getdata/{id}" , formCont.FetchFormDataController).Methods("GET")
	return router
}
