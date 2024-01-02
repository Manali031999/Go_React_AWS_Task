package db

import (
	"fmt"
	"new_be/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetDBInstance() (*gorm.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %s", err.Error())
	}
	return db, nil
}

func InsertInstance(db *gorm.DB, instanceData models.Instance) error {
	// Check if the instance already exists
	var existingInstance models.Instance
	result := db.Where("instance_id = ?", instanceData.InstanceID).First(&existingInstance)
	if result.Error == nil {
		// Instance already exists, update the existing record with new data
		result := db.Model(&existingInstance).Updates(instanceData)
		if result.Error != nil {
			return fmt.Errorf("error updating existing instance data: %v", result.Error)
		}
		fmt.Println("Existing instance updated.")
		return nil
	}

	// Insert the instance data
	result = db.Create(&instanceData)
	if result.Error != nil {
		return fmt.Errorf("error inserting instance data: %v", result.Error)
	}
	return nil
}

func InsertMetricData(db *gorm.DB, metricData models.MetricData) error {
	// Check if the metric data already exists
	var existingMetricData models.MetricData
	result := db.Where("instance_id = ?", metricData.InstanceID).First(&existingMetricData)
	if result.Error == nil {
		// Metric Data already exists, update the existing record with new data
		result := db.Model(&existingMetricData).Updates(metricData)
		if result.Error != nil {
			return fmt.Errorf("error updating existing metric data: %v", result.Error)
		}
		fmt.Println("Existing metric data updated.")
		return nil
	}

	// Insert the new metric data
	result = db.Create(&metricData)
	if result.Error != nil {
		return fmt.Errorf("error inserting metric data: %v", result.Error)
	}
	return nil
}
