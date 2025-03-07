package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Srujankm12/SRproject/internal/middlewares"
	"github.com/Srujankm12/SRproject/pkg/database"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("unable to load env: %v", err)
	}

	conn := NewConnection()
	defer conn.DB.Close()

	query := database.NewQuery(conn.DB)
	err := query.CreateTables()
	if err != nil {
		log.Fatalf("Unable to create database tables: %v", err)
	}

	router := mux.NewRouter()
	router.Use(middlewares.CorsMiddleware)

	router = registerRouter(conn.DB)

	server := &http.Server{
		Addr:    os.Getenv("PORT"),
		Handler: router,
	}

	if os.Getenv("PORT") == "" {
		log.Fatal("PORT environment variable not set")
	}

	log.Printf("Server is running at port %s", os.Getenv("PORT"))
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Unable to start the server: %v", err)
	}
}
