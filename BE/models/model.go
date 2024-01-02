package models

import (
	"gorm.io/gorm"
)

// Instance represents the database model for instances
type Instance struct {
	gorm.Model
	InstanceID   string
	InstanceType string
	Region       string
}

// MetricData represents the database model for metrics
type MetricData struct {
	gorm.Model
	InstanceID string
	CPU        float64
	GraphData  []byte `gorm:"type:json"`
}
