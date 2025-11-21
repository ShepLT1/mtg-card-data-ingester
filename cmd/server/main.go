package main

import (
	"fmt"
"log"
	"net/http"

	"SWE1-project-data-ingester/internal/api"
	"SWE1-project-data-ingester/internal/data"
)

func main() {
	data.Connect()
	sqlDB, err := data.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database object: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	log.Println("✅ Connected to database")

	router := api.NewRouter()

	fmt.Println("Server running on http://localhost:8085")
	log.Fatal(http.ListenAndServe(":8085", router))
}
