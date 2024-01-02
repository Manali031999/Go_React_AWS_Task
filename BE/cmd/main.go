package main

import (
	"fmt"
	"log"
	"net/http"

	"new_be/db"
	"new_be/handlers"
	"new_be/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Load environment variables from the .env file
	if err := utils.LoadEnv(); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()

	// Migrate DB
	if err := db.MigrateDatabase(); err != nil {
		log.Fatal("Error migrating database:", err)
	}
	// Create DB with AWS data
	if err := utils.AwsAccess(); err != nil {
		log.Fatal("Error updating DB with AWS data:", err)
	}

	// Use the CORS middleware
	r.Use(handlers.CorsMiddleware)

	// Create routes to serve data
	handlers.InitRoutes(r)

	// Start the server
	port := ":8080"
	fmt.Println("Server is running on port", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
