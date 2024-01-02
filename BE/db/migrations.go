package db

import (
	"fmt"
	"log"
	"new_be/models"
)

func MigrateDatabase() error {
	db, err := GetDBInstance()
	if err != nil {
		log.Fatalf("Error getting database instance: %v", err)
	}

	// AutoMigrate tables
	db.AutoMigrate(&models.Instance{}, &models.MetricData{})

	fmt.Println("Database migration successful")
	return nil
}
