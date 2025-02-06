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
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("unable to load env: %v", err)
	}

	// Create a new database connection
	conn := NewConnection()
	defer conn.DB.Close()

	// Create tables before starting the server
	query := database.NewQuery(conn.DB)
	err := query.CreateTables()
	if err != nil {
		log.Fatalf("Unable to create database tables: %v", err)
	}

	// Create a new router and apply the CORS middleware globally
	router := mux.NewRouter()
	router.Use(middlewares.CorsMiddleware)

	// Register your routes
	router = registerRouter(conn.DB)

	// Create and start the server
	server := &http.Server{
		Addr:    os.Getenv("PORT"), // Default to port if not set
		Handler: router,
	}

	// If PORT is not set in the .env file, default to 8080
	if os.Getenv("PORT") == "" {
		log.Fatal("PORT environment variable not set")
	}

	log.Printf("Server is running at port %s", os.Getenv("PORT"))
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Unable to start the server: %v", err)
	}
}
